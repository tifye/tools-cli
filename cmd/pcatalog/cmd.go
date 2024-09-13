package pcatalog

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/cmd/pcatalog/list"
	"github.com/spf13/cobra"
)

func NewProductCatalogCommand(tCli *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "Product catalog commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newDownloadCommand(tCli),
		list.NewListCommand(tCli),
	)

	return cmd
}
