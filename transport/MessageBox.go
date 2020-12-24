package transport

import (
	"sync"
)

type MessageBox struct {
	queue [][]byte
	mtx   *sync.Cond
}

func newMessageBox() *MessageBox {
	msgBox := &MessageBox{}
	msgBox.mtx = sync.NewCond(&sync.Mutex{})
	msgBox.queue = make([][]byte, 0)
	return msgBox
}

func (msgBox *MessageBox) push(packet []byte) {
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
