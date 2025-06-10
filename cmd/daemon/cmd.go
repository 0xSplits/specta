package daemon

import (
	"github.com/spf13/cobra"
)

const (
	use = "daemon"
	sho = "Execute Specta's long running process for exposing RPC handlers."
	lon = "Execute Specta's long running process for exposing RPC handlers."
)

func New() *cobra.Command {
	var flg *flag
	{
		flg = &flag{}
	}

	var cmd *cobra.Command
	{
		cmd = &cobra.Command{
			Use:   use,
			Short: sho,
			Long:  lon,
			RunE:  (&run{flag: flg}).runE,
		}
	}

	{
		flg.Init(cmd)
	}

	return cmd
}
