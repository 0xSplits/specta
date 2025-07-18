package endpoint

import (
	"fmt"

	"github.com/0xSplits/otelgo/recorder"
	"github.com/0xSplits/otelgo/registry"
	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

const (
	Metric = "http_endpoint_health"
)

var (
	mapping = map[string]map[string][]string{
		"explorer": {
			"testing":    {"https://explorer.testing.splits.org"},
			"staging":    {"https://explorer.staging.splits.org"},
			"production": {"https://explorer.production.splits.org", "https://app.splits.org"},
		},
		"server": {
			"testing":    {"https://server.testing.splits.org/metrics"},
			"staging":    {"https://server.staging.splits.org/metrics"},
			"production": {"https://server.production.splits.org/metrics", "https://api.splits.org/metrics"},
		},
		"specta": {
			"testing":    {"https://specta.testing.splits.org/metrics"},
			"staging":    {"https://specta.staging.splits.org/metrics"},
			"production": {"https://specta.production.splits.org/metrics"},
		},
		"teams": {
			"testing":    {"https://teams.testing.splits.org"},
			"staging":    {"https://teams.staging.splits.org"},
			"production": {"https://teams.production.splits.org", "https://teams.splits.org"},
		},
	}
)

type Config struct {
	Env envvar.Env
	Log logger.Interface
	Met metric.Meter
}

type Handler struct {
	env envvar.Env
	log logger.Interface
	reg registry.Interface
}

func New(c Config) *Handler {
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Met == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Met must not be empty", c)))
	}

	cou := map[string]recorder.Interface{}

	gau := map[string]recorder.Interface{}

	{
		gau[Metric] = recorder.NewGauge(recorder.GaugeConfig{
			Des: "the health status of an http endpoint",
			Lab: map[string][]string{
				"service": {"explorer", "server", "specta", "teams"},
			},
			Met: c.Met,
			Nam: Metric,
		})
	}

	his := map[string]recorder.Interface{}

	var reg registry.Interface
	{
		reg = registry.New(registry.Config{
			Env: c.Env.Environment,
			Log: c.Log,

			Cou: cou,
			Gau: gau,
			His: his,
		})
	}

	return &Handler{
		env: c.Env,
		log: c.Log,
		reg: reg,
	}
}
