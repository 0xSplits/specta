package keypair

import (
	"fmt"

	"github.com/0xSplits/otelgo/recorder"
	"github.com/0xSplits/otelgo/registry"
	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

const (
	Metric = "aws_ec2_keypair_total"
)

type Config struct {
	Aws aws.Config
	Env envvar.Env
	Log logger.Interface
	Met metric.Meter
}

type Handler struct {
	ec2 *ec2.Client
	env envvar.Env
	log logger.Interface
	reg registry.Interface
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

	{
		cou[Metric] = recorder.NewCounter(recorder.CounterConfig{
			Des: "the total amount of EC2 key-pairs",
			Met: c.Met,
			Nam: Metric,
		})
	}

	gau := map[string]recorder.Interface{}
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
		ec2: ec2.NewFromConfig(c.Aws),
		env: c.Env,
		log: c.Log,
		reg: reg,
	}
}
