package sites

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/spf13/cobra"
)

func newListCommand(tCli *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sites",
		Run: func(cmd *cobra.Command, args []string) {
			for site, tagLinks := range sites {
				tempLogger := tCli.Log.WithPrefix(site)
				for tag, link := range tagLinks {
					tempLogger.Printf("%s = %s", tag, link)
				}
			}
		},
	}

	return cmd
}
