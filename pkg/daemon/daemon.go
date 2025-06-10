package daemon

import (
	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/xh3b4sd/logger"
)

type Config struct {
	Env envvar.Env
}

type Daemon struct {
	env envvar.Env
	log logger.Interface
}

func New(c Config) *Daemon {
	var log logger.Interface
	{
		log = logger.New(logger.Config{
			Filter: logger.NewLevelFilter(c.Env.LogLevel),
		})
	}

	log.Log(
		"level", "info",
		"message", "daemon is starting",
		"environment", c.Env.Environment,
	)

	return &Daemon{
		env: c.Env,
		log: log,
	}
}
