package rediscmd

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestConnectRedis(t *testing.T) {
	service, err := ConnectRedis(redisURL)
	t.Run("group", func(t *testing.T) {
		t.Run("Test1", func(t *testing.T) {
			if err != nil {
				t.Errorf("ConnectRedis should return without error")
			}
		})
		t.Run("Test2", func(t *testing.T) {
			if service.client == nil {
				t.Errorf("ConnectRedis should return service with not nil client")
			}
		})
	})
}

func TestSubscribe(t *testing.T) {
	service, _ := ConnectRedis(redisURL)
	ch, errs := service.subscribe(channelIN)

	service.Publish(channelIN, "STOP")
	msg := <-ch

	t.Run("group", func(t *testing.T) {
		t.Run("Test1", func(t *testing.T) {
			if errs != nil {
				t.Errorf("subscribe should return without error")
			}
		})
		t.Run("Test2", func(t *testing.T) {
			if msg.Payload != "STOP" {
				t.Errorf("Payload should be STOP")
			}
		})
	})

	service.Unsubscribe(channelIN)
}

func TestSubAndManage(t *testing.T) {
	service, _ := ConnectRedis(redisURL)

	errorOK := errors.New("STOP Listening")

	onMsg := func(channel, message string) error {
		if message == "STOP" {
			return errorOK
		}
		return errors.New("Received something else")
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		service.Publish(channelIN, "STOP")
	}()

	err := service.SubAndManage(onMsg, channelIN)

	if err != errorOK {
		t.Errorf("Payload should be STOP")
	}
}

func TestUnsubscribe(t *testing.T) {
	service, _ := ConnectRedis(redisURL)

	err := service.Unsubscribe(channelIN)
	if fmt.Sprintf("%s", err) != "There is no Subscription to unsubscribe from" {
		t.Errorf("Unsubscribe should not be able to unsubscribe when there is no subscription")
	}

	service.subscribe(channelIN)
	if err = service.Unsubscribe(channelIN); err != nil {
		t.Errorf("Unsubscribe should be able to unsubscribe")
	}
}
