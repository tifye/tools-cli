package gsim_web_launch

import (
	"fmt"
	"os"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/registry"
)

type updateRegistryOptions struct {
	exePath string
}

func newUpdateRegistryCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &updateRegistryOptions{}
	cmd := &cobra.Command{
		Use:   "update-registry",
		Short: "Update the gsim web server registry",
		Run: func(cmd *cobra.Command, args []string) {
			if opts.exePath == "" {
				exePath, err := os.Executable()
				if err != nil {
					tCli.Log.Error("Error getting executable path", "err", err)
					return
				}
				opts.exePath = exePath
			}

			if _, err := os.Stat(opts.exePath); os.IsNotExist(err) {
				tCli.Log.Error("Executable does not exist", "path", opts.exePath)
				return
			}

			// TODO: document
			pKey, _, err := registry.CreateKey(registry.CLASSES_ROOT, "gsim-web-launch", registry.ALL_ACCESS)
			if err != nil {
				tCli.Log.Error("Error creating registry key", "err", err)
				return
			}
			defer pKey.Close()
			pKey.SetStringValue("", "URL: GSim Web Launch Protocol")
			pKey.SetStringValue("URL Protocol", "")

			cKey, _, err := registry.CreateKey(pKey, "shell\\open\\command", registry.ALL_ACCESS)
			if err != nil {
				tCli.Log.Error("Error creating registry key", "err", err)
				return
			}
			defer cKey.Close()
			cKey.SetStringValue("", fmt.Sprintf(`"%s" "gsim-web-launch" "%%1" "--debug"`, opts.exePath))

			tCli.Log.Debug("Updated registry to", "path", opts.exePath)
		},
	}

	cmd.Flags().StringVarP(&opts.exePath, "exe", "e", "", "Path to the cli executable")
	cmd.MarkFlagFilename("exe", "exe")

	return cmd
}
