package handlers

import "C"
import (
	"github.com/saichler/syncit/model"
	"github.com/saichler/syncit/transport"
	log "github.com/saichler/utils/golang"
	"io/ioutil"
	"os"
)

type Fetch struct {
}

func (h *Fetch) Cli() string {
	return "fetch"
}

func (h *Fetch) HandleCommand(c *model.Command, tc *transport.Connection) {
	data, err := ioutil.ReadFile(c.Args[0])
	if err != nil {
		c.Response = []byte(err.Error())
	} else {
		c.Response = data
	}
	transport.Send(c, tc)
}

func (h *Fetch) HandleResponse(c *model.Command, tc *transport.Connection) {
	file, err := os.Create(c.Args[1])
	if err != nil {
		log.Error(err)
	} else {
		file.Write(c.Response)
		file.Close()
	}
	log.Info("Written: ", c.Args[1])
}

func (h *Fetch) Exec(args []string, tc *transport.Connection) {
	source := "/home/saichler/demo.zip"
	target := "/tmp/demo.zip"
	c := &model.Command{}
	c.Cli = h.Cli()
	c.Args = []string{source, target}
	transport.Send(c, tc)
}
