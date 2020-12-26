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
		if file != nil {
			file.Close()
		}
	}()
	if err != nil {
		c.Response = []byte(err.Error())
		return
	}
	if f.Size() > transport.LARGE_PACKET {
		responseID := 0
		responseCount := 0
		file, err = os.Open(c.Args[0])
		if err != nil {
			c.Response = []byte(err.Error())
			return
		}
		responseCount = int(f.Size() / transport.LARGE_PACKET)
		if f.Size()%transport.LARGE_PACKET > 0 {
			responseCount++
		}
		for responseID < responseCount-1 {
			data := make([]byte, transport.LARGE_PACKET)
			file.Read(data)
			c.ResponseCount = int32(responseCount)
			c.ResponseId = int32(responseID)
			c.Response = data
			transport.Send(c, tc)
			responseID++
		}
		left := int(f.Size()) - transport.LARGE_PACKET*responseID
		if left > 0 {
			data := make([]byte, left)
			file.Read(data)
			c.ResponseId = int32(responseID)
			c.ResponseCount = int32(responseCount)
			c.Response = data
			transport.Send(c, tc)
		}
	} else {
		data, err := ioutil.ReadFile(c.Args[0])
		if err != nil {
			c.Response = []byte(err.Error())
		} else {
			c.Response = data
		}
		transport.Send(c, tc)
	}
}

func (h *Fetch) HandleResponse(command *model.Command, tc *transport.Connection) {
	h.mtx.Lock()
	fetchJob := h.jobs[command.Args[0]]
	h.mtx.Unlock()

	fetchJob.cond.L.Lock()
	if command.ResponseId != fetchJob.last+1 {
		fetchJob.waiting[command.ResponseId] = command
		fetchJob.hadOrderIssue = true
		log.Info("Part ", command.ResponseId, " of file:", command.Args[1])
		fetchJob.cond.L.Unlock()
		return
	}

	var file *os.File
	var err error

	if command.ResponseId == 0 {
		fmt.Print("Receiving ", command.Args[1], " with ", command.ResponseCount, " parts:.")
		index := strings.LastIndex(command.Args[1], "/")
		if index != -1 {
			dirPath := command.Args[1][0:index]
			_, exist := os.Stat(dirPath)
			if exist != nil {
				os.MkdirAll(dirPath, 0777)
			}
		}
	}

	defer func() {
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
			if command.ResponseId%10 == 0 {
				pr := float64(command.ResponseId) / float64(command.ResponseCount) * 100
				fmt.Print(".", int(pr), "%")
			} else {
				fmt.Print(".")
			}
		}
		fetchJob.cond.L.Unlock()
	}()

	dataToWrite := make([]byte, 0)

	if command.ResponseCount == 0 || command.ResponseId == 0 {
		file, err = os.Create(command.Args[1])
		if err != nil {
			return
		} else {
			if fetchJob.hadOrderIssue {
				log.Info("Writing part ", command.ResponseId, " of file:", command.Args[1])
			}
			dataToWrite = append(dataToWrite, command.Response...)
		}
	} else {
		file, err = os.OpenFile(command.Args[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			log.Error(err)
			return
		}
		if fetchJob.hadOrderIssue {
			log.Info("Writing part ", command.ResponseId, " size:", len(command.Response), " of file:", command.Args[1])
		}
		dataToWrite = append(dataToWrite, command.Response...)
	}
	fetchJob.last = command.ResponseId

	waitingCommand, ok := fetchJob.waiting[fetchJob.last+1]
	for ok {
		fmt.Print(".")
		log.Info("Writing part ", waitingCommand.ResponseId, " size:", len(waitingCommand.Response), " of file:", waitingCommand.Args[1])
		dataToWrite = append(dataToWrite, waitingCommand.Response...)
		fetchJob.last = waitingCommand.ResponseId
		delete(fetchJob.waiting, waitingCommand.ResponseId)
		waitingCommand, ok = fetchJob.waiting[fetchJob.last+1]
	}
	file.Write(dataToWrite)
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
		//panic(c.Args[1] + " had order issues, please check!")
	}
	delete(h.jobs, c.Args[0])
	h.mtx.Unlock()
}
