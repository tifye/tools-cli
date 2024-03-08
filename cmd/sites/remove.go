package sites

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/spf13/cobra"
)

type removeOptions struct {
	name string
}

func newRemoveCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &removeOptions{}
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a site",
		Run: func(cmd *cobra.Command, args []string) {
			_, ok := sites[opts.name]
			if !ok {
				tCli.Log.Print("Site does not exist", "site", opts.name)
				return
			}

			delete(sites, opts.name)
		},
	}

	cmd.Flags().StringVarP(&opts.name, "site", "s", "", "The site name")
	cmd.MarkFlagRequired("site")

	return cmd
}
