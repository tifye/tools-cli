package tifdefinition

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/spf13/cobra"
)

func NewTifDefinitionCommand(tCli *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tifdef",
		Short: "Tif definition subcommands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newListCommand(tCli),
	)

	return cmd
}
