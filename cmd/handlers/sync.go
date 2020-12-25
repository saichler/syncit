package handlers

import (
	"github.com/golang/protobuf/proto"
	"github.com/saichler/syncit/files"
	"github.com/saichler/syncit/model"
	"github.com/saichler/syncit/transport"
	log "github.com/saichler/utils/golang"
)

type Sync struct {
	commandHandler CommandHandler
}

func NewSync(commadHandler CommandHandler) *Sync {
	sy := &Sync{}
	sy.commandHandler = commadHandler
	return sy
}

func (h *Sync) Cli() string {
	return "sync"
}

func (h *Sync) HandleCommand(c *model.Command, tc *transport.Connection) {
	log.Info("Scanning requested directory:", c.Args[0])
	dir := files.Scan(c.Args[0])
	err := SetResponse(c, dir)
	if err != nil {
		c.Response = []byte(err.Error())
	}
	transport.Send(c, tc)
}

func (h *Sync) HandleResponse(c *model.Command, tc *transport.Connection) {
	dir := &model.File{}
	err := proto.Unmarshal(c.Response, dir)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Received Requested library:", c.Args[0], " and comparing to local copy:", c.Args[1])
	files.Stat(dir, c.Args[1])
	log.Info("Fetching missing Files...")
	h.fetch(dir, tc)
	log.Info("Done Synching Directories")
}

var count = 0

func (h *Sync) fetch(file *model.File, tc *transport.Connection) {
	if file.Files == nil {
		if file.SizeA != file.SizeZ {
			h.commandHandler.Execute("fetch", []string{file.NameA, file.NameZ}, tc)
		}
	} else {
		for _, subFile := range file.Files {
			h.fetch(subFile, tc)
		}
	}
}

func (h *Sync) Exec(args []string, tc *transport.Connection) {
	if args == nil || len(args) != 2 {
		log.Error("Sync need 2 args, source & destination")
		return
	}
	c := &model.Command{}
	c.Cli = h.Cli()
	c.Args = args
	transport.Send(c, tc)
}
