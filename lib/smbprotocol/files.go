package smbprotocol

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/ilightthings/shareblair/lib/options"
)

type folder struct {
	Path            []string
	ListOfFolders   []folder
	ListOfFiles     []file
	ReadAccess      bool
	WriteAccess     bool
	NumberOfFiles   int
	NumberOfFolders int
	NumberOfItems   int
}

type file struct {
	Path       string
	FolderPath string
	FilePath   string
	FileName   string
	Size       int
}

type Target struct {
	HostDestination  string
	ResolvedIP       net.IP
	UserFlag         *options.UserFlags
	ConnectionTCP    net.Conn
	ConnectionTCP_OK bool
	ConnectionSMB    *smb2.Session
	ConnectionSMB_OK bool
	GuestOnly        bool
	GuestAccess      bool
	ListOfShares     []Share
}

// TODO, add verbose
func (r *Target) Initialize(f *options.UserFlags, target string) error {
	r.UserFlag = f
	r.HostDestination = target

	if r.HostDestination == "" {
		return errors.New("no target has been given")
	}

	// Check for username + password or hash
	if r.UserFlag.User != "" {
		if r.UserFlag.Hash == "" && r.UserFlag.Password == "" {
			return errors.New("no password or hash supplied")
		}
	}

	// Check for guest only test
	if r.UserFlag.User == "" {
		r.GuestOnly = true
	} else {
		r.GuestOnly = false
	}

	return nil

}

func (r *Target) InitTCP() error {
	dstNet := fmt.Sprintf("%s:%d", r.HostDestination, r.UserFlag.Port)
	if r.UserFlag.Verbose {
		fmt.Printf("Attempting TCP Connection to %s\n", dstNet)
	}
	conn, err := net.DialTimeout("tcp", dstNet, 3*time.Second)
	if err != nil {
		if r.UserFlag.Verbose {
			fmt.Printf("Failed TCP Connection to %s\n", dstNet)
		}

		r.ConnectionTCP_OK = false
		return err
	} else {
		if r.UserFlag.Verbose {
			fmt.Printf("Sucessful TCP Connection to %s\n", dstNet)
		}
		r.ConnectionTCP_OK = true
		r.ConnectionTCP = conn
		return nil
	}

}

func (r *Target) CloseTCP() error {
	err := r.ConnectionTCP.Close()
	if err != nil {
		return err
	}
	r.ConnectionTCP_OK = false
	return nil

}

func (r *Target) InitSMBAuth(f *options.UserFlags) error {
	smbConnectionOptions := &smb2.NTLMInitiator{
		User:   f.User,
		Domain: f.Domain,
	}

	if r.UserFlag.Password != "" {
		smbConnectionOptions.Password = r.UserFlag.Password
	} else {
		var newerr error
		smbConnectionOptions.Hash, newerr = hex.DecodeString(r.UserFlag.Hash)
		if newerr != nil {
			return errors.New("could not encode hash")
		}
	}

	smbConnection := &smb2.Dialer{
		Initiator: smbConnectionOptions,
	}
	s, err := smbConnection.Dial(r.ConnectionTCP)
	if err != nil {
		r.ConnectionSMB_OK = false
		return err
	} else {
		r.ConnectionSMB_OK = true
		r.ConnectionSMB = s
		return nil

	}
}

func (r *Target) CloseSMBSession() error {
	err := r.ConnectionSMB.Logoff()
	if err != nil {
		return err
	}
	r.ConnectionSMB_OK = false
	return nil
}

func (r *Target) GetShareList() ([]string, error) {
	list, err := r.ConnectionSMB.ListSharenames()
	if err != nil {
		return nil, err
	}
	for _, x := range list {
		var shareFolder Share
		shareFolder.ShareName = x
		r.ListOfShares = append(r.ListOfShares, shareFolder)
	}
	return list, nil
}

func (r *Target) GuestAccessCheck() bool {
	guestOptions := &smb2.NTLMInitiator{
		User:     "Guest",
		Password: "",
		Domain:   "",
	}
	guestConnect := smb2.Dialer{
		Initiator: guestOptions,
	}

	conn, err := guestConnect.Dial(r.ConnectionTCP)
	if err != nil {
		r.GuestAccess = false
		return false
	} else {
		conn.Logoff() // TODO Maybe don't close this
		//Might need to just give this connection back if I want to do guest access to shares testing
		return true

	}

}
