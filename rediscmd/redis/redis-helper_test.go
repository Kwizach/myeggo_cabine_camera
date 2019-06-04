package redis

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestConnectRedis(t *testing.T) {
	service, err := connectRedis(AllSettings["redisURL"])
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
	service, _ := connectRedis(AllSettings["redisURL"])
	ch, errs := service.subscribe(AllSettings["channelIN"])

	service.publish(AllSettings["channelIN"], "STOP")
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

	service.unsubscribe(AllSettings["channelIN"])
}

func TestSubAndManage(t *testing.T) {
	service, _ := connectRedis(AllSettings["redisURL"])

	errorOK := errors.New("STOP Listening")

	onMsg := func(channel, message string) error {
		if message == "STOP" {
			return errorOK
		}
		return errors.New("Received something else")
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		service.publish(AllSettings["channelIN"], "STOP")
	}()

	err := service.subAndManage(onMsg, AllSettings["channelIN"])

	if err != errorOK {
		t.Errorf("Payload should be STOP")
	}
}

func TestUnsubscribe(t *testing.T) {
	service, _ := connectRedis(AllSettings["redisURL"])

	err := service.unsubscribe(AllSettings["channelIN"])
	if fmt.Sprintf("%s", err) != "There is no Subscription to unsubscribe from" {
		t.Errorf("Unsubscribe should not be able to unsubscribe when there is no subscription")
	}

	service.subscribe(AllSettings["channelIN"])
	if err = service.unsubscribe(AllSettings["channelIN"]); err != nil {
		t.Errorf("Unsubscribe should be able to unsubscribe")
	}
}
