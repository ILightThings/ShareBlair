package options

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserFlags struct {
	Target          string `validate:"ip|cidr|file"`
	targetsParsed   []string
	Threads         int
	Verbose         bool
	User            string
	Domain          string
	Password        string
	Hash            string
	Port            int
	MaxDepth        int
	OutFileLocation string
	Logging         logrus.Logger `validate:"required"`
}

func (x *UserFlags) DetermineTarget() []string {
	v := validator.New()
	err := v.Var(x.Target, "cidr")
	if err == nil {
		hosts := hosts_cird(x.Target)
		x.Logging.Info("Parsed CIDR from input")
		for y, z := range hosts {
			x.Logging.Info(fmt.Sprintf("%d - %s", y, z))

		}
		x.targetsParsed = hosts
		return hosts

	}
	err = v.Var(x.Target, "file")
	if err == nil {
		x.Logging.Info("Parsed Targets from file.")
		filelines := readTargetsFromFile(x.Target, x)
		x.targetsParsed = filelines
		if filelines != nil {
			return filelines
		}

	}

	err1 := v.Var(x.Target, "ip|hostname")
	if err1 == nil {
		var single []string
		single = append(single, x.Target)

		x.targetsParsed = single
		x.Logging.Info("Parsed Target from input.")
		for z := range single {
			x.Logging.Info(single[z])
		}

		return single

	}
	x.Logging.Fatal("No Targets Proccessed. Invalid Input")
	return nil

}

func hosts_cird(cidr string) []string { // https://gist.github.com/kotakanbe/d3059af990252ba89a82
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1]
}

//  http://play.golang.org/p/m8TNTtygK0
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func readTargetsFromFile(s string, userflag *UserFlags) []string {
	var lines []string
	file, err := os.Open(s)
	if err != nil {
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		v := validator.New()
		x := strings.TrimSpace(scanner.Text())
		err := v.Var(x, "ip|hostname")
		if err == nil {
			userflag.Logging.Info(fmt.Sprintf("%s added", x))
			lines = append(lines, x)

		} else {

			userflag.Logging.Info(fmt.Sprintf("%s is invalid. Skipping", x))

		}
	}
	userflag.Logging.Info("End of File")
	return lines

}
