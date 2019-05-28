package main

import (
	"rpi-client/rediscmd"
	"rpi-client/system"
)

func main() {
	// Check if the RPI has default hostname
	if system.IsDefaultHostName() {
		system.CreateNewHostName()
	}

	rediscmd.SubRedis()
}
