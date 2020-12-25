package cmd

import (
	"github.com/golang/protobuf/proto"
	"github.com/saichler/syncit/cmd/handlers"
	"github.com/saichler/syncit/model"
	"github.com/saichler/syncit/transport"
	log "github.com/saichler/utils/golang"
)

var cmdHandlers = make(map[string]handlers.Handler)
var ini = initHandlers()

type CommandHandler struct {
}

func initHandlers() bool {
	ls := &handlers.LS{}
	cmdHandlers[ls.Cli()] = ls
	fetch := handlers.NewFetch()
	cmdHandlers[fetch.Cli()] = fetch
	sync := handlers.NewSync(&CommandHandler{})
	cmdHandlers[sync.Cli()] = sync
	return true
}

func (ch *CommandHandler) Execute(command string, args []string, tc *transport.Connection) {
	h, ok := cmdHandlers[command]
	if ok {
		h.Exec(args, tc)
	} else {
		log.Error("Unknown command ", command)
	}
}

func (ch *CommandHandler) HandleMessage(msg []byte, tc *transport.Connection) {
	cmd := &model.Command{}
	err := proto.Unmarshal(msg, cmd)
	if err != nil {
		log.Error("Failed to unmarshal command:", err)
		return
	}
	ch.handleCommand(cmd, tc)
}

func (ch *CommandHandler) handleCommand(c *model.Command, tc *transport.Connection) {
	if c.Response == nil {
		h, ok := cmdHandlers[c.Cli]
		if ok {
			h.HandleCommand(c, tc)
			return
		}
		c.Response = []byte("my response")
		transport.Send(c, tc)
	} else {
		h, ok := cmdHandlers[c.Cli]
		if ok {
			h.HandleResponse(c, tc)
			return
		}
		log.Error(string(c.Response))
	}
}
