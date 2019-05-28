package rediscmd

import (
	"errors"
	"fmt"
	"os/exec"
)

func init() {
	allCommands["REBOOT"] = func(_ []string) error { return reboot() }
	allCommands["SHUTDOWN"] = func(_ []string) error { return shutdown() }
	allCommands["LIGHT"] = func(strs []string) error {
		action := "default"
		if len(strs) > 0 {
			action = strs[0]
		}
		return light(action)
	}
}

// reboot the RPI
func reboot() error {
	cmd := exec.Command("reboot")
	if err := cmd.Run(); err != nil {
		return rpiMsgWithError("Error while trying to reboot", err)
	}
	// returning error will stop the program
	return errors.New("Rebooting RPI")
}

// shutdown the RPI
func shutdown() error {
	cmd := exec.Command("shutdown", "now")
	if err := cmd.Run(); err != nil {
		return rpiMsgWithError("Error while trying to shutdown", err)
	}
	// returning error will stop the program
	return errors.New("Shutting Down RPI")
}

// light turns the RPI led on, off or default
func light(action string) error {
	var cmd *exec.Cmd

	switch action {
	case "on":
		cmd = exec.Command("sh", "-c", "echo none > /sys/class/leds/led0/trigger; echo 0 > /sys/class/leds/led0/brightness")
	case "off":
		cmd = exec.Command("sh", "-c", "echo none > /sys/class/leds/led0/trigger; echo 1 > /sys/class/leds/led0/brightness")
	case "heartbeat":
		cmd = exec.Command("sh", "-c", "echo heartbeat > /sys/class/leds/led0/trigger")
	case "default":
		cmd = exec.Command("sh", "-c", "echo mmc0 > /sys/class/leds/led0/trigger")
	default:
		return rpiMsg(fmt.Sprintf("Unknown parameter %s for LIGHT", action))
	}

	if err := cmd.Run(); err != nil {
		return rpiMsgWithError(fmt.Sprintf("Error while trying to LIGHT %s :", action), err)
	}

	return nil
}
