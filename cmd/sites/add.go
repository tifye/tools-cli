package sites

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type addOptions struct {
	name     string
	tagLinks map[string]string
}

func newAddCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &addOptions{}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a site",
		Run: func(cmd *cobra.Command, args []string) {
			_, ok := sites[opts.name]
			if ok {
				log.Print("Site already exists", "site", opts.name)
				return
			}

			sites[opts.name] = opts.tagLinks
		},
	}

	cmd.Flags().StringVarP(&opts.name, "site", "s", "", "The site name")
	cmd.MarkFlagRequired("site")

	cmd.Flags().StringToStringVarP(&opts.tagLinks, "tags", "t", nil, "Comma separated tag=link pairs")
	cmd.MarkFlagRequired("tag")

	return cmd
}
