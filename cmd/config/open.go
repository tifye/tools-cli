package config

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/spf13/cobra"
)

func newOpenCommand(tCli *cli.ToolsCli) *cobra.Command {
	openCmd := &cobra.Command{
		Use:   "open",
		Short: "Open a configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			err := pkg.OpenURL(cli.GetConfigPath())
			if err != nil {
				tCli.Log.Fatal("Error opening config", "err", err)
			}
		},
	}

	return openCmd
}
