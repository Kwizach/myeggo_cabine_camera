package rediscmd

import (
	"errors"
	"fmt"
	"strings"
)

var (
	myID string
	serv *service
	// allCommands is a hash table with commands coming from channelIN
	// and their associated function
	allCommands = make(map[string]func(_ []string) error)
)

// SubRedis subscribe to Redis pubsub and manage incoming messages
func SubRedis(rpiID string) error {
	myID = rpiID

	var err error
	// Connect to Redis
	serv, err = connectRedis(allSettings["redis_server"])
	if err != nil {
		return err
	}
	// Tell the server that we are listening
	rpiMsg("Connected")

	// Subscribe to channelIN and wait for messages
	// This function won't exit until there is an error
	err = serv.subAndManage(onMsg, allSettings["channelIN"])
	Log("subAndManage ended with error", err)

	// Tell the server that we are not listening anymore
	rpiMsg("Disconnected")

	var errorMsg string

	// We are done... unsubscribe
	errU := serv.unsubscribe()
	if errU != nil {
		errorMsg = fmt.Sprintf("subAndManage ended with error: %s\nAnd could not unsubscribe with error: %s", err, errU)
	} else {
		errorMsg = fmt.Sprintf("subAndManage ended with error: %s\n", err)
	}

	return errors.New(errorMsg)
}

func onMsg(channel string, message string) error {
	if channel == allSettings["channelIN"] {
		msgs := strings.Split(message, " ")
		if msgs[0] != "" {
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
		return rpiMsg("Received empty command")
	}
	return nil
}

// rpiMsg publish a message on channelOUT
func rpiMsg(msg string) error {
	if serv != nil {
		return serv.publish(allSettings["channelOUT"], fmt.Sprintf("RPI [%s] - %s", myID, msg))
	}
	return errors.New("Service is down")
}

// Log send message and error on channelLOG
func Log(msg string, err error) error {
	if serv != nil {
		if err != nil {
			return serv.publish(allSettings["channelLOG"], fmt.Sprintf("RPI [%s] - %s %s", myID, msg, err))
		}
		return serv.publish(allSettings["channelLOG"], fmt.Sprintf("RPI [%s] - %s", myID, msg))
	}
	return errors.New("Service is down")
}

// getMyKey retrieve a key from Redis
// 1st check if we have a dedicated one first
// if not get the global one
func getMyKey(key string) (string, error) {
	dedicatedKey := fmt.Sprintf("%s:%s", myID, key)
	res, err := serv.getKey(dedicatedKey)
	if res == "" { // dedicatedKey doesn't exist in Redis
		res, err = serv.getKey(key)
		if res == "" && err == nil { // key doesn't exist
			return "", nil
		} else if err != nil {
			return "", err
		}
	}
	return res, nil
}
