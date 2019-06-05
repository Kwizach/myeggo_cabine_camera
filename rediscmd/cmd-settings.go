package rediscmd

import (
	"os"
	"strings"
)

// allSettings are variable that will be use through out the program
var allSettings = make(map[string]string)

func init() {
	allSettings["redis_server"] = "redis://redis.egg:6379/"
	allSettings["channelIN"] = "commands"
	allSettings["channelOUT"] = "out"
	allSettings["channelLOG"] = "log"

	allSettings["ntpd_server"] = "ntpd.egg"

	allSettings["shoot_folder"] = "/tmp/shoot"
	allSettings["photo_ext"] = ".jpg"

	// NEW_SETTING
	allCommands["NEW_SETTING"] = func(args []string) error {
		if len(args) > 0 {
			key := args[0]
			val, err := getMyKey(key)
			if err != nil {
				return err
			}
			if val != "" {
				allSettings[key] = val
			}
		} else {
			Log("Received NEW_SETTING with no argument", nil)
		}

		return nil
	}

	modifySettingsFromEnv()
}

func modifySettingsFromEnv() {
	envToCheck := []string{"redis_server", "channelIN", "channelOUT", "channelLOG", "ntpd_server"}

	for _, v := range envToCheck {
		envVal := os.Getenv("EGG_" + strings.ToUpper(v))
		if envVal != "" {
			allSettings[v] = envVal
		}
	}
}
