package rediscmd

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestSubRedis(t *testing.T) {
	str := "TEST"

	go func() {
		err := SubRedis(str)
		if err != nil {
			t.Errorf("SubRedis returned with error")
		}
	}()
	time.Sleep(2019 * time.Millisecond)
	if service != nil {
		service.Unsubscribe(channelIN)
	}

	if myID != str {
		t.Errorf("SubRedis should set myID")
	}
}

func TestOnMsg(t *testing.T) {
	errTest1 := errors.New("SUCCESS TEST1")
	allCommands["TEST1"] = func(_ []string) error {
		return errTest1
	}
	allCommands["TEST2"] = func(params []string) error {
		return fmt.Errorf("Success TEST2 %s %s", params[0], params[1])
	}

	// var errC error
	// service, errC = ConnectRedis(redisURL)
	// if errC != nil {
	// 	t.Errorf("TestOnMsg can't ConnectRedis")
	// }

	chTIN, errIN := service.subscribe(channelIN)
	if errIN != nil {
		t.Errorf("TestOnMsg can't subscribe")
	}
	chTOUT, errOUT := service.subscribe(channelOUT)
	if errOUT != nil {
		t.Errorf("TestOnMsg can't subscribe")
	}
	defer service.Unsubscribe(channelIN, channelOUT)

	go func() {
		var (
			err      error
			wait4OUT string
		)
		for err == nil {
			select {
			case msgI := <-chTIN:
				switch msgI.Payload {
				case "TEST1", "TEST1 with 4 useless parameters":
					if onMsg(msgI.Channel, msgI.Payload) != errTest1 {
						t.Errorf("onMsg TEST1 failed")
						err = errors.New("TEST1 failed")
					}
				case "TEST2 is successfull":
					if fmt.Sprintf("%s", onMsg(msgI.Channel, msgI.Payload)) != "Success TEST2 is successfull" {
						t.Errorf("onMsg TEST2 failed")
						err = errors.New("TEST2 failed")
					}
				case "", "TEST_UNKNOWN":
					if msgI.Payload == "TEST_UNKNOWN" {
						wait4OUT = "RPI [] - Unknown command: TEST_UNKNOWN"
					} else {
						wait4OUT = "RPI [] - Received empty command"
					}
					errO := onMsg(msgI.Channel, msgI.Payload)
					if errO != nil {
						// check that publish is done
						// retrieve it with chTOUT
						t.Errorf("onMsg TEST_UNKNOWN Publish failed")
						err = errors.New("TEST4 failed")
					}
				}
			case msgO := <-chTOUT:
				if msgO.Payload != wait4OUT {
					t.Errorf("onMsg TEST_UNKNOWN Received failed")
					err = errors.New("TEST4 failed")
				}
			}
		}
	}()
	time.Sleep(2019 * time.Millisecond)

	t.Run("group", func(t *testing.T) {
		t.Run("Test1", func(t *testing.T) {
			service.Publish(channelIN, "TEST1")
		})
		t.Run("Test2", func(t *testing.T) {
			service.Publish(channelIN, "TEST2 is successfull")
		})
		t.Run("Test3", func(t *testing.T) {
			service.Publish(channelIN, "TEST1 with 4 useless parameters")
		})
		t.Run("Test4", func(t *testing.T) {
			service.Publish(channelIN, "TEST_UNKNOWN")
		})
		t.Run("Test5", func(t *testing.T) {
			time.Sleep(555 * time.Millisecond)
			service.Publish(channelIN, "")
		})
	})
}

func TestRpiMsg(t *testing.T) {
	// var err error
	// service, err = ConnectRedis(redisURL)
	// if err != nil {
	// 	t.Errorf("TestRpiMsg can't ConnectRedis")
	// }

	chTOUT, errOUT := service.subscribe(channelOUT)
	if errOUT != nil {
		t.Errorf("TestRpiMsg can't subscribe")
	}
	defer service.Unsubscribe(channelOUT)

	time.Sleep(2019 * time.Millisecond)
	rpiMsg("TEST_MSG")
	msgI := <-chTOUT
	if msgI.Payload != "RPI [] - TEST_MSG" {
		t.Errorf("rpiMsg failed")
	}
}

func TestRpiMsgWithError(t *testing.T) {
	// var err error
	// service, err = ConnectRedis(redisURL)
	// if err != nil {
	// 	t.Errorf("TestRpiMsgWithError can't ConnectRedis")
	// }

	chTOUT, errOUT := service.subscribe(channelOUT)
	if errOUT != nil {
		t.Errorf("TestRpiMsgWithError can't subscribe")
	}
	defer service.Unsubscribe(channelOUT)

	time.Sleep(2019 * time.Millisecond)
	rpiMsgWithError("TEST_MSG", errors.New("WITH ERROR"))
	msgI := <-chTOUT
	if msgI.Payload != "RPI [] - TEST_MSG WITH ERROR" {
		t.Errorf("rpiMsgWithError failed")
	}
}
