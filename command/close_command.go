package command

import(
	"os"
)

type CloseCommand struct{
	CloseSignal chan<- os.Signal
}

func (command *CloseCommand) Name() string {
	return "close"
}

func (command *CloseCommand) Help() string {
	return "close leaf process"
}

func (command *CloseCommand) Run([]string) string {
	command.CloseSignal <- os.Interrupt

	return ""
}
