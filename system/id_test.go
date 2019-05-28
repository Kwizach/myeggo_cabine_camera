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
