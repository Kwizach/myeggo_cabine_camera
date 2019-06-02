package main

import (
	"rpi-client/rediscmd"
	"rpi-client/system"
)

func main() {
	// Check if the RPI has hostname set has expected
	if !system.IsHostNameGood() {
		if err := system.CreateNewHostName(); err != nil {
			log(err)
			return
		}
	}

	// Check if the RPI has MAC Address set has expected
	if !system.IsCurrentMACGood() {
		if err := system.SetMACAddress(); err != nil {
			log(err)
			return
		}
	}

	// Launch Redis PubSub channels:
	//		* 'commands' channel from the server
	//		* 'out' channel to the server
	//		* 'log' channel to the server
	log(rediscmd.SubRedis(system.MyID()))
}
