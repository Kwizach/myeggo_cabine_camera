package redis

import (
	"errors"
	"fmt"
	"strings"
)

const (
	redisURL   string = "redis://redis.myeggo.com:6379/"
	channelIN  string = "commands"
	channelOUT string = "out"
	channelLOG string = "log"
)

var (
	myID    string
	service *Service
	// AllCommands is a hash table with commands coming from channelIN
	// and their associated function
	AllCommands = make(map[string]func(_ []string) error)
)

// SubRedis subscribe to Redis pubsub and manage incoming messages
func SubRedis(rpiID string) error {
	myID = rpiID

	var err error
	// Connect to Redis
	service, err = connectRedis(redisURL)
	if err != nil {
		return err
	}
	// Tell the server that we are listening
	RpiMsg("Connected")
	service.incrKey("connected-rpis")

	// Subscribe to channelIN and wait for messages
	// This function won't exit until there is an error
	err = service.subAndManage(onMsg, channelIN)
	Log("subAndManage ended with error", err)

	// Tell the server that we are not listening anymore
	RpiMsg("Disconnected")
	service.decrKey("connected-rpis")

	var errorMsg string

	// We are done... unsubscribe
	errU := service.unsubscribe()
	if errU != nil {
		errorMsg = fmt.Sprintf("subAndManage ended with error: %s\nAnd could not unsubscribe with error: %s", err, errU)
	} else {
		errorMsg = fmt.Sprintf("subAndManage ended with error: %s\n", err)
	}

	return errors.New(errorMsg)
}

func onMsg(channel string, message string) error {
	if channel == channelIN {
		msgs := strings.Split(message, " ")
		if msgs[0] != "" {
			if _, ok := AllCommands[msgs[0]]; ok {
				// If there are parameters
				if len(msgs[1:]) > 0 {
					return AllCommands[msgs[0]](msgs[1:])
				}
				// Command only, with no parameter
				return AllCommands[msgs[0]](nil)
			}
			// This command doesn't exist
			return RpiMsg(fmt.Sprintf("Unknown command: %s", message))
		}
		// Empty command
		return RpiMsg("Received empty command")
	}
	return nil
}

// RpiMsg publish a message on channelOUT
func RpiMsg(msg string) error {
	if service != nil {
		return service.publish(channelOUT, fmt.Sprintf("RPI [%s] - %s", myID, msg))
	}
	return errors.New("Service is down")
}

// Log send message and error on channelLOG
func Log(msg string, err error) error {
	if service != nil {
		return service.publish(channelLOG, fmt.Sprintf("RPI [%s] - %s %s", myID, msg, err))
	}
	return errors.New("Service is down")
}

// getMyKey retrieve a key from Redis
// 1st check if we have a dedicated one first
// if not get the global one
func getMyKey(key string) (string, error) {
	dedicatedKey := fmt.Sprintf("%s:%s", myID, key)
	res, err := service.getKey(dedicatedKey)
	if res == "" { // dedicatedKey doesn't exist in Redis
		res, err = service.getKey(key)
		if res == "" && err == nil { // key doesn't exist
			return "", nil
		} else if err != nil {
			return "", err
		}
	}
	return res, nil
}
