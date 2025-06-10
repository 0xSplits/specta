package daemon

import (
	"github.com/0xSplits/specta/pkg/runtime"
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"
)

type flag struct {
	Env string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.Env, "env", "development", "the environment file to load, e.g. development for env.development")
}

func (f *flag) Validate() error {
	if f.Env != "development" && f.Env != "testing" && f.Env != "staging" && f.Env != "production" {
		return tracer.Maskf(runtime.ExecutionFailedError, "--env must be one of [development testing staging production]")
	}

	return nil
}
