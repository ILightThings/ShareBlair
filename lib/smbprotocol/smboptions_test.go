package smbprotocol

import (
	"testing"

	"github.com/ilightthings/shareblair/lib/options"
)

func TestConnectionAuthenticated(t *testing.T) {
	flags := &options.UserFlags{
		Target:   "127.0.0.1",
		User:     "gameandwatch",
		Password: "password",
		Port:     445,
	}
	_, guestaccess, err := ConnectionAuthenticated(flags)
	if err != nil {
		t.Error(err)
	}
	if guestaccess == false {
		t.Errorf("Cannot log into account %s with password %s", flags.User, flags.Password)
	}

}

func TestConnectionGuest(t *testing.T) {
	flags := &options.UserFlags{
		Target:   "127.0.0.1",
		User:     "Guest",
		Password: "",
		Port:     445,
	}
	_, guestaccess, err := ConnectionAuthenticated(flags)
	if err != nil {
		t.Error(err)
	}
	if guestaccess == true {
		t.Error("Guest account was able to log in.... That should not have been able to happen")
	}

}
