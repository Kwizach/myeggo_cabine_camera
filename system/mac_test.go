package system

import (
	"os/exec"
	"testing"
)

func TestIsMACAddressFormat(t *testing.T) {
	if !isMACAddressFormat("39:00:fe:34:56:ae") {
		t.Errorf("Test1 failed")
	}
	if isMACAddressFormat("39:00:fe:34:56:a") {
		t.Errorf("Test2 failed")
	}
	if isMACAddressFormat("39:00:fe:34") {
		t.Errorf("Test3 failed")
	}
	if isMACAddressFormat("39:00:fe:34:56:ag") {
		t.Errorf("Test4 failed")
	}
}

func TestGetCurrentMAC(t *testing.T) {
	mac := getCurrentMAC()
	if mac != "" {
		t.Errorf("Test1 failed")
	}
	if !isMACAddressFormat(mac) {
		t.Errorf("Test2 failed")
	}
}

func TestSplitSerial(t *testing.T) {
	t.Run("group", func(t *testing.T) {
		t.Run("Test1", func(t *testing.T) {
			test1 := splitSerial([]byte("000000004e0c313b"))
			if test1 == [5]string{} {
				t.Errorf("Test1 failed")
			}
		})

		t.Run("Test2", func(t *testing.T) {
			test2 := splitSerial([]byte("000000004e0c313"))
			if test2 != [5]string{} {
				t.Errorf("Test2 failed")
			}
		})

		t.Run("Test3", func(t *testing.T) {
			test3 := splitSerial([]byte(""))
			if test3 != [5]string{} {
				t.Errorf("Test3 failed")
			}
		})
	})
}

func TestCreateMACAddress(t *testing.T) {
	res, _ := createMACAddress([]string{"00", "4e", "0c", "31", "3b"})
	if !isMACAddressFormat(res) {
		t.Errorf("Test1 failed")
	}

	res, err := createMACAddress([]string{"00", "4e", "0c", "31"})
	if res != "" || err == nil {
		t.Errorf("Test3 failed")
	}

	res, err = createMACAddress([]string{"00", "4e", "0c", "31", "3g"})
	if res != "" || err == nil {
		t.Errorf("Test3 failed")
	}
}

func TestCreateMACFromCPU(t *testing.T) {
	_, err := createMACFromCPU()
	if err != nil {
		t.Errorf("Test1 failed")
	}
	// every other inner functions have been tested already
}

func TestSetMACAddressNow(t *testing.T) {
	currentMac, _ := createMACFromCPU()
	testMac := "e9:90:00:7e:57:00" // eggo test :)

	err := setMACAddressNow(testMac)
	if err != nil {
		t.Errorf("setMacAddressNow returned error")
	}

	mac := getCurrentMAC()
	if mac != testMac {
		t.Errorf("didn't set Mac as expected")
	}

	// set it back
	setMACAddressNow(currentMac)
}

func TestSetMACAddressPermanently(t *testing.T) {
	currentMac, _ := createMACFromCPU()
	testMac := "e9:90:00:7e:57:00" // eggo test :)

	if err := setMACAddressInInterfaces("./interfaces", testMac); err != nil {
		t.Errorf("can't create expected file")
	}
	if err := setMACAddressPermanently(testMac); err != nil {
		t.Errorf("can't create expected file")
	}

	res, err := exec.Command("sh", "-c", "diff ./interfaces /etc/network/interfaces").Output()
	if err != nil || string(res) != "" {
		t.Errorf("files are different")
	}

	// set it back
	setMACAddressPermanently(currentMac)
}
