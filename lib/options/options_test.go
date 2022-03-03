package options

import "testing"

func TestReturnTargets(t *testing.T) {

	testCidr := UserFlags{
		Target: "192.168.1.0/28",
	}
	testFile := UserFlags{
		Target: "targets_test.txt",
	}

	testSingleIP := UserFlags{
		Target: "192.168.60.60",
	}

	singleHost := UserFlags{
		Target: "foxtrot.ggray.info",
	}
	singleHost.DetermineTarget()
	testSingleIP.DetermineTarget()
	testCidr.DetermineTarget()
	testFile.DetermineTarget()

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

}
