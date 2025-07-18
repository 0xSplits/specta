package deployment

import (
	"fmt"
	"time"

	"github.com/0xSplits/otelgo/recorder"
	"github.com/0xSplits/otelgo/registry"
	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

const (
	MetricTotal    = "deployment_execution_total"
	MetricDuration = "deployment_execution_duration_seconds"
)

type Config struct {
	Aws aws.Config
	Env envvar.Env
	Log logger.Interface
	Met metric.Meter
}

type Handler struct {
	acp *codepipeline.Client
	cac *lru.Cache[string, struct{}]
	env envvar.Env
	log logger.Interface
	reg registry.Interface
	sta time.Time
	tag *resourcegroupstaggingapi.Client
}

func New(c Config) *Handler {
	if c.Aws.Region == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Aws must not be empty", c)))
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
			Des: "the total amount of deployment pipeline executions",
			Lab: map[string][]string{
				"platform": {"codepipeline"},
				"success":  {"true", "false"},
			},
			Met: c.Met,
			Nam: MetricTotal,
		})
	}

	gau := map[string]recorder.Interface{}

	his := map[string]recorder.Interface{}

	{
		his[MetricDuration] = recorder.NewHistogram(recorder.HistogramConfig{
			Des: "the time it takes for deployment pipeline executions to complete",
			Lab: map[string][]string{
				"platform": {"codepipeline"},
				"success":  {"true", "false"},
			},
			Buc: []float64{
				60,  //  1 min
				120, //  2 min
				180, //  3 min
				240, //  4 min
				300, //  5 min

				420, //  7 min
				540, //  9 min
				660, // 11 min
				780, // 13 min
				900, // 15 min
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

	var cac *lru.Cache[string, struct{}]
	{
		cac, err = lru.New[string, struct{}](max + 1)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	var sta time.Time
	{
		sta = time.Now().UTC()
	}

	return &Handler{
		acp: codepipeline.NewFromConfig(c.Aws),
		cac: cac,
		env: c.Env,
		log: c.Log,
		reg: reg,
		sta: sta,
		tag: resourcegroupstaggingapi.NewFromConfig(c.Aws),
	}
}
