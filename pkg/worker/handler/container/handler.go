package container

import (
	"fmt"

	"github.com/0xSplits/otelgo/recorder"
	"github.com/0xSplits/otelgo/registry"
	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

const (
	Metric = "aws_ecs_container_health"
)

type Config struct {
	Aws aws.Config
	Env envvar.Env
	Log logger.Interface
	Met metric.Meter
}

type Handler struct {
	ecs *ecs.Client
	env envvar.Env
	log logger.Interface
	reg registry.Interface
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

	cou := map[string]recorder.Interface{}

	gau := map[string]recorder.Interface{}

	{
		gau[Metric] = recorder.NewGauge(recorder.GaugeConfig{
			Des: "the health status of ecs service containers",
			Lab: map[string][]string{
				"service": {"alloy", "graphql", "kayron", "otel-collector", "server", "specta", "worker"},
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
		ecs: ecs.NewFromConfig(c.Aws),
		env: c.Env,
		log: c.Log,
		reg: reg,
		tag: resourcegroupstaggingapi.NewFromConfig(c.Aws),
	}
}
