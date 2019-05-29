package system

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const defaultHostname = "raspberrypi"

var myID string

// MyID returns the RPI ID
func MyID() string {
	return myID
}

// IsDefaultHostName Check if the name of the RPI is the default name
func IsDefaultHostName() bool {
	if name, err := os.Hostname(); err == nil {
		return name == defaultHostname
	}
	return true
}

// CreateNewHostName Set the a new random name, based on a uuid to the RPI
func CreateNewHostName() error {
	newName, err := getCPUSerial()
	if err != nil {
		return err
	}
	if newName == "rpi-" {
		return errors.New("Can't retrieve CPU serial")
	}

	return changeHostName(newName)
}

func getCPUSerial() (string, error) {
	cpuID, err := exec.Command("sh", "-c", "cat /proc/cpuinfo | grep Serial | cut -d':' -f2 | tr -d ' ' | tr -d '\n'").Output()

	if err != nil {
		fmt.Println("Error in getting cpuID")
		return "", err
	}

	return fmt.Sprintf("rpi-%s", cpuID), nil
}

func changeHostName(newName string) error {

	if err := changeEtcHosts(newName); err != nil {
		fmt.Println("Error in changeEtcHosts")
		return err
	}

	if err := changeEtcHostname(newName); err != nil {
		fmt.Println("Error in changeEtcHostname")
		return err
	}

	if err := changeHostNameCtl(newName); err != nil {
		fmt.Println("Error in changeHostNameCtl")
		return err
	}

	myID = newName

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

// 4) Restart the mDNS daemon
// To be able to use mynewhostname.local from other machines, we need to restart the mDNS daemon to respond to the new hostname.
// sudo systemctl restart avahi-daemon
