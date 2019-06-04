package rediscmd

import (
	"os"
	"rpi-client/rediscmd/commands"
	"rpi-client/rediscmd/redis"
	"strings"
)

// AllSettings are variable that will be use through out the program
var AllSettings = make(map[string]string)

func init() {
	modifySettingsFromEnv()

	for k, v := range commands.AllSettings {
		AllSettings[k] = v
		redis.AllSettings[k] = v
	}

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

func modifySettingsFromEnv() {
	envToCheck := []string{"redis_server", "channelIN", "channelOUT", "channelLOG", "ntpd_server"}

	for _, v := range envToCheck {
		envValue := os.Getenv("EGG_" + strings.ToUpper(v))
		if envValue != "" {
			commands.AllSettings[v] = envValue
		}
	}
}
