package tests

import (
	"fmt"
	"github.com/saichler/syncit/transport"
	"testing"
	"time"
)

func TestTrsansport(t *testing.T) {
	key := "qHYsJuloczNsFrbqlhlffjkRuHWfrCtH"
	key2 := "qHYsJuloczNsFrbqlhlffjkRuHWfrCtH"
	server := transport.NewListener(45621, "Hello World", key, nil)
	go server.Listen()

	time.Sleep(time.Second)

	service, err := transport.Connect("127.0.0.1", key2, "Hello World", 45621, nil)
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}

	err = service.Send([]byte("Welcome"))
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	time.Sleep(time.Second * 10)
}
