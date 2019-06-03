package system

import (
	"testing"
)

func TestMyID(t *testing.T) {
	cpuID = "TEST_ID"
	if MyID() != "rpi-"+cpuID {
		t.Errorf("MyID should return myID")
	}
	cpuID = ""
}

func TestIsHostNameGood(t *testing.T) {
	if !IsHostNameGood() {
		goodHostname, err := getCPUSerial()
		if err != nil {
			t.Errorf("Can't retriever CPUSerial")
			return
		}
		t.Errorf("IsHostNameGood should be %s", goodHostname)
	}
}

func TestCreateNewHostName(t *testing.T) {}

func TestGetCPUSerial(t *testing.T) {
	_, err := getCPUSerial()
	if err != nil {
		t.Errorf("Can't retrieve CPU serial")
	}
}

func TestChangeHostName(t *testing.T) {
	err := changeHostName("rpi-test")
	if err != nil {
		t.Errorf("Can't set hostname to rpi-test")
		return
	}

	err = CreateNewHostName()
	if err != nil {
		t.Errorf("Can't set new hostname based on cpuInfo")
		return
	}

	if !IsHostNameGood() {
		t.Errorf("Hostname is still default")
	}
}

func TestChangeEtcHosts(t *testing.T) {}

func TestChangeEtcHostname(t *testing.T) {}

func TestChangeHostNameCtl(t *testing.T) {}
