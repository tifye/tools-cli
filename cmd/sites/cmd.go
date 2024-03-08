package sites

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type siteLinks map[string]string

var sites map[string]siteLinks

type sitesOptions struct {
	site string
	tag  string
}

func NewSitesCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &sitesOptions{}

	cmd := &cobra.Command{
		Use:   "sites",
		Short: "Access to various tool sites",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			sites = make(map[string]siteLinks)
			err := viper.UnmarshalKey("sites", &sites)
			if err != nil {
				tCli.Log.Fatal("Error reading sites from config", "err", err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(sites) == 0 {
				tCli.Log.Print("No sites registered.")
				tCli.Log.Print("Add a site with", "command", "tools-cli add site --site <site> --url <url>")
				return
			}

			links, ok := sites[opts.site]
			if !ok {
				tCli.Log.Print("Site not found.", "site", opts.site)
				tCli.Log.Print("View all sites with", "command", "tools-cli sites list")
				return
			}

			link, ok := links[opts.tag]
			if !ok {
				tCli.Log.Print("Tag not found.", "tag", opts.tag)
				tCli.Log.Print("View all tags with", "command", "tools-cli sites list --site <site>")
				return
			}

			tCli.Log.Print("Opening", "link", link)
			pkg.OpenURL(link)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			viper.Set("sites", sites)
			err := viper.WriteConfig()
			if err != nil {
				tCli.Log.Fatal("Error saving to config", "err", err)
			}
		},
	}

	cmd.Flags().StringVarP(&opts.site, "site", "s", "", "Open link corresponding to site")
	cmd.MarkFlagRequired("site")
	cmd.Flags().StringVarP(&opts.tag, "tag", "t", "", "Open link corresponding to tag")
	cmd.MarkFlagRequired("tag") // Todo: Moke this optional, if optional then list all links and use charm form to select

	// Todo: Register for autocompletion

	cmd.AddCommand(
		newAddCommand(tCli),
		newListCommand(tCli),
	)

	return cmd
}
