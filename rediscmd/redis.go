package rediscmd

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	redisURL   string = "redis://redis.myeggo.com:6379/"
	channelIN  string = "commands"
	channelOUT string = "out"
)

var (
	myID    string
	service *Service
	// allCommands is a hash table with commands and their associated function
	allCommands = make(map[string]func(_ []string) error)
)

func init() {
	// START
	allCommands["START"] = func(_ []string) error {
		when := fmt.Sprintf("%f", float64(time.Now().UnixNano())/1e9)
		fmt.Printf("%s - %s\n", when, "START")
		return service.Publish(channelOUT, when)
	}
	// STOP
	allCommands["STOP"] = func(_ []string) error {
		return errors.New("STOP Listening")
	}
	// PING
	allCommands["PING"] = func(_ []string) error {
		return rpiMsg("PONG")
	}
}

// SubRedis subscribe to Redis pubsub and manage incoming messages
func SubRedis(rpiID string) error {
	myID = rpiID

	var err error
	// Connect to Redis
	service, err = ConnectRedis(redisURL)
	if err != nil {
		return err
	}
	// Subscribe to channelIN and wait for messages
	// This function won't exit until there is an error
	err = service.SubAndManage(onMsg, channelIN)
	fmt.Println(err)
	// We are done... unsubscribe
	return service.Unsubscribe()
}

func onMsg(channel string, message string) error {
	if channel == channelIN {
		msgs := strings.Split(message, " ")
		if len(msgs) > 0 {
			if _, ok := allCommands[msgs[0]]; ok {
				// If there are parameters
				if len(msgs[1:]) > 0 {
					return allCommands[msgs[0]](msgs[1:])
				}
				// Command only, with no parameter
				return allCommands[msgs[0]](nil)
			}
			// This command doesn't exist
			return rpiMsg(fmt.Sprintf("Unknown command: %s", message))
		}
		// Empty command
		return rpiMsg(fmt.Sprintf("Received empty command: %s", message))
	}
	return nil
}

// rpiMsg publish a message on channelOUT
func rpiMsg(msg string) error {
	return service.Publish(channelOUT, fmt.Sprintf("RPI [%s] - %s", myID, msg))
}

// rpiMsgWithError publish a message and the error we have encountered on channelOUT
func rpiMsgWithError(msg string, err error) error {
	return service.Publish(channelOUT, fmt.Sprintf("RPI [%s] - %s %s", myID, msg, err))
}
