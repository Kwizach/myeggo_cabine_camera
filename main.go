package main

import (
	"fmt"
	"rpi-client/rediscmd"
	"rpi-client/system"
)

func main() {
	// Check if the RPI has default hostname
	if !system.IsGoodHostName() {
		if err := system.CreateNewHostName(); err != nil {
			fmt.Println(err)
			return
		}
	}

	rediscmd.SubRedis(system.MyID())
}
