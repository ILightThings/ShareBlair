package options

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestReturnTargets(t *testing.T) {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
	})

	testCidr := UserFlags{
		Target:  "192.168.1.0/28",
		Logging: *log,
	}
	testFile := UserFlags{
		Target:  "targets_test.txt",
		Logging: *log,
	}

	testSingleIP := UserFlags{
		Target:  "192.168.60.60",
		Logging: *log,
	}

	singleHost := UserFlags{
		Target:  "foxtrot.ggray.info",
		Logging: *log,
	}

	Invalid := UserFlags{
		Target:  "whatevercouldgowrong",
		Logging: *log,
	}

	cidrResult := testCidr.DetermineTarget()
	fileResult := testFile.DetermineTarget()
	ipResult := testSingleIP.DetermineTarget()
	singleHostResult := singleHost.DetermineTarget()

	if len(cidrResult) != 14 {
		t.Errorf("CIDR expect len 16, got %d", len(cidrResult))
	}
	if len(fileResult) != 3 {
		t.Errorf("File expect len 16, got %d", len(fileResult))
	}
	if len(ipResult) != 1 {
		t.Error("Should be 1")
	}
	if len(singleHostResult) != 1 {
		t.Error("Should be 1")
	}
	Invalid.DetermineTarget()

}
