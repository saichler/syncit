package main

import (
	"bufio"
	"fmt"
	"github.com/saichler/syncit/cmd"
	"github.com/saichler/syncit/transport"
	log "github.com/saichler/utils/golang"
	"math/rand"
	"os"
	"strings"
	"time"
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

	if args[0] == "service" {
		go con.startService(args[1:])
	} else if args[0] == "connect" {
		con.connect(args[1:])
	} else {
		con.commandHanler.Execute(args[0], args[1:], con.tc)
	}
}

func (con *Console) connect(args []string) {
	host := "127.0.0.1"
	port := 45454
	key := "qHYsJuloczNsFrbqlhlffjkRuHWfrCtH"
	secret := "syncit"
	tc, err := transport.Connect(host, key, secret, port, con.commandHanler)
	if err != nil {
		log.Error("Unable to connect:", err)
		return
	}
	con.tc = tc
	log.Info("Connected!")
}

func (con *Console) startService(args []string) {
	if args == nil || len(args) != 1 {
		log.Error("service command needs 1 argument as secret")
		return
	}

	port := 45454
	key := GenerateAES256Key()
	secret := args[0]

	log.Info("Key=", key, " secret=", secret, " port=", port)
	con.service = transport.NewListener(port, secret, key, con.commandHanler)
	err := con.service.Listen()
	if err != nil {
		log.Error("Failed to start service:", err)
	}
}

var l = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateAES256Key() string {
	rand.Seed(time.Now().UnixNano())
	key := make([]rune, 32)
	for i := range key {
		key[i] = l[rand.Intn(len(l))]
	}
	return string(key)
}
