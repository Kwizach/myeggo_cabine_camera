package system

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var cpuID string

// MyID returns the RPI ID
func MyID() string {
	return "rpi-" + cpuID
}

// IsHostNameGood Check if the name of the RPI is the default name
func IsHostNameGood() bool {
	if name, err := os.Hostname(); err == nil {
		goodHostname, err := getGoodName()
		if err != nil {
			return false
		}
		return name == goodHostname
	}
	return false
}

// CreateNewHostName Set the a new random name, based on a uuid to the RPI
func CreateNewHostName() error {
	newName, err := getGoodName()
	if err != nil {
		return err
	}
	if newName == "rpi-" {
		return errors.New("[id.CreateNewHostName] Can't retrieve CPU serial")
	}

	return changeHostName(newName)
}

func getGoodName() (string, error) {
	_, err := getCPUSerial()
	if err != nil {
		return "", err
	}
	return MyID(), nil
}

func getCPUSerial() ([]byte, error) {
	if cpuID != "" {
		return []byte(cpuID), nil
	}

	cpuInfo, err := exec.Command("sh", "-c", "cat /proc/cpuinfo | grep Serial | cut -d':' -f2 | tr -d ' ' | tr -d '\n'").Output()
	if err != nil {
		return nil, err
	}
	cpuID = string(cpuInfo)

	return cpuInfo, nil
}

func changeHostName(newName string) error {
	if err := changeEtcHosts(newName); err != nil {
		return err
	}

	if err := changeEtcHostname(newName); err != nil {
		return err
	}

	if err := changeHostNameCtl(newName); err != nil {
		return err
	}

	if err := restartAvahiDaemon(); err != nil {
		return err
	}

	return nil
}

func changeEtcHosts(name string) error {
	return exec.Command("sh", "-c", fmt.Sprintf("sed -i -e 's/raspberrypi/%s/g' /etc/hosts", name)).Run()
}
func changeEtcHostname(name string) error {
	return exec.Command("sh", "-c", fmt.Sprintf("echo \"%s\" > /etc/hostname", name)).Run()
}
func changeHostNameCtl(name string) error {
	return exec.Command("sh", "-c", fmt.Sprintf("hostnamectl set-hostname \"%s\"", name)).Run()
}

func restartAvahiDaemon() error {
	// To be able to use mynewhostname.local from other machines,
	// we need to restart the mDNS daemon to respond to the new hostname.
	return exec.Command("sh", "-c", "systemctl restart avahi-daemon").Run()
}
