package rediscmd

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestConnectRedis(t *testing.T) {
	fmt.Println(allSettings)
	service, err := connectRedis(allSettings["redis_server"])
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
	service, _ := connectRedis(allSettings["redis_server"])
	ch, errs := service.subscribe(allSettings["channelIN"])

	service.publish(allSettings["channelIN"], "STOP")
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

	service.unsubscribe(allSettings["channelIN"])
}

func TestSubAndManage(t *testing.T) {
	service, _ := connectRedis(allSettings["redis_server"])

	errorOK := errors.New("STOP Listening")

	onMsg := func(channel, message string) error {
		if message == "STOP" {
			return errorOK
		}
		return errors.New("Received something else")
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		service.publish(allSettings["channelIN"], "STOP")
	}()

	err := service.subAndManage(onMsg, allSettings["channelIN"])

	if err != errorOK {
		t.Errorf("Payload should be STOP")
	}
}

func TestUnsubscribe(t *testing.T) {
	service, _ := connectRedis(allSettings["redis_server"])

	err := service.unsubscribe(allSettings["channelIN"])
	if fmt.Sprintf("%s", err) != "There is no Subscription to unsubscribe from" {
		t.Errorf("Unsubscribe should not be able to unsubscribe when there is no subscription")
	}

	service.subscribe(allSettings["channelIN"])
	if err = service.unsubscribe(allSettings["channelIN"]); err != nil {
		t.Errorf("Unsubscribe should be able to unsubscribe")
	}
}
