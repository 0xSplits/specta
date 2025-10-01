package daemon

import (
	"context"
	"os"

	"github.com/0xSplits/specta/pkg/worker/handler/build"
	"github.com/0xSplits/specta/pkg/worker/handler/container"
	"github.com/0xSplits/specta/pkg/worker/handler/deployment"
	"github.com/0xSplits/specta/pkg/worker/handler/endpoint"
	"github.com/0xSplits/specta/pkg/worker/handler/keypair"
	"github.com/0xSplits/specta/pkg/worker/handler/stack"
	"github.com/0xSplits/workit/handler"
	"github.com/0xSplits/workit/registry"
	"github.com/0xSplits/workit/worker/parallel"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/google/go-github/v75/github"
	"github.com/xh3b4sd/tracer"
)

func (d *Daemon) Worker() *parallel.Worker {
	var cfg aws.Config
	{
		cfg = musAws()
	}

	var git *github.Client
	{
		git = newGit(d.env)
	}

	var reg *registry.Registry
	{
		reg = registry.New(registry.Config{
			Env: d.env.Environment,
			Log: d.log,
			Met: d.met,
		})
	}

	var wor *parallel.Worker
	{
		wor = parallel.New(parallel.Config{
			Han: []handler.Cooler{
				build.New(build.Config{Env: d.env, Git: git, Log: d.log, Met: d.met}),
				deployment.New(deployment.Config{Aws: cfg, Env: d.env, Log: d.log, Met: d.met}),
				container.New(container.Config{Aws: cfg, Env: d.env, Log: d.log, Met: d.met}),
				endpoint.New(endpoint.Config{Env: d.env, Log: d.log, Met: d.met}),
				keypair.New(keypair.Config{Aws: cfg, Env: d.env, Log: d.log, Met: d.met}),
				stack.New(stack.Config{Aws: cfg, Env: d.env, Log: d.log, Met: d.met}),
			},
			Log: d.log,
			Reg: reg,
		})
	}

	return wor
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
