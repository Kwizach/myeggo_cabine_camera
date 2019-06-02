package rediscmd

import (
	"rpi-client/rediscmd/commands"
	"rpi-client/rediscmd/redis"
)

func init() {
	for k, v := range commands.AllCommands {
		redis.AllCommands[k] = v
	}
}

// SubRedis export SubRedis from redis package
func SubRedis(rpiID string) error {
	return redis.SubRedis(rpiID)
}

// Log export Log from redis package
func Log(msg string, err error) error {
	return redis.Log(msg, err)
}
