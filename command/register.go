package command

import (
	"github.com/ida-wong/leaf/chanrpc"
	. "github.com/ida-wong/leaf/internal"
	"github.com/ida-wong/leaf/log"
)

var commands = Commands

// you must call the function before calling console.Init
// goroutine not safe
func Register(name string, help string, f interface{}, server *chanrpc.Server) {
	if _, exists := commands[name]; exists {
		log.Fatal("command %v is already registered", name)
	}

	server.Register(name, f)

	command := new(ExternalCommand)
	command.name = name
	command.help = help
	command.server = server
	commands[name] = command
}