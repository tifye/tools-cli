package config

import (
	"os/exec"
	"path/filepath"

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
			fpath := cli.ConfigPath()

			if opts.openInVSCode {
				tCli.Log.Debug("Opening config in vscode", "path", fpath)
				err := exec.Command("code", fpath).Run()
				if err != nil {
					tCli.Log.Fatal("Error opening config in VSCode", "err", err)
				}
				return
			}

			dir := filepath.Dir(fpath)
			tCli.Log.Debug("Opening config directory", "dir", dir)
			err := pkg.OpenURL(dir)
			if err != nil {
				tCli.Log.Fatal("Error opening config", "err", err)
			}
		},
	}

	openCmd.Flags().BoolVarP(&opts.openInVSCode, "code", "c", false, "Open in VSCode")

	return openCmd
}
