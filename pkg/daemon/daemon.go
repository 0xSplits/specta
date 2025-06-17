package daemon

import (
	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/0xSplits/specta/pkg/recorder"
	"github.com/xh3b4sd/logger"
	"go.opentelemetry.io/otel/metric"
	otelmetric "go.opentelemetry.io/otel/metric"
)

type Config struct {
	Env envvar.Env
}

type Daemon struct {
	env envvar.Env
	log logger.Interface
	met otelmetric.Meter
}

func New(c Config) *Daemon {
	var log logger.Interface
	{
		log = logger.New(logger.Config{
			Filter: logger.NewLevelFilter(c.Env.LogLevel),
		})
	}

	var met metric.Meter
	{
		met = recorder.NewMeter(recorder.MeterConfig{
			Env: c.Env.Environment,
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
		met: met,
	}
}
