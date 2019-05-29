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

func TestIsDefaultHostName(t *testing.T) {
	if IsDefaultHostName() {
		t.Errorf("IsDefaultHostName should not be raspberrypi")
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
	err := changeHostName(defaultHostname)
	if err != nil {
		t.Errorf("Can't set hostname to raspberrypi")
	}

	err = CreateNewHostName()
	if err != nil {
		t.Errorf("Can't set new hostname based on cpuInfo")
	}

	if IsDefaultHostName() {
		t.Errorf("Hostname is still default")
	}
}

func TestChangeEtcHosts(t *testing.T) {}

func TestChangeEtcHostname(t *testing.T) {}

func TestChangeHostNameCtl(t *testing.T) {}
