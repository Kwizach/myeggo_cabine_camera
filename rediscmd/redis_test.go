package rediscmd

import (
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
	time.Sleep(1)

	if myID != str {
		t.Errorf("SubRedis should set myID")
	}
}
