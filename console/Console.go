package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/saichler/security"
	"github.com/saichler/syncit/cmd"
	"github.com/saichler/syncit/model"
	"github.com/saichler/syncit/transport"
	log "github.com/saichler/utils/golang"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type Console struct {
	commandHanler *cmd.CommandHandler
	service       *transport.Listener
	tc            *transport.Connection
}

var prompt = "Sync-it->"
var running = true

func main() {
	con := &Console{}
	con.commandHanler = &cmd.CommandHandler{}
	model.InitSt()

	if len(os.Args) > 1 && os.Args[1] == "service" {
		go con.startService()
		waitForSignals()
		return
	} else if len(os.Args) > 1 {
		args := os.Args[1:]
		for _, arg := range args {
			fmt.Println("Processing command:", arg)
			con.processCommand(arg)
		}
		waitForSignals()
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

func waitForSignals() {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	waitForSingle := <-osSignal
	fmt.Println("Received: ", waitForSingle.String(), ", shutting down.")
}

func getCommandAndArgs(str string) (string, []string) {
	index := strings.Index(str, " ")
	if index == -1 {
		return str, []string{}
	}
	command := str[0:index]
	str = str[index+1:]
	args := make([]string, 0)
	qo := false
	buff := bytes.Buffer{}
	for _, c := range str {
		if c == '"' || c == '\'' {
			qo = !qo
		} else if c == ' ' && !qo {
			args = append(args, buff.String())
			buff = bytes.Buffer{}
		} else {
			buff.WriteString(string(c))
		}
	}

	if buff.String() != "" {
		args = append(args, buff.String())
	}

	return command, args
}

func (con *Console) processCommand(input string) {
	command, args := getCommandAndArgs(input)
	if command == "exit" || command == "quit" {
		running = false
		return
	} else if command == "" {
		return
	}

	if command == "gk" {
		st := security.InitSecureStore(model.IO_FILENAME)
		k, _ := st.Get(model.MYK)
		log.Info("MYK=", k)
		return
	} else if command == "sk" {
		st := security.InitSecureStore(model.IO_FILENAME)
		st.Put(model.MYK, security.GenerateAES256Key())
		return
	} else if command == "gs" {
		st := security.InitSecureStore(model.IO_FILENAME)
		s, _ := st.Get(model.MYS)
		log.Info("MYS=", s)
		return
	} else if command == "ss" {
		st := security.InitSecureStore(model.IO_FILENAME)
		st.Put(model.MYS, args[0])
		model.Secret = args[0]
		return
	} else if command == "gp" {
		st := security.InitSecureStore(model.IO_FILENAME)
		p, _ := st.Get(model.MYP)
		log.Info("MYP=", p)
		return
	} else if command == "sp" {
		st := security.InitSecureStore(model.IO_FILENAME)
		st.Put(model.MYP, args[0])
		return
	}
	if command == "service" {
		go con.startService()
	} else if command == "connect" {
		con.connect(args)
	} else {
		con.commandHanler.Execute(command, args, con.tc)
	}
}

func (con *Console) connect(args []string) {
	st := security.InitSecureStore(model.IO_FILENAME)
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

	tc, err := transport.Connect(host, key, secret, port, con.commandHanler)
	if err != nil {
		log.Error("Unable to connect:", err)
		return
	}
	con.tc = tc
	log.Info("Connected!")
}

func (con *Console) startService() {
	st := security.InitSecureStore(model.IO_FILENAME)
	p, _ := st.Get(model.MYP)
	port, _ := strconv.Atoi(p)
	key, _ := st.Get(model.MYK)
	secret, _ := st.Get(model.MYS)

	con.service = transport.NewListener(port, secret, key, con.commandHanler)
	err := con.service.Listen()
	if err != nil {
		log.Error("Failed to start service:", err)
	}
}
