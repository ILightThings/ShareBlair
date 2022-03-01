package smbprotocol

import (
	"errors"
	"fmt"
	"net"

	"github.com/hirochachacha/go-smb2"
	"github.com/ilightthings/shareblair/lib/options"
)

type folder struct {
	listoffolders   []folder
	ListOfFiles     []file
	ReadAccess      bool
	WriteAccess     bool
	NumberOfFiles   int
	NumberOfFolders int
	NumberOfItems   int
}

type file struct {
	FolderPath string
	FilePath   string
	FileName   string
	Size       int
}

type Target struct {
	HostDestination  string
	ResolvedIP       net.IP
	User             string
	Password         string
	Domain           string
	Hash             string
	ConnectionTCP    net.Conn
	ConnectionTCP_OK bool
	ConnectionSMB    *smb2.Session
	ConnectionSMB_OK bool
	GuestOnly        bool
	GuestAccess      bool
	ListOfShares     []string
	// TODO, replace user,password,domain,hash with userflag object. This will increase memory usages as it will copy it per target.
}

// TODO, add verbose
func (r *Target) Initialize(f *options.UserFlags, target string) error {
	r.User = f.User
	r.Password = f.Password
	r.Domain = f.Domain
	r.Hash = f.Hash
	r.HostDestination = target

	if r.HostDestination == "" {
		return errors.New("no target has been given")
	}

	// Check for username + password or hash
	if r.User != "" {
		if r.Password == "" && r.Hash == "" {
			return errors.New("no password or hash supplied")
		}
	}

	// Check for guest only test
	if r.User == "" {
		r.GuestOnly = true
	} else {
		r.GuestOnly = false
	}

	return nil

}

// TODO, Rebuild code so each target is an object and the methods defined in smboptions.go are implemented as object methods.

func (r *Target) InitTCP(f *options.UserFlags) error {
	dstNet := fmt.Sprintf("%s:%d", f.Target, f.Port)
	conn, err := net.Dial("tcp", dstNet)
	if err != nil {
		r.ConnectionTCP_OK = false
		return err
	} else {
		r.ConnectionTCP_OK = false
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
		User:     f.User,
		Password: f.Password,
		Domain:   f.Domain,
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
	r.ListOfShares = list
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
