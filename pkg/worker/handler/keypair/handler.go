package keypair

import (
	"fmt"

	"github.com/0xSplits/specta/pkg/envvar"
	"github.com/0xSplits/specta/pkg/recorder"
	"github.com/0xSplits/specta/pkg/registry"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/metric"
)

type Config struct {
	Aws aws.Config
	Log logger.Interface
	Met metric.Meter
}

type Handler struct {
	ec2 *ec2.Client
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
		nam := "aws_ec2_keypair_total"
		cou[nam] = recorder.NewCounter(recorder.CounterConfig{
			Des: "the total amount of EC2 key-pairs",
			Met: c.Met,
			Nam: nam,
		})
	}

	gau := map[string]recorder.Interface{}
	his := map[string]recorder.Interface{}

	var reg registry.Interface
	{
		reg = registry.New(registry.Config{
			Env: envvar.Env{},
			Log: c.Log,

			Cou: cou,
			Gau: gau,
			His: his,
		})
	}

	return &Handler{
		ec2: ec2.NewFromConfig(c.Aws),
		log: c.Log,
		reg: reg,
	}
}
