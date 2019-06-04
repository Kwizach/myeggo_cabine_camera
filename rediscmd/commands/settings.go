package commands

import (
	"rpi-client/rediscmd/redis"
)

// allSettings
var allSettings = make(map[string]string)

func init() {
	allSettings["shoot_folder"] = "/tmp/shoot"
	allSettings["photo_ext"] = ".jpg"

	// NEW_SETTING
	AllCommands["NEW_SETTING"] = func(args []string) error {
		if len(args) > 0 {
			key := args[0]
			val, err := redis.GetMyKey(key)
			if err != nil {
				return err
			}
			if val != "" {
				allSettings[key] = val
			}
		} else {
			redis.Log("Received NEW_SETTING with no argument", nil)
		}

		return nil
	}
}