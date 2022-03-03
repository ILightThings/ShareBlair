package options

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type UserFlags struct {
	Target          string
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
}

func (x *UserFlags) DetermineTarget() []string {
	hosts, err := hosts(x.Target)
	if err == nil {
		if x.Verbose {
			fmt.Printf("Parsed CIDR from input.\nTargets:\n")
			for y, x := range hosts {
				fmt.Printf("%d - %s\n", y, x)
			}

		}
		x.targetsParsed = hosts
		return hosts
	}
	filelines, err := readTargetsFromFile(x.Target)
	if err == nil {
		if x.Verbose {
			fmt.Printf("Parsed Targets from file.\nTargets:\n")
			for y, x := range filelines {
				fmt.Printf("%d - %s\n", y, x)
			}
		}
		x.targetsParsed = filelines
		return filelines
	}

	var single []string
	single = append(single, x.Target)

	x.targetsParsed = single
	if x.Verbose {
		fmt.Printf("Parsed Target from input.\nTarget:\n")
		for _, x := range single {
			fmt.Println(x)
		}
	}
	return single

}

func hosts(cidr string) ([]string, error) { // https://gist.github.com/kotakanbe/d3059af990252ba89a82
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
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

func readTargetsFromFile(s string) ([]string, error) {
	var lines []string
	file, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		x := strings.TrimSpace(scanner.Text())
		lines = append(lines, x)

	}
	return lines, nil

}
