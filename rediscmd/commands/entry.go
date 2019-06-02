package commands

import (
	"errors"
	"rpi-client/rediscmd/redis"
)

// AllCommands that we will share with redis package
var AllCommands = make(map[string]func(_ []string) error)

func init() {
	// STOP
	AllCommands["STOP"] = func(_ []string) error {
		return errors.New("STOP Listening")
	}
	// PING
	AllCommands["PING"] = func(_ []string) error {
		return redis.RpiMsg("PONG")
	}
}
