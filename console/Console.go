package main

import (
	"bufio"
	"fmt"
	"github.com/saichler/security"
	"github.com/saichler/syncit/cmd"
	"github.com/saichler/syncit/transport"
	log "github.com/saichler/utils/golang"
	"os"
	"strconv"
	"strings"
)

type Console struct {
	commandHanler *cmd.CommandHandler
	service       *transport.Listener
	tc            *transport.Connection
}

const (
	FILENAME = "./console.io"
	MYK      = "/my/k"
	MYS      = "/my/s"
	MYP      = "/my/p"
)

var prompt = "Sync-it->"
var running = true

func main() {
	con := &Console{}
	con.commandHanler = &cmd.CommandHandler{}

	_, err := os.Stat(FILENAME)
	if err != nil {
		st := security.InitSecureStore(FILENAME)
		st.Put(MYK, security.GenerateAES256Key())
		st.Put(MYS, "sync-it")
		st.Put(MYP, "45454")
	}

	if len(os.Args) > 1 && os.Args[1] == "service" {
		con.startService()
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for running {
		fmt.Print(prompt)
		cmd, _ := reader.ReadString('\n')
		cmd = cmd[0 : len(cmd)-1]
		con.processCommand(cmd)
	}
	fmt.Println("Goodbye!")
}

func (con *Console) processCommand(command string) {
	if command == "exit" || command == "quit" {
		running = false
		return
	} else if command == "" {
		return
	}
	args := strings.Split(command, " ")
	if args == nil || len(args) == 0 {
		return
	}

	if args[0] == "gk" {
		st := security.InitSecureStore(FILENAME)
		k, _ := st.Get(MYK)
		log.Info("MYK=", k)
		return
	} else if args[0] == "sk" {
		st := security.InitSecureStore(FILENAME)
		st.Put(MYK, security.GenerateAES256Key())
		return
	} else if args[0] == "gs" {
		st := security.InitSecureStore(FILENAME)
		s, _ := st.Get(MYS)
		log.Info("MYS=", s)
		return
	} else if args[0] == "ss" {
		st := security.InitSecureStore(FILENAME)
		st.Put(MYS, args[1])
		return
	} else if args[0] == "gp" {
		st := security.InitSecureStore(FILENAME)
		p, _ := st.Get(MYP)
		log.Info("MYP=", p)
		return
	} else if args[0] == "sp" {
		st := security.InitSecureStore(FILENAME)
		st.Put(MYP, args[1])
		return
	}
	if args[0] == "service" {
		go con.startService()
	} else if args[0] == "connect" {
		con.connect(args[1:])
	} else {
		con.commandHanler.Execute(args[0], args[1:], con.tc)
	}
}

func (con *Console) connect(args []string) {
	st := security.InitSecureStore(FILENAME)
	if args == nil || len(args) == 0 {
		log.Error("To connect you need the following args <host> <port> <key> <secret>")
		return
	}

	host := ""
	port := 45454
	key := ""
	secret := ""

	if len(args) == 4 {
		host = args[0]
		port, _ = strconv.Atoi(args[1])
		key = args[2]
		secret = args[3]
		prefix := "/" + host + "/"
		st.Put(prefix+"p", strconv.Itoa(port))
		st.Put(prefix+"k", key)
		st.Put(prefix+"s", secret)
	} else if len(args) == 1 {
		host = args[0]
		prefix := "/" + host + "/"
		p, _ := st.Get(prefix + "p")
		port, _ = strconv.Atoi(p)
		key, _ = st.Get(prefix + "k")
		secret, _ = st.Get(prefix + "s")
		if port == 0 || key == "" || secret == "" {
			log.Error("To connect you need the following args <host> <port> <key> <secret>")
			return
		}
	} else {
		log.Error("To connect you need the following args <host> <port> <key> <secret>")
		return
	}

	if args == nil || len(args) != 3 {
		log.Error("Connectg needs 3 args <host> <key> <secret>")
		return
	}

	tc, err := transport.Connect(host, key, secret, port, con.commandHanler)
	if err != nil {
		log.Error("Unable to connect:", err)
		return
	}
	con.tc = tc
	log.Info("Connected!")
}

func (con *Console) startService() {
	st := security.InitSecureStore(FILENAME)
	p, _ := st.Get(MYP)
	port, _ := strconv.Atoi(p)
	key, _ := st.Get(MYK)
	secret, _ := st.Get(MYS)

	con.service = transport.NewListener(port, secret, key, con.commandHanler)
	err := con.service.Listen()
	if err != nil {
		log.Error("Failed to start service:", err)
	}
}
