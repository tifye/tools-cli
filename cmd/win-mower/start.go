package winmower

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/spf13/cobra"
)

type startOptions struct {
	platform pkg.Platform
	detach   bool
}

func newStartCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &startOptions{}
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the WinMower",
		RunE: func(cmd *cobra.Command, args []string) error {
			winMower, err := tCli.WinMowerRegistry.GetWinMower(opts.platform, cmd.Context())
			if err != nil {
				tCli.Log.Error("Error getting winmower", "err", err)
				return err
			}

			if opts.detach {
				err = pkg.OpenURL(winMower.Path)
				if err != nil {
					tCli.Log.Error("Error opening winmower", "err", err)
					return err
				}
				return nil
			}

			exCmd := exec.CommandContext(cmd.Context(), winMower.Path)
			exCmd.Dir = filepath.Dir(winMower.Path)
			exCmd.Stdout = os.Stdout
			exCmd.Stderr = os.Stderr
			exCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
			err = exCmd.Start()
			if err != nil {
				tCli.Log.Error("Error starting winmower", "err", err)
				return err
			}
			err = exCmd.Wait()
			if err != nil {
				tCli.Log.Error("Error waiting for winmower", "err", err)
				return err
			}

			return nil
		},
	}

	cmd.Flags().VarP(&opts.platform, "platform", "p", "Platform to start")
	cmd.MarkFlagRequired("platform")

	cmd.Flags().BoolVarP(&opts.detach, "detach", "d", false, "Detach from the process")

	return cmd
}
