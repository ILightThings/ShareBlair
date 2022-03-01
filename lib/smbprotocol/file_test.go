package smbprotocol

import (
	"testing"

	"github.com/ilightthings/shareblair/lib/options"
)

func TestFullObject(t *testing.T) {
	authuser := &options.UserFlags{
		User:     "Gameandwatch",
		Target:   "127.0.0.1",
		Password: "password",
		Domain:   "",
		Port:     445,
	}

	testTarget := Target{}
	err := testTarget.InitTCP(authuser)
	if err != nil {
		t.Errorf("Failed to connect: %s\n", err)
	}

	err1 := testTarget.InitSMBAuth(authuser)
	if err1 != nil {
		t.Errorf("Failed to authenticate: %s \n", err)
	}

	list, err := testTarget.GetShareList()
	if err != nil {
		t.Error("Error getting shares")
	}
	if len(list) < 1 {
		t.Errorf("Incorrect amount of shares. Expecting 1. Got %d\n", len(list))
	}

	guestAccess := testTarget.GuestAccessCheck()
	if guestAccess == true {
		t.Error("Somehow guest access is enabled here.....")
	}

	testTarget.CloseSMBSession()
	testTarget.CloseTCP()

}
