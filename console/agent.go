package console

import (
	"bufio"
	"strings"

	"github.com/ida-wong/leaf/config"
	. "github.com/ida-wong/leaf/internal"
	"github.com/ida-wong/leaf/network"
)

type Agent struct {
	conn   *network.TCPConn
	reader *bufio.Reader
}

func newAgent(conn *network.TCPConn) network.Agent {
	a := new(Agent)
	a.conn = conn
	a.reader = bufio.NewReader(conn)
	return a
}

func (a *Agent) Run() {
	for {
		if config.ConsolePrompt != "" {
			a.conn.Write([]byte(config.ConsolePrompt))
		}

		line, err := a.reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line[:len(line)-1], "\r")

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		name := args[0]
		if name == "quit" {
			break
		}

		command, exists := Commands[name]
		if !exists {
			a.conn.Write([]byte("command not found, try `help` for help\r\n"))
			continue
		}

		output := command.Run(args[1:])
		if output != "" {
			a.conn.Write([]byte(output + "\r\n"))
		}
	}
}

func (a *Agent) OnClose() {

}
