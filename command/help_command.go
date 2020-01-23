package command

type HelpCommand struct{}

func (command *HelpCommand) Name() string {
	return "help"
}

func (command *HelpCommand) Help() string {
	return "this help text"
}

func (command *HelpCommand) Run([]string) string {
	output := "commands:\r\n"
	for _, command := range commands {
		output += command.Name() + " - " + command.Help() + "\r\n"
	}
	output += "quit - exit console"

	return output
}
