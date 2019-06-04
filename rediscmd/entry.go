package rediscmd

import (
	"os"
	"rpi-client/rediscmd/commands"
	"rpi-client/rediscmd/redis"
	"strings"
)

func init() {
	// link redis settings to commands package
	redis.AllSettings = commands.AllSettings
	// link redis commands to commands package
	redis.AllCommands = commands.AllCommands

	modifySettingsFromEnv()
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
