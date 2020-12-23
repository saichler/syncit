package tests

import (
	"fmt"
	"github.com/saichler/syncit/cmd"
	"github.com/saichler/syncit/model"
	"github.com/saichler/syncit/transport"
	"testing"
	"time"
)

func TestCommandHandler(t *testing.T) {
	key := "qHYsJuloczNsFrbqlhlffjkRuHWfrCtH"
	key2 := "qHYsJuloczNsFrbqlhlffjkRuHWfrCtH"
	ch := &cmd.CommandHandler{}
	server := transport.NewListener(45621, "Hello World", key, ch)
	go server.Listen()

	time.Sleep(time.Second)

	service, err := transport.Connect("127.0.0.1", key2, "Hello World", 45621, ch)
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}

	c := &model.Command{}
	c.Cli = "ls"
	c.Args = []string{"/home/saichler"}
	err = transport.Send(c, service)
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	time.Sleep(time.Second * 10)
}
