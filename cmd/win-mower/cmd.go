package winmower

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/spf13/cobra"
)

func NewWinMowerCommand(tCli *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "winmower",
		Short: "Winmower subcommands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newDownloadCommand(tCli),
		newStartCommand(tCli),
		// TODO: Add list flag
	)

	return cmd
}
