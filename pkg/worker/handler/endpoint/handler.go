package endpoint

import (
	"fmt"

	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/0xSplits/specta/pkg/recorder"
	"github.com/0xSplits/specta/pkg/registry"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

const (
	Metric = "http_endpoint_health"
)

var (
	endpoint = map[string]map[string]string{
		"explorer": {
			"testing":    "https://test.app.splits.org",
			"staging":    "https://beta.app.splits.org",
			"production": "https://app.splits.org",
		},
		"server": {
			"testing":    "https://test.api.splits.org/metrics",
			"staging":    "https://beta.api.splits.org/metrics",
			"production": "https://api.splits.org/metrics",
		},
		"specta": {
			"testing":    "https://specta.testing.splits.org/metrics",
			"staging":    "https://specta.staging.splits.org/metrics",
			"production": "https://specta.production.splits.org/metrics",
		},
		"teams": {
			"testing":    "https://test.teams.splits.org",
			"staging":    "https://beta.teams.splits.org",
			"production": "https://teams.splits.org",
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
			Env: c.Env,
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
