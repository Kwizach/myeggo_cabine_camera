package system

import (
	"testing"
)

func TestMyID(t *testing.T) {
	myID = "toto"
	if MyID() != myID {
		t.Errorf("MyID should return myID")
	}
}

func TestIsGoodHostName(t *testing.T) {
	if !IsGoodHostName() {
		goodHostname, err := getCPUSerial()
		if err != nil {
			t.Errorf("Can't retriever CPUSerial")
			return
		}
		t.Errorf("IsGoodHostName should be %s", goodHostname)
	}
}

func TestCreateNewHostName(t *testing.T) {}

func TestGetCPUSerial(t *testing.T) {
	newName, err := getCPUSerial()
	if err != nil || newName == "rpi-" {
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

	if !IsGoodHostName() {
		t.Errorf("Hostname is still default")
	}
}

func TestChangeEtcHosts(t *testing.T) {}

func TestChangeEtcHostname(t *testing.T) {}

func TestChangeHostNameCtl(t *testing.T) {}
