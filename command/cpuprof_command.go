package command

import (
	"fmt"
	"os"
	"path"
	"runtime/pprof"
	"time"

	"github.com/ida-wong/leaf/config"
)

type CpuProfCommand struct{}

func (command *CpuProfCommand) Name() string {
	return "cpuprof"
}

func (command *CpuProfCommand) Help() string {
	return "CPU profiling for the current process"
}

func (command *CpuProfCommand) Run(args []string) string {
	if len(args) == 0 {
		return command.usage()
	}

	switch args[0] {
	case "start":
		fn := profileName() + ".cpuprof"
		f, err := os.Create(fn)
		if err != nil {
			return err.Error()
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			f.Close()
			return err.Error()
		}
		return fn
	case "stop":
		pprof.StopCPUProfile()
		return ""
	default:
		return command.usage()
	}
}

func (command *CpuProfCommand) usage() string {
	return "cpuprof writes runtime profiling data in the format expected by \r\n" +
		"the pprof visualization tool\r\n\r\n" +
		"Usage: cpuprof start|stop\r\n" +
		"  start - enables CPU profiling\r\n" +
		"  stop  - stops the current CPU profile"
}

func profileName() string {
	now := time.Now()
	return path.Join(config.ProfilePath,
		fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second()))
}
