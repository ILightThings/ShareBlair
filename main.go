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
	target := parser.String("t", "target", &argparse.Options{Required: true, Help: "Hostname, IP, CIDR or file of targets"})
	user := parser.String("u", "user", &argparse.Options{Required: false, Help: "User to authenticate with"})
	domain := parser.String("d", "domain", &argparse.Options{Required: false, Help: "Domain to authenticate with"})
	password := parser.String("p", "password", &argparse.Options{Required: false, Help: "Password to authenticate with"})
	hash := parser.String("", "hash", &argparse.Options{Required: false, Help: "Hash to authenticate with"})
	port := parser.Int("", "port", &argparse.Options{Required: false, Default: 445, Help: "Port to connect to"})
	verbose := parser.Flag("v", "verbose", &argparse.Options{Required: false, Help: "Add verbosity", Default: false})
	//TODO add timeout for func (r *Target) InitTCP()

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	userflags := &options.UserFlags{
		Target:   *target,
		User:     *user,
		Domain:   *domain,
		Password: *password,
		Hash:     *hash,
		Port:     *port,
		Verbose:  *verbose,
	}

	//var hostsArray []smbprotocol.Target

	userflags.DetermineTarget()
	var scope []smbprotocol.Target

	for _, x := range userflags.TargetsParsed {
		var singleTarget smbprotocol.Target
		singleTarget.Initialize(userflags, x)
		scope = append(scope, singleTarget)

	}

	for _, y := range scope {
		err := y.InitTCP()
		if err == nil {
			err := y.InitSMBAuth()
			if err == nil {
				shares, err := y.GetShareList()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Guest Access: %t\n", y.GuestAccessCheck())
				for _, share := range shares {
					fmt.Println(share)
				}
				y.CloseSMBSession()
			}
			y.CloseTCP()
		}

	}
}
