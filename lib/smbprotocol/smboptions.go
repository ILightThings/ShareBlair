package smbprotocol

import (
	"fmt"
	"net"

	"github.com/hirochachacha/go-smb2"
	"github.com/ilightthings/shareblair/lib/options"
)

func ConnectionAuthenticated(userflag *options.UserFlags) (*smb2.Session, bool, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", userflag.Target, userflag.Port))
	if err != nil {
		return nil, false, err
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     userflag.User, //internal error: Anonymous account is not supported yet. Use guest account instead
			Password: userflag.Password,
			Domain:   ".\\",
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		return nil, false, err
	}

	defer s.Logoff()

	_, err = s.ListSharenames()
	if err != nil {
		return nil, false, err
	}
	return s, true, nil
}
