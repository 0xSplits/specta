package cmd

import (
	"github.com/0xSplits/specta/cmd/daemon"
	"github.com/0xSplits/specta/cmd/version"
	"github.com/spf13/cobra"
)

var (
	use = "specta"
	sho = "Golang based RPC microservice."
	lon = "Golang based RPC microservice."
)

func New() *cobra.Command {
	var cmd *cobra.Command
	{
		cmd = &cobra.Command{
			Use:   use,
			Short: sho,
			Long:  lon,
			Run:   (&run{}).run,
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
			// We slience errors because we do not want to see spf13/cobra printing.
			// The errors returned by the commands will be propagated to the main.go
			// anyway, where we have custom error printing for the command line
			// tool.
			SilenceErrors: true,
			SilenceUsage:  true,
		}
	}

	{
		cmd.SetHelpCommand(&cobra.Command{Hidden: true})
	}

	{
		cmd.AddCommand(daemon.New())
		cmd.AddCommand(version.New())
	}

	return cmd
}
