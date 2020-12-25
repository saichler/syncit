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
	"time"
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

func (h *Fetch) HandleResponse(c *model.Command, tc *transport.Connection) {
	h.mtx.Lock()
	fetchJob := h.jobs[c.Args[0]]
	h.mtx.Unlock()

	fetchJob.cond.L.Lock()
	if c.ResponseId != fetchJob.last+1 {
		fetchJob.queue = append(fetchJob.queue, c)
		fetchJob.cond.L.Unlock()
		log.Error("Not Last ", c.ResponseId, fetchJob.last)
		time.Sleep(time.Second)
		return
	}

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
			fmt.Print("Receiving ", c.Args[1], " with ", c.ResponseCount, " parts:.")
		}
		if file != nil {
			file.Close()
		}
		if err != nil {
			log.Error(err)
		}
		if fetchJob.last == c.ResponseCount-1 || c.ResponseCount == 0 {
			fmt.Println("Done!")
			fetchJob.cond.Broadcast()
		} else {
			fmt.Print(".")
		}
		fetchJob.cond.L.Unlock()
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
			log.Error(err)
			return
		}
		file.Write(c.Response)
	}
	fetchJob.last = c.ResponseId
	if len(fetchJob.queue) > 0 {
		fetchJob.hadOrderIssue = true
		found := true
		for found && len(fetchJob.queue) > 0 {
			found = false
			index := -1
			for i, c := range fetchJob.queue {
				if c.ResponseId == fetchJob.last+1 {
					found = true
					file.Write(c.Response)
					index = i
					fetchJob.last = c.ResponseId
					break
				}
			}
			if found {
				tmp := make([]*model.Command, 0)
				tmp = append(tmp, fetchJob.queue[0:index]...)
				tmp = append(tmp, fetchJob.queue[index+1:]...)
				fetchJob.queue = tmp
			}
		}
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
	h.jobs[c.Args[0]].queue = make([]*model.Command, 0)
	h.jobs[c.Args[0]].cond = sync.NewCond(&sync.Mutex{})
	h.jobs[c.Args[0]].last = -1
	h.mtx.Unlock()

	h.jobs[c.Args[0]].cond.L.Lock()
	transport.Send(c, tc)
	h.jobs[c.Args[0]].cond.Wait()
	h.jobs[c.Args[0]].cond.L.Unlock()

	h.mtx.Lock()
	if h.jobs[c.Args[0]].hadOrderIssue {
		panic(c.Args[1] + " had order issue")
	}
	delete(h.jobs, c.Args[0])
	h.mtx.Unlock()
}
