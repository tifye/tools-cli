package winmower

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/spf13/cobra"
)

type downloadOptions struct {
	platform pkg.Platform
}

func newDownloadCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &downloadOptions{}
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download winmower",
		Run: func(cmd *cobra.Command, args []string) {
			winMower, err := tCli.WinMowerRegistry.GetWinMower(opts.platform, cmd.Context())
			if err != nil {
				tCli.Log.Error("Error getting winmower", "err", err)
				return
			}

			tCli.Log.Info("Winmower", "winmower", winMower.Path)
		},
	}

	cmd.Flags().VarP(&opts.platform, "platform", "p", "Platform to download")
	cmd.MarkFlagRequired("platform")

	return cmd
}
