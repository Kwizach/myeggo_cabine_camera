package rediscmd

import (
	"rpi-client/rediscmd/redis"
	"testing"
)

func TestInit(t *testing.T) {
	if _, ok := redis.AllCommands["PING"]; !ok {
		t.Errorf("STOP command should exist")
	}

	if _, ok := redis.AllCommands["PINGervkbnev"]; ok {
		t.Errorf("weird command shouldn't exist")
	}
}

// Other functions have been tested in redis package
