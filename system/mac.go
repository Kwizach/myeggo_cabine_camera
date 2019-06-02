package system

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const myEggoSignature string = "e9"

// MACAddress of the interface... need to check if we need to export it (in case for now)
var MACAddress string

// IsCurrentMACGood check if MAC address is based on CPU serial
func IsCurrentMACGood() bool {
	currentMAC := getCurrentMAC()
	if currentMAC == "" {
		return false
	}

	espectedMAC, err := createMACFromCPU()
	if err != nil {
		return false
	}

	return espectedMAC == currentMAC
}

// SetMACAddress on the RPI
func SetMACAddress() error {
	mac, err := createMACFromCPU()
	if err != nil {
		return err
	}

	if err := setMACAddressNow(mac); err != nil {
		return err
	}
	if err := setMACAddressPermanently(mac); err != nil {
		return err
	}

	return nil
}

func getCurrentMAC() string {
	res, err := exec.Command("sh", "-c", "ifconfig eth0 | grep eth0 | awk '{print $NF}'").Output()
	if err != nil {
		return ""
	}
	return string(res)
}

func createMACFromCPU() (string, error) {
	if MACAddress != "" {
		return MACAddress, nil
	}

	var (
		serial []byte
		err    error
	)

	if cpuID == "" {
		serial, err = getCPUSerial()
		if err != nil {
			return "", err
		}
	} else {
		serial = []byte(cpuID)
	}

	splS := splitSerial(serial)
	if splS == [5]string{} {
		return "", errors.New("Wrong serial format")
	}

	mac, err := createMACAddress(splS[:])
	if err != nil {
		return "", err
	}
	if !isMACAddressFormat(mac) {
		return "", errors.New("Wrong MAC Format")
	}
	MACAddress = mac

	return mac, nil
}

// splitSerial split CPU serial number into [5]string
func splitSerial(serial []byte) [5]string {
	var tmp [9]string
	var res [5]string

	if len(serial) != 16 {
		return [5]string{}
	}

	spl := serial
	for i := 0; len(spl) > 0; i += 2 {
		tmp[i/2] = string(spl[:2])
		spl = serial[i:]
	}

	if len(tmp) != 9 {
		return [5]string{}
	}

	copy(res[:], tmp[4:])

	return res
}

// createMAC from splitSerial
// splS should be [5]string
func createMACAddress(splS []string) (string, error) {
	if len(splS) != 5 {
		return "", errors.New("Wrong input length")
	}

	re := regexp.MustCompile(`^([0-9A-Fa-f]{2})$`)
	for _, v := range splS {
		if !re.MatchString(v) {
			return "", errors.New("Wrong input datas")
		}
	}

	return strings.Join(append([]string{myEggoSignature}, splS...), ":"), nil
}

func isMACAddressFormat(mac string) bool {
	re := regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$`)
	return re.MatchString(mac)
}

func setMACAddressNow(mac string) error {
	// Down the interface
	err := exec.Command("sh", "-c", "ifconfig eth0 down").Run()
	if err != nil {
		return err
	}

	// Change MAC Address
	err = exec.Command("sh", "-c", fmt.Sprintf("ifconfig eth0 hw %s", mac)).Run()
	if err != nil {
		return err
	}

	// Up the interface
	err = exec.Command("sh", "-c", "ifconfig eth0 up").Run()
	if err != nil {
		return err
	}
	return nil
}

func setMACAddressPermanently(mac string) error {
	return setMACAddressInInterfaces("/etc/network/interfaces", mac)
}

func setMACAddressInInterfaces(fileURL string, mac string) error {
	input := "allow-hotplug eth0\n" +
		"iface eth0 inet dhcp\n" +
		fmt.Sprintf("  hwaddress ether %s", mac)

	f, err := os.OpenFile(fileURL, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer f.Close()

	// Erase file content
	err = f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(input))
	if err != nil {
		return err
	}

	return nil
}
