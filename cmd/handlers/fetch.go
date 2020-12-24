package handlers

import "C"
import (
	"fmt"
	"github.com/saichler/syncit/model"
	"github.com/saichler/syncit/transport"
	log "github.com/saichler/utils/golang"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

const (
	MAX_PART_SIZE = 5 * 1024 * 1024
)

type Fetch struct {
	sy  *sync.Mutex
	mtx map[string]*sync.Cond
}

func (h *Fetch) Cli() string {
	return "fetch"
}

func (h *Fetch) HandleCommand(c *model.Command, tc *transport.Connection) {
	f, err := os.Stat(c.Args[0])
	var file *os.File
	defer func() {
		if err != nil {
			log.Error(err)
		}
		err := transport.Send(c, tc)
		if err != nil {
			log.Error("err:", err)
		}
		if file != nil {
			file.Close()
		}
	}()
	if err != nil {
		c.Response = []byte(err.Error())
		return
	}
	if f.Size() > MAX_PART_SIZE {
		file, err = os.Open(c.Args[0])
		if err != nil {
			c.Response = []byte(err.Error())
			return
		}
		c.ResponseCount = int32(f.Size()/MAX_PART_SIZE + 1)
		for c.ResponseId < c.ResponseCount-1 {
			data := make([]byte, MAX_PART_SIZE)
			file.Read(data)
			c.Response = data
			transport.Send(c, tc)
			c.ResponseId++
		}
		data := make([]byte, f.Size()-int64(MAX_PART_SIZE*c.ResponseId))
		file.Read(data)
		c.Response = data
	} else {
		data, err := ioutil.ReadFile(c.Args[0])
		if err != nil {
			c.Response = []byte(err.Error())
		} else {
			c.Response = data
		}
	}
}

func (h *Fetch) HandleResponse(c *model.Command, tc *transport.Connection) {
	var file *os.File
	var err error
	index := strings.LastIndex(c.Args[1], "/")
	if index != -1 {
		dirPath := c.Args[1][0:index]
		_, exist := os.Stat(dirPath)
		if exist != nil {
			os.MkdirAll(dirPath, 0777)
		}
	}
	defer func() {
		if c.ResponseId == 0 {
			fmt.Print("Receiving ", c.Args[1], ".")
		}
		if file != nil {
			file.Close()
		}
		if err != nil {
			log.Error(err)
		}
		if c.ResponseId == c.ResponseCount-1 || c.ResponseCount == 0 {
			fmt.Println("Done!")
			h.sy.Lock()
			cond := h.mtx[c.Args[0]]
			h.sy.Unlock()
			cond.Broadcast()
		} else {
			fmt.Print(".")
		}
	}()
	if c.ResponseCount == 0 || c.ResponseId == 0 {
		file, err = os.Create(c.Args[1])
		if err != nil {
			return
		} else {
			file.Write(c.Response)
		}
	} else {
		file, err := os.OpenFile(c.Args[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			return
		}
		file.Write(c.Response)
	}
}

func (h *Fetch) Exec(args []string, tc *transport.Connection) {
	if args == nil || len(args) != 2 {
		log.Error("Fetch requiers source and destination")
		return
	}
	if h.sy == nil {
		h.sy = &sync.Mutex{}
		h.mtx = make(map[string]*sync.Cond)
	}
	c := &model.Command{}
	c.Cli = h.Cli()
	c.Args = args
	cond := sync.NewCond(&sync.Mutex{})
	h.sy.Lock()
	h.mtx[c.Args[0]] = cond
	h.sy.Unlock()
	h.mtx[c.Args[0]].L.Lock()
	transport.Send(c, tc)
	h.mtx[c.Args[0]].Wait()
	h.mtx[c.Args[0]].L.Unlock()
	h.sy.Lock()
	delete(h.mtx, c.Args[0])
	h.sy.Unlock()
}
