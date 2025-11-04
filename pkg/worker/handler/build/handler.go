package build

import (
	"fmt"
	"time"

	"github.com/0xSplits/otelgo/recorder"
	"github.com/0xSplits/otelgo/registry"
	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/google/go-github/v75/github"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

const (
	MetricTotal    = "build_execution_total"
	MetricDuration = "build_execution_duration_seconds"
)

const (
	Organization = "0xSplits"
)

type Config struct {
	Env envvar.Env
	Git *github.Client
	Log logger.Interface
	Met metric.Meter
}

type Handler struct {
	cac *lru.Cache[int64, struct{}]
	env envvar.Env
	git *github.Client
	log logger.Interface
	reg registry.Interface
	sta time.Time
}

func New(c Config) *Handler {
	if c.Git == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Git must not be empty", c)))
	}
	if c.Log == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Log must not be empty", c)))
	}
	if c.Met == nil {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Met must not be empty", c)))
	}

	var err error

	cou := map[string]recorder.Interface{}

	{
		cou[MetricTotal] = recorder.NewCounter(recorder.CounterConfig{
			Des: "the total amount of build container executions",
			Lab: map[string][]string{
				"platform":   {"github"},
				"repository": {"kayron", "pulsar", "server", "specta", "splits-lite"},
				"success":    {"true", "false"},
				"workflow":   {"check", "image"},
			},
			Met: c.Met,
			Nam: MetricTotal,
		})
	}

	gau := map[string]recorder.Interface{}

	his := map[string]recorder.Interface{}

	{
		his[MetricDuration] = recorder.NewHistogram(recorder.HistogramConfig{
			Des: "the time it takes for build container executions to complete",
			Lab: map[string][]string{
				"platform":   {"github"},
				"repository": {"kayron", "pulsar", "server", "specta", "splits-lite"},
				"success":    {"true", "false"},
				"workflow":   {"check", "image"},
			},
			Buc: []float64{
				30,  // 0.5 min
				60,  // 1.0 min
				90,  // 1.5 min
				120, // 2.0 min
				150, // 2.5 min

				180, // 3.0 min
				210, // 3.5 min
				240, // 4.0 min
				270, // 4.5 min
				300, // 5.0 min
			},
			Met: c.Met,
			Nam: MetricDuration,
		})
	}

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

	var cac *lru.Cache[int64, struct{}]
	{
		cac, err = lru.New[int64, struct{}](max + 1)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	var sta time.Time
	{
		sta = time.Now().UTC()
	}

	return &Handler{
		cac: cac,
		env: c.Env,
		git: c.Git,
		log: c.Log,
		reg: reg,
		sta: sta,
	}
}
