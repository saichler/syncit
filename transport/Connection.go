package transport

import (
	"errors"
	"fmt"
	log "github.com/saichler/utils/golang"
	"net"
	"strconv"
	"sync"
)

type Connection struct {
	key         string
	inbox       *MessageBox
	outbox      *MessageBox
	conn        net.Conn
	running     bool
	msgListener MessageListener
	writeMutex  *sync.Cond
	name        string
}

func newConnection(con net.Conn, key string, ml MessageListener) *Connection {
	c := &Connection{}
	c.conn = con
	c.name = con.RemoteAddr().String()
	c.inbox = newMessageBox("Inbox "+c.name, 1000)
	c.outbox = newMessageBox("Outbox "+c.name, 20)
	c.running = true
	c.key = key
	c.msgListener = ml
	c.writeMutex = sync.NewCond(&sync.Mutex{})
	return c
}

func connect(conn net.Conn, key, secret string, ml MessageListener) error {
	service := newConnection(conn, key, ml)
	initData, err := readPacket(service.conn)
	if err != nil {
		return err
	}

	data, err := decode(string(initData), key)
	if err != nil {
		conn.Close()
		return err
	}

	if string(data) != secret {
		conn.Close()
		return errors.New("Incorrect Secret/Key, aborting connection")
	}
	writePacket([]byte("OK"), conn)
	go service.read()
	go service.write()
	go service.process()
	return nil
}

func Connect(host, key, secret string, port int, ml MessageListener) (*Connection, error) {
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	data, err := encode([]byte(secret), key)
	if err != nil {
		return nil, err
	}

	err = writePacket([]byte(data), conn)
	if err != nil {
		return nil, err
	}

	inData, err := readPacket(conn)

	if string(inData) != "OK" {
		return nil, errors.New("Failed to connect, incorrect Key/Secret")
	}

	service := newConnection(conn, key, ml)
	go service.read()
	go service.write()
	go service.process()
	return service, nil
}

func (c *Connection) read() {
	for c.running {
		packet, err := readPacket(c.conn)
		if err != nil {
			log.Error(err)
			break
		}
		if packet != nil {
			if len(packet) == 2 && string(packet) == "WC" {
				c.writeMutex.L.Lock()
				c.writeMutex.Broadcast()
				c.writeMutex.L.Unlock()
				continue
			} else if len(packet) >= LARGE_PACKET {
				/*
					c.writeMutex.L.Lock()
					writePacket([]byte("WC"), c.conn)
					c.writeMutex.L.Unlock()
				*/
			}
			if c.running {
				c.inbox.push(packet)
			}
		} else {
			break
		}
	}
	log.Info("Connection Read for ", c.name, " ended.")
	c.Shutdown()
}

func (c *Connection) Shutdown() {
	c.running = false
	if c.conn != nil {
		c.conn.Close()
	}
	c.inbox.Shutdown()
	c.outbox.Shutdown()
	c.writeMutex.Broadcast()
}

func (c *Connection) write() {
	for c.running {
		packet := c.outbox.pop()
		if packet != nil {
			c.writeMutex.L.Lock()
			if c.running {
				writePacket(packet, c.conn)
			}
			if len(packet) >= LARGE_PACKET {
				//c.writeMutex.Wait()
			}
			c.writeMutex.L.Unlock()
		} else {
			break
		}
	}
	log.Info("Connection Write for ", c.name, " ended.")
	c.Shutdown()
}

func (c *Connection) Send(data []byte) error {
	if c.running {
		encData, err := encode(data, c.key)
		if err != nil {
			return err
		}
		c.outbox.push([]byte(encData))
	}
	return nil
}

func (c *Connection) process() {
	for c.running {
		packet := c.inbox.pop()
		if packet != nil {
			encString := string(packet)
			data, err := decode(encString, c.key)
			if err != nil {
				break
			}
			c.handleMessage(data)
		}
	}
	log.Info("Message Processing for ", c.name, " Ended")
	c.Shutdown()
}

func (c *Connection) handleMessage(msg []byte) {
	if c.running {
		if c.msgListener != nil {
			go c.msgListener.HandleMessage(msg, c)
		} else {
			fmt.Println(string(msg))
		}
	}
}
