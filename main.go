package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/ilightthings/shareblair/lib/options"
	"github.com/ilightthings/shareblair/lib/smbprotocol"
)

func main() {
	parser := argparse.NewParser("shareblair", "")
	target := parser.String("t", "target", &argparse.Options{Required: true, Help: "IP Target to scan for SMB Shares"})
	user := parser.String("u", "user", &argparse.Options{Required: false, Help: "User to authenticate with"})
	domain := parser.String("d", "domain", &argparse.Options{Required: false, Help: "Domain to authenticate with"})
	password := parser.String("p", "password", &argparse.Options{Required: false, Help: "Password to authenticate with"})
	hash := parser.String("", "hash", &argparse.Options{Required: false, Help: "Hash to authenticate with"})
	port := parser.Int("", "port", &argparse.Options{Required: false, Default: 445, Help: "Port to connect to"})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	userflags := &options.UserFlags{
		Target:   *target,
		User:     *user,
		Domain:   *domain,
		Password: *password,
		Hash:     *hash, //
		Port:     *port,
	}

	//var hostsArray []smbprotocol.Target

	singlehost := &smbprotocol.Target{}
	err = singlehost.Initialize(userflags, userflags.Target)
	if err != nil {
		log.Fatal(err)
	}
	err = singlehost.InitTCP(userflags)
	if err != nil {
		log.Fatal(err)
	}
	err = singlehost.InitSMBAuth(userflags)
	if err != nil {
		log.Fatal(err)
	}

	shares, newerr := singlehost.GetShareList()
	if newerr != nil {
		log.Fatal(newerr)
	}

	for _, x := range shares {
		fmt.Println(x)
	}

	if singlehost.GuestAccessCheck() {
		fmt.Println("Guest Access is enabled")
	} else {
		fmt.Println("Guest access is disabled")
	}

}

// TODO, Add a parser for the target that will detect hostname vs IP vs CIDR vs File
