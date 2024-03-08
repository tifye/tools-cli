package sites

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type siteDetails struct {
	Prod    string
	Staging string
}

type sitesOptions struct {
	site        string
	openStaging bool
}

func NewSitesCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &sitesOptions{}
	var sites = map[string]siteDetails{}

	cmd := &cobra.Command{
		Use:   "sites",
		Short: "Access to various tool sites",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			sites = make(map[string]siteDetails)
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

			site, ok := sites[opts.site]
			if !ok {
				tCli.Log.Print("Site not found.", "site", opts.site)
				tCli.Log.Print("View all sites with", "command", "tools-cli sites list")
				return
			}

			var url string
			if opts.openStaging {
				url = site.Staging
			} else {
				url = site.Prod
			}

			if url == "" {
				tCli.Log.Print("No URL found for site", "site", opts.site, "staging", opts.openStaging)
				return
			}

			tCli.Log.Print("Opening site", "site", opts.site, "url", url)
			pkg.OpenURL(url)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			viper.Set("sites", sites)
		},
	}

	cmd.Flags().StringVarP(&opts.site, "site", "s", "", "The site to open")
	cmd.MarkFlagRequired("site")

	cmd.Flags().BoolVar(&opts.openStaging, "staging", false, "Open staging site")

	//Todo: Register for autocompletion

	cmd.AddCommand(newAddCommand())

	return cmd
}
