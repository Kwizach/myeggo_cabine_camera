package rediscmd

import (
	"errors"
)

func init() {
	// STOP
	allCommands["STOP"] = func(_ []string) error {
		return errors.New("STOP Listening")
	}
	// PING
	allCommands["PING"] = func(_ []string) error {
		return rpiMsg("PONG")
	}
}
