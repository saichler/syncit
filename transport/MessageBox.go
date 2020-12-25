package transport

import (
	log "github.com/saichler/utils/golang"
	"sync"
	"time"
)

type MessageBox struct {
	queue   [][]byte
	mtx     *sync.Cond
	maxSize int
}

func newMessageBox(maxSize int) *MessageBox {
	msgBox := &MessageBox{}
	msgBox.mtx = sync.NewCond(&sync.Mutex{})
	msgBox.queue = make([][]byte, 0)
	msgBox.maxSize = maxSize
	return msgBox
}

func (msgBox *MessageBox) push(packet []byte) {
	for len(msgBox.queue) >= msgBox.maxSize {
		log.Info("Queue Wait")
		time.Sleep(time.Second)
		msgBox.mtx.Broadcast()
	}
	log.Info("Pushing ", len(packet))
	msgBox.mtx.L.Lock()
	defer msgBox.mtx.L.Unlock()
	msgBox.queue = append(msgBox.queue, packet)
	msgBox.mtx.Broadcast()
}

func (msgBox *MessageBox) pop() []byte {
	for {
		msgBox.mtx.L.Lock()
		if len(msgBox.queue) == 0 {
			msgBox.mtx.Wait()
		}
		if len(msgBox.queue) > 0 {
			data := msgBox.queue[0]
			msgBox.queue = msgBox.queue[1:]
			msgBox.mtx.L.Unlock()
			return data
		}
		msgBox.mtx.L.Unlock()
	}
}
