package command

import (
	"os"
	"runtime/pprof"
)

type ProfCommand struct{}

func (command *ProfCommand) Name() string {
	return "prof"
}

func (command *ProfCommand) Help() string {
	return "writes a pprof-formatted snapshot"
}

func (command *ProfCommand) Run(args []string) string {
	if len(args) == 0 {
		return command.usage()
	}

	var (
		p  *pprof.Profile
		fn string
	)
	switch args[0] {
	case "goroutine":
		p = pprof.Lookup("goroutine")
		fn = profileName() + ".gprof"
	case "heap":
		p = pprof.Lookup("heap")
		fn = profileName() + ".hprof"
	case "thread":
		p = pprof.Lookup("threadcreate")
		fn = profileName() + ".tprof"
	case "block":
		p = pprof.Lookup("block")
		fn = profileName() + ".bprof"
	default:
		return command.usage()
	}

	f, err := os.Create(fn)
	if err != nil {
		return err.Error()
	}
	defer f.Close()
	err = p.WriteTo(f, 0)
	if err != nil {
		return err.Error()
	}

	return fn
}

func (command *ProfCommand) usage() string {
	return "prof writes runtime profiling data in the format expected by \r\n" +
		"the pprof visualization tool\r\n\r\n" +
		"Usage: prof goroutine|heap|thread|block\r\n" +
		"  goroutine - stack traces of all current goroutines\r\n" +
		"  heap      - a sampling of all heap allocations\r\n" +
		"  thread    - stack traces that led to the creation of new OS threads\r\n" +
		"  block     - stack traces that led to blocking on synchronization primitives"
}
