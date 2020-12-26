package transport

import (
	log "github.com/saichler/utils/golang"
	"sync"
)

type MessageBox struct {
	name    string
	queue   [][]byte
	mtx     *sync.Cond
	maxSize int
	running bool
}

func newMessageBox(name string, maxSize int) *MessageBox {
	msgBox := &MessageBox{}
	msgBox.mtx = sync.NewCond(&sync.Mutex{})
	msgBox.queue = make([][]byte, 0)
	msgBox.maxSize = maxSize
	msgBox.running = true
	msgBox.name = name
	return msgBox
}

func (msgBox *MessageBox) push(packet []byte) {
	msgBox.mtx.L.Lock()
	defer msgBox.mtx.L.Unlock()

	for len(msgBox.queue) >= msgBox.maxSize && msgBox.running {
		msgBox.mtx.Broadcast()
		msgBox.mtx.Wait()
	}
	if msgBox.running {
		msgBox.queue = append(msgBox.queue, packet)
	} else {
		msgBox.queue = msgBox.queue[0:0]
	}
	msgBox.mtx.Broadcast()
}

func (msgBox *MessageBox) pop() []byte {
	for msgBox.running {
		msgBox.mtx.L.Lock()
		if len(msgBox.queue) == 0 {
			msgBox.mtx.Broadcast()
			msgBox.mtx.Wait()
		}
		if len(msgBox.queue) > 0 {
			data := msgBox.queue[0]
			msgBox.queue = msgBox.queue[1:]
			msgBox.mtx.Broadcast()
			msgBox.mtx.L.Unlock()
			return data
		}
		msgBox.mtx.L.Unlock()
	}
	log.Info("Message Box", msgBox.name, " has stopped.")
	return nil
}

func (msgBox *MessageBox) Shutdown() {
	msgBox.running = false
	msgBox.mtx.Broadcast()
}
