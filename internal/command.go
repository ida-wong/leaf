package internal

type Command interface {
	// must goroutine safe
	Name() string
	// must goroutine safe
	Help() string
	// must goroutine safe
	Run(args []string) string
}

var Commands = make(map[string]Command)
