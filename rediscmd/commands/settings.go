package commands

import (
	"rpi-client/rediscmd/redis"
)

// AllSettings are variable that will be use through out the program
var AllSettings = make(map[string]string)

func init() {
	AllSettings["redis_server"] = "redis://redis.egg:6379/"
	AllSettings["channelIN"] = "commands"
	AllSettings["channelOUT"] = "out"
	AllSettings["channelLOG"] = "log"

	AllSettings["ntpd_server"] = "ntpd.egg"

	AllSettings["shoot_folder"] = "/tmp/shoot"
	AllSettings["photo_ext"] = ".jpg"

	// NEW_SETTING
	AllCommands["NEW_SETTING"] = func(args []string) error {
		if len(args) > 0 {
			key := args[0]
			val, err := redis.GetMyKey(key)
			if err != nil {
				return err
			}
			if val != "" {
				AllSettings[key] = val
			}
		} else {
			redis.Log("Received NEW_SETTING with no argument", nil)
		}

		return nil
	}
}
