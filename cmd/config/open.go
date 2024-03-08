package config

import (
	"os/exec"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/spf13/cobra"
)

type openOptions struct {
	openInVSCode bool
}

func newOpenCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &openOptions{}

	openCmd := &cobra.Command{
		Use:   "open",
		Short: "Open a configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			filepath := cli.GetConfigPath()

			if opts.openInVSCode {
				err := exec.Command("code", filepath).Run()
				if err != nil {
					tCli.Log.Fatal("Error opening config in VSCode", "err", err)
				}
				return
			}

			err := pkg.OpenURL(filepath)
			if err != nil {
				tCli.Log.Fatal("Error opening config", "err", err)
			}
		},
	}

	openCmd.Flags().BoolVarP(&opts.openInVSCode, "code", "c", false, "Open in VSCode")

	return openCmd
}
