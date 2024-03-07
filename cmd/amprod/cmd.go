package amprod

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/spf13/cobra"
)

func NewAmProdCommand(toolsCli *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "amprod",
		Short: "AM Production subcommands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(newValidateTestIndexCommand(toolsCli))

	return cmd
}
