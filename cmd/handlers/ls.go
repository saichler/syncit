package handlers

import (
	"github.com/golang/protobuf/proto"
	"github.com/saichler/syncit/files"
	"github.com/saichler/syncit/model"
	"github.com/saichler/syncit/transport"
	log "github.com/saichler/utils/golang"
)

type LS struct {
}

func (ls *LS) Cli() string {
	return "ls"
}

func (ls *LS) HandleCommand(c *model.Command, tc *transport.Connection) {
	dir := files.Scan(c.Args[0])
	err := SetResponse(c, dir)
	if err != nil {
		c.Response = []byte(err.Error())
	}
	transport.Send(c, tc)
}

func (ls *LS) HandleResponse(c *model.Command, tc *transport.Connection) {
	dir := &model.File{}
	err := proto.Unmarshal(c.Response, dir)
	if err != nil {
		log.Error(err)
		return
	}
	files.Print(dir, 2, true, true)
}

func (ls *LS) Exec(args []string, tc *transport.Connection) {
	c:=&model.Command{}
	c.Cli = ls.Cli()
	c.Args = args
	transport.Send(c,tc)
}
