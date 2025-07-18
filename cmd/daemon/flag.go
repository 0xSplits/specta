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
	cmd.Flags().StringVar(&f.Env, "env", ".env", "the environment file to load")
}

func (f *flag) Validate() error {
	if f.Env == "" {
		return tracer.Mask(runtime.InvalidFlagError, tracer.Context{Key: "reason", Value: "--env must not be empty"})
	}

	return nil
}
