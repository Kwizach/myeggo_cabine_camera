package rediscmd

import (
	"errors"
)

var (
	// ErrorSTOP is returned when we receive STOP on channelIN
	ErrorSTOP = errors.New("STOP Listening")
)

func init() {
	// STOP
	allCommands["STOP"] = func(_ []string) error {
		return ErrorSTOP
	}
	// PING
	allCommands["PING"] = func(_ []string) error {
		return rpiMsg("PONG")
	}
}
