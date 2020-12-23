package transport

import (
	"errors"
	log "github.com/saichler/utils/golang"
	"net"
	"strconv"
	"sync"
)

type Listener struct {
	port        int
	secret      string
	key         string
	socket      net.Listener
	mtx         *sync.Mutex
	running     bool
	msgListener MessageListener
}

func NewListener(port int, secret, key string, ml MessageListener) *Listener {
	l := &Listener{}
	l.port = port
	l.secret = secret
	l.key = key
	l.msgListener = ml
	return l
}

func (l *Listener) bind() error {
	log.Info("Listening on 0.0.0.0:", l.port)

	socket, e := net.Listen("tcp", ":"+strconv.Itoa(l.port))

	if e != nil {
		return e
	}
	l.socket = socket
	l.mtx = &sync.Mutex{}
	return nil
}

func (l *Listener) Listen() error {
	if l.port == 0 {
		return errors.New("Listener does not have a port defined")
	}
	if l.secret == "" {
		return errors.New("Listener does not have a secret")
	}
	if l.key == "" {
		return errors.New("Listener does not have a key")
	}
	err := l.bind()
	if err != nil {
		return err
	}
	l.running = true
	for l.running {
		conn, e := l.socket.Accept()
		if e != nil {
			return e
		}
		l.addService(conn)
	}
	return nil
}

func (l *Listener) addService(conn net.Conn) {
	err := connect(conn, l.key, l.secret, l.msgListener)
	if err != nil {
		log.Error("Failed to connect to ", conn.RemoteAddr().String(), " ", err)
	}
}
