package daemon

import (
	"fmt"

	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/0xSplits/specta/pkg/registry"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type Config struct {
	Env envvar.Env
}

type Daemon struct {
	env envvar.Env
	log logger.Interface
	reg registry.Interface
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

	var reg registry.Interface
	{
		reg = registry.New(registry.Config{
			Log: log,
			Met: musMet(c),
		})
	}

	return &Daemon{
		env: c.Env,
		log: log,
		reg: reg,
	}
}

func metNam(env string) string {
	return fmt.Sprintf("specta.%s.splits.org", env)
}

func musMet(cfg Config) otelmetric.Meter {
	exp, err := prometheus.New()
	if err != nil {
		tracer.Panic(tracer.Mask(err))
	}

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exp),
	).Meter(metNam(cfg.Env.Environment))
}
