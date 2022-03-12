package main

import (
	"os"

	"github.com/akamensky/argparse"
	"github.com/ilightthings/shareblair/lib/options"
	"github.com/ilightthings/shareblair/lib/report"
	"github.com/ilightthings/shareblair/lib/smbprotocol"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	parser := argparse.NewParser("shareblair", "")
	target := parser.String("t", "target", &argparse.Options{Required: true, Help: "Hostname, IP, CIDR or file of targets"})
	user := parser.String("u", "user", &argparse.Options{Required: false, Help: "User to authenticate with"})
	domain := parser.String("d", "domain", &argparse.Options{Required: false, Help: "Domain to authenticate with"})
	password := parser.String("p", "password", &argparse.Options{Required: false, Help: "Password to authenticate with"})
	hash := parser.String("", "hash", &argparse.Options{Required: false, Help: "Hash to authenticate with"})
	port := parser.Int("", "port", &argparse.Options{Required: false, Default: 445, Help: "Port to connect to"})
	verbose := parser.Int("v", "verbose", &argparse.Options{Required: false, Help: "Verbosity Level. 0 - Debug, 1 - Info, Warn - 3, Error - 4", Default: 2})
	maxDepth := parser.Int("", "maxdepth", &argparse.Options{Required: false, Help: "Max Recursive Depth for Share Scanning. 0 will only scan the top level folders.", Default: 5})

	//TODO add timeout for func (r *Target) InitTCP()
	//TODO add out file location

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	switch *verbose {
	case 6:
		logger.SetLevel(logrus.TraceLevel)
	case 0:
		logger.SetLevel(logrus.DebugLevel)
	case 1:
		logger.SetLevel(logrus.InfoLevel)
	case 2:
		logger.SetLevel(logrus.WarnLevel)
	case 3:
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.WarnLevel)
	}

	userflags := &options.UserFlags{
		Target:   *target,
		User:     *user,
		Domain:   *domain,
		Password: *password,
		Hash:     *hash,
		Port:     *port,
		MaxDepth: *maxDepth,
		Logging:  *logger,
	}
	// TODO Implement better logging system. Maybe logrus
	//var hostsArray []smbprotocol.Target

	userflags.DetermineTarget()
	var scope []smbprotocol.Target

	for _, x := range userflags.DetermineTarget() {
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
				for z := range shares {
					err := y.ListOfShares[z].InitializeShare(y.ConnectionSMB, y.UserFlag)
					if err != nil {
						y.ListOfShares[z].UserRead = false
					} else {
						y.ListOfShares[z].DirWalk(y.HostDestination)
						y.ListOfShares[z].UnmountShare()
					}

				}

			}
			y.CloseTCP()
		}
		if y.ConnectionTCP_OK {
			report.MakeJSON(&y, userflags)
		}

	}
}
