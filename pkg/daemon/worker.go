package daemon

import (
	"context"
	"os"

	"github.com/0xSplits/specta/pkg/worker"
	"github.com/0xSplits/specta/pkg/worker/handler"
	"github.com/0xSplits/specta/pkg/worker/handler/container"
	"github.com/0xSplits/specta/pkg/worker/handler/endpoint"
	"github.com/0xSplits/specta/pkg/worker/handler/keypair"
	"github.com/0xSplits/specta/pkg/worker/handler/stack"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/xh3b4sd/tracer"
)

func (d *Daemon) Worker() *worker.Worker {
	var cfg aws.Config
	{
		cfg = musAws()
	}

	return worker.New(worker.Config{
		Env: d.env,
		Han: []handler.Interface{
			container.New(container.Config{Aws: cfg, Env: d.env, Log: d.log, Met: d.met}),
			endpoint.New(endpoint.Config{Env: d.env, Log: d.log, Met: d.met}),
			keypair.New(keypair.Config{Aws: cfg, Env: d.env, Log: d.log, Met: d.met}),
			stack.New(stack.Config{Aws: cfg, Env: d.env, Log: d.log, Met: d.met}),
		},
		Log: d.log,
		Met: d.met,
	})
}

func musAws() aws.Config {
	reg := os.Getenv("AWS_REGION")
	if reg == "" {
		reg = "us-west-2"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(reg))
	if err != nil {
		tracer.Panic(tracer.Mask(err))
	}

	return cfg
}
