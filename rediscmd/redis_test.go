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
	if serv != nil {
		serv.publish(allSettings["channelIN"], "STOP")
		serv.unsubscribe(allSettings["channelIN"])
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

	chTIN, errIN := serv.subscribe(allSettings["channelIN"])
	if errIN != nil {
		t.Errorf("TestOnMsg can't subscribe")
	}
	chTOUT, errOUT := serv.subscribe(allSettings["channelOUT"])
	if errOUT != nil {
		t.Errorf("TestOnMsg can't subscribe")
	}
	defer serv.unsubscribe(allSettings["channelIN"], allSettings["channelOUT"])

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
						wait4OUT = "RPI [TEST] - Unknown command: TEST_UNKNOWN"
					} else {
						wait4OUT = "RPI [TEST] - Received empty command"
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
			serv.publish(allSettings["channelIN"], "TEST1")
		})
		t.Run("Test2", func(t *testing.T) {
			serv.publish(allSettings["channelIN"], "TEST2 is successfull")
		})
		t.Run("Test3", func(t *testing.T) {
			serv.publish(allSettings["channelIN"], "TEST1 with 4 useless parameters")
		})
		t.Run("Test4", func(t *testing.T) {
			serv.publish(allSettings["channelIN"], "TEST_UNKNOWN")
		})
		t.Run("Test5", func(t *testing.T) {
			time.Sleep(555 * time.Millisecond)
			serv.publish(allSettings["channelIN"], "")
		})
	})
}

func TestRpiMsg(t *testing.T) {
	chTOUT, errOUT := serv.subscribe(allSettings["channelOUT"])
	if errOUT != nil {
		t.Errorf("TestRpiMsg can't subscribe")
	}
	defer serv.unsubscribe(allSettings["channelOUT"])

	time.Sleep(2019 * time.Millisecond)
	rpiMsg("TEST_MSG")
	msgI := <-chTOUT
	if msgI.Payload != "RPI [TEST] - TEST_MSG" {
		t.Errorf("rpiMsg failed - received: %s", msgI.Payload)
	}
}

func TestLog(t *testing.T) {
	chTLOG, errLOG := serv.subscribe(allSettings["channelLOG"])
	if errLOG != nil {
		t.Errorf("TestLog can't subscribe")
	}
	defer serv.unsubscribe(allSettings["channelLOG"])

	time.Sleep(2019 * time.Millisecond)
	Log("TEST_MSG", errors.New("WITH ERROR"))
	msgI := <-chTLOG
	if msgI.Payload != "RPI [TEST] - TEST_MSG WITH ERROR" {
		t.Errorf("TestLog failed - received: %s", msgI.Payload)
	}
}

func TestGetMyKey(t *testing.T) {
	myKey := "clef:test"
	globalKey := "clef:test2"

	serv.setKeyValue(myID+":"+myKey, "42", 10*time.Second)
	serv.setKeyValue(globalKey, "84", 10*time.Second)
	time.Sleep(2019 * time.Millisecond)

	if val, err := getMyKey(myKey); err != nil || val != "42" {
		t.Errorf("TEST1 failed expected 42 received %s\n error: %s", val, err)
	}
	if val, err := getMyKey(globalKey); err != nil || val != "84" {
		t.Errorf("TEST2 failed expected 84 received %s\n error: %s", val, err)
	}
	if val, err := getMyKey("unknown_key"); err != nil || val != "" {
		t.Errorf("TEST3 failed expected '' received %s\n error: %s", val, err)
	}
}
