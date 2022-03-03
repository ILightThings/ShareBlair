package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ilightthings/shareblair/lib/options"
	"github.com/ilightthings/shareblair/lib/smbprotocol"
)

func MakeJSON(target *smbprotocol.Target, o *options.UserFlags) {
	b, err := json.Marshal(target)
	if err != nil {
		log.Fatal(err)
	}
	fileDst := fmt.Sprintf("%s_SMBShares.json", target.HostDestination)
	ioutil.WriteFile(fileDst, b, 0644)
}
