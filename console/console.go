package console

import (
	"fmt"
	"github.com/ida-wong/leaf/log"
	"math"

	. "github.com/ida-wong/leaf/command"
	"github.com/ida-wong/leaf/config"
	. "github.com/ida-wong/leaf/internal"
	"github.com/ida-wong/leaf/network"
)

var server *network.TCPServer

func Init() {
	consolePort := config.ConsolePort
	if consolePort == 0 {
		return
	}

	log.Release("Console port: %d", consolePort)

	server = new(network.TCPServer)
	server.Addr = fmt.Sprintf("localhost:%d", consolePort)
	server.MaxConnNum = int(math.MaxInt32)
	server.PendingWriteNum = 100
	server.NewAgent = newAgent

	server.Start()

	commands := []Command{
		new(CpuProfCommand),
		new(ExternalCommand),
		new(HelpCommand),
		new(ProfCommand),
	}

	for _, item := range commands {
		Commands[item.Name()] = item
	}
}

func Destroy() {
	if server != nil {
		server.Close()
	}
}
