package daemon

import (
	"context"
	"os"

	"github.com/0xSplits/specta/pkg/worker"
	"github.com/0xSplits/specta/pkg/worker/handler"
	"github.com/0xSplits/specta/pkg/worker/handler/keypair"
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
		Han: []handler.Interface{
			keypair.New(keypair.Config{Aws: cfg, Log: d.log, Met: d.met}),
		},
		Log: d.log,
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
