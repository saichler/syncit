package handlers

import (
	"github.com/golang/protobuf/proto"
	"github.com/saichler/syncit/model"
	"github.com/saichler/syncit/transport"
)

type CommandHandler interface {
	Execute(command string, args []string, tc *transport.Connection)
}

type Handler interface {
	Cli() string
	HandleCommand(command *model.Command, tc *transport.Connection)
	HandleResponse(command *model.Command, tc *transport.Connection)
	Exec([]string, *transport.Connection)
}

func SetResponse(c *model.Command, pb proto.Message) error {
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	c.Response = data
	return nil
}
