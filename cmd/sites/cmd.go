package sites

import (
	"strings"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type siteLinks map[string]string

var sites map[string]siteLinks

type sitesOptions struct {
	siteTag string
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
			runSites(tCli, opts)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			viper.Set("sites", sites)
			err := viper.WriteConfig()
			if err != nil {
				tCli.Log.Fatal("Error saving to config", "err", err)
			}
		},
	}

	cmd.Flags().StringVarP(&opts.siteTag, "site", "s", "", "Open link corresponding to site")
	cmd.MarkFlagRequired("site")

	// Todo: Add base tag command
	// Todo: Register for autocompletion

	cmd.AddCommand(
		newAddCommand(tCli),
		newListCommand(tCli),
		newRemoveCommand(tCli),
	)

	return cmd
}

func runSites(tCli *cli.ToolsCli, opts *sitesOptions) {
	if len(sites) == 0 {
		tCli.Log.Error("No sites registered.")
		tCli.Log.Info("Add a site with", "command", "tools-cli add site --site <site> --url <url>")
		return
	}

	parts := strings.Split(opts.siteTag, ":")
	site := parts[0]
	tag := ""
	if len(parts) > 1 {
		tag = parts[1]
	}

	links, ok := sites[site]
	if !ok {
		tCli.Log.Error("Site not found.", "site", site)
		tCli.Log.Info("View all sites with", "command", "tools-cli sites list")
		return
	}

	link, ok := links[tag]
	if !ok {
		tCli.Log.Error("Tag not found.", "tag", tag)
		tCli.Log.Info("View all tags with", "command", "tools-cli sites list --site <site>")
		return
	}

	tCli.Log.Print("Opening", "link", link)
	pkg.OpenURL(link)
}
