package leaf

// reference: https://github.com/name5566/leaf
import (
	"os"
	"os/signal"

	"github.com/ida-wong/leaf/cluster"
	"github.com/ida-wong/leaf/config"
	"github.com/ida-wong/leaf/console"
	"github.com/ida-wong/leaf/log"
	"github.com/ida-wong/leaf/module"
)

func Run(mods ...module.Module) {
	// logger
	if config.LogLevel != "" {
		logger, err := log.New(config.LogLevel, config.LogPath, config.LogFlag)
		if err != nil {
			panic(err)
		}

		log.Export(logger)
		defer logger.Close()
	}

	log.Release("Leaf %v starting up", version)

	// close signal
	closeSig := make(chan os.Signal, 1)
	signal.Notify(closeSig, os.Interrupt, os.Kill)

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init(closeSig)

	// close
	sig := <-closeSig
	log.Release("Leaf closing down (signal: %v)", sig)
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}
