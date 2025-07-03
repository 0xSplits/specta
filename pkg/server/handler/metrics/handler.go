package metrics

import (
	"fmt"

	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/0xSplits/specta/pkg/recorder"
	"github.com/0xSplits/specta/pkg/registry"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
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

	his := map[string]recorder.Interface{}

	{
		nam := "page_ready_duration_seconds"
		his[nam] = recorder.NewHistogram(recorder.HistogramConfig{
			Des: "the time it takes for the app to show a page",
			Lab: map[string][]string{
				"page": {"root"},
			},
			Buc: []float64{
				0.1, //   100ms
				0.2, //   200ms
				0.3, //   300ms
				0.4, //   400ms
				0.5, //   500ms

				1.0, // 1,000ms
				2.0, // 2,000ms
				3.0, // 3,000ms
				4.0, // 4,000ms
				5.0, // 5,000ms
			},
			Met: c.Met,
			Nam: nam,
		})
	}

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
