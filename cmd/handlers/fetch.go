package handlers

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
	mtx  *sync.Mutex
	jobs map[string]*FetchJob
}

func NewFetch() *Fetch {
	f := &Fetch{}
	f.mtx = &sync.Mutex{}
	f.jobs = make(map[string]*FetchJob)
	return f
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
		c.ResponseCount = int32(f.Size() / MAX_PART_SIZE)
		if f.Size()&MAX_PART_SIZE > 0 {
			c.ResponseCount++
		}
		for c.ResponseId < c.ResponseCount-1 {
			data := make([]byte, MAX_PART_SIZE)
			file.Read(data)
			c.Response = data
			transport.Send(c, tc)
			c.ResponseId++
		}
		left := f.Size() - int64(MAX_PART_SIZE*c.ResponseId)
		if left > 0 {
			data := make([]byte, left)
			file.Read(data)
			c.Response = data
		}
	} else {
		data, err := ioutil.ReadFile(c.Args[0])
		if err != nil {
			c.Response = []byte(err.Error())
		} else {
			c.Response = data
		}
	}
}

func (h *Fetch) HandleResponse(command *model.Command, tc *transport.Connection) {
	h.mtx.Lock()
	fetchJob := h.jobs[command.Args[0]]
	h.mtx.Unlock()

	fetchJob.cond.L.Lock()
	if command.ResponseId != fetchJob.last+1 {
		fetchJob.waiting[command.ResponseId] = command
		fetchJob.cond.L.Unlock()
		return
	}

	var file *os.File
	var err error

	index := strings.LastIndex(command.Args[1], "/")
	if index != -1 {
		dirPath := command.Args[1][0:index]
		_, exist := os.Stat(dirPath)
		if exist != nil {
			os.MkdirAll(dirPath, 0777)
		}
	}

	defer func() {
		if command.ResponseId == 0 {
			fmt.Print("Receiving ", command.Args[1], " with ", command.ResponseCount, " parts:.")
		}
		if file != nil {
			file.Close()
		}
		if err != nil {
			log.Error(err)
		}
		if fetchJob.last == command.ResponseCount-1 || command.ResponseCount == 0 {
			fmt.Println("Done!")
			fetchJob.cond.Broadcast()
		} else {
			fmt.Print(".")
		}
		fetchJob.cond.L.Unlock()
	}()
	if command.ResponseCount == 0 || command.ResponseId == 0 {
		file, err = os.Create(command.Args[1])
		if err != nil {
			return
		} else {
			file.Write(command.Response)
		}
	} else {
		file, err := os.OpenFile(command.Args[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			log.Error(err)
			return
		}
		file.Write(command.Response)
	}
	fetchJob.last = command.ResponseId

	waitingCommand, ok := fetchJob.waiting[fetchJob.last+1]
	for ok {
		fetchJob.hadOrderIssue = true
		file.Write(waitingCommand.Response)
		fetchJob.last = waitingCommand.ResponseId
		delete(fetchJob.waiting, waitingCommand.ResponseId)
		waitingCommand, ok = fetchJob.waiting[fetchJob.last+1]
	}
}

func (h *Fetch) Exec(args []string, tc *transport.Connection) {
	if args == nil || len(args) != 2 {
		log.Error("Fetch requiers source and destination")
		return
	}

	c := &model.Command{}
	c.Cli = h.Cli()
	c.Args = args
	h.mtx.Lock()
	h.jobs[c.Args[0]] = &FetchJob{}
	h.jobs[c.Args[0]].waiting = make(map[int32]*model.Command, 0)
	h.jobs[c.Args[0]].cond = sync.NewCond(&sync.Mutex{})
	h.jobs[c.Args[0]].last = -1
	h.mtx.Unlock()

	h.jobs[c.Args[0]].cond.L.Lock()
	transport.Send(c, tc)
	h.jobs[c.Args[0]].cond.Wait()
	h.jobs[c.Args[0]].cond.L.Unlock()

	h.mtx.Lock()
	if h.jobs[c.Args[0]].hadOrderIssue {
		panic(c.Args[1] + " had order issues, please check!")
	}
	delete(h.jobs, c.Args[0])
	h.mtx.Unlock()
}
