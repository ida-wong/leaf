package command

import "github.com/ida-wong/leaf/chanrpc"

type ExternalCommand struct {
	name   string
	help   string
	server *chanrpc.Server
}

func (command *ExternalCommand) Name() string {
	return command.name
}

func (command *ExternalCommand) Help() string {
	return command.help
}

func (command *ExternalCommand) Run(_args []string) string {
	args := make([]interface{}, len(_args))
	for i, v := range _args {
		args[i] = v
	}

	ret, err := command.server.Call1(command.name, args...)
	if err != nil {
		return err.Error()
	}
	output, ok := ret.(string)
	if !ok {
		return "invalid output type"
	}

	return output
}
