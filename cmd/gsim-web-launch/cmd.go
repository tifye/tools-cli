package gsim_web_launch

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/spf13/cobra"
)

type gsimWebLaunchOptions struct {
	serialNumber uint
	platform     pkg.Platform
}

func NewGsimWebLaunchCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &gsimWebLaunchOptions{}
	cmd := &cobra.Command{
		Use:   "gsim-web-launch",
		Short: "Launch the gsim web server",
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				var input string
				fmt.Println("Press enter to exit...")
				fmt.Scanln(&input)
			}()

			if len(args) == 0 {
				tCli.Log.Error("Excepted args: {serialNumber}/{platform}")
				return
			}

			payload := strings.TrimPrefix(args[0], "gsim-web-launch:")
			parts := strings.Split(payload, "/")
			if len(parts) != 2 {
				tCli.Log.Error("Invalid payload format, expected '{serialNumber}/{platform}'", "payload", payload)
				return
			}

			serialNumber, err := strconv.ParseUint(parts[0], 10, 32)
			if err != nil {
				tCli.Log.Error("Error parsing serial number", "serialNumber", parts[0], "err", err)
				return
			}
			var platform pkg.Platform
			err = platform.Set(parts[1])
			if err != nil {
				tCli.Log.Error("Error setting platform", "platform", parts[1], "err", err, "valid platforms", pkg.GetPlatforms())
				return
			}

			opts.serialNumber = uint(serialNumber)
			opts.platform = platform

			runGsimWebLaunch(tCli, *opts)
		},
	}

	cmd.AddCommand(newUpdateRegistryCommand(tCli))

	return cmd
}

func runGsimWebLaunch(tCli *cli.ToolsCli, opts gsimWebLaunchOptions) {
	tCli.Log.Debug("Args", "serialNumber", opts.serialNumber, "platform", opts.platform)

	_, err := tCli.WinMowerRegistry.DownloadWinMower(opts.platform, context.TODO())
	if err != nil {
		tCli.Log.Error("Error getting winmower", "err", err)
		return
	}
	tCli.Log.Debug("Winmower", "platform", opts.platform)

	simRegistry := pkg.NewSimulatorRegistry(filepath.Join(cli.ConfigDir(), "simulators"), tCli.BundleRegistry, tCli.Client, tCli.Log)
	simMeta, err := simRegistry.DownloadSimulator(context.TODO())
	if err != nil {
		tCli.Log.Error("Error getting simulator", "err", err)
		return
	}
	tCli.Log.Debug("Simulator", "simulator", simMeta.Path)

	baseUrl := "https://hqvrobotics.azure-api.net/gardensimulatorpacket"
	gspRegistry := pkg.NewGSPacketRegistry(filepath.Join(cli.ConfigDir(), "gspackets"), baseUrl, tCli.Client, tCli.Log)
	gspMeta, err := gspRegistry.DownloadGSPacket(opts.serialNumber, opts.platform, context.Background())
	if err != nil {
		tCli.Log.Error("Error getting gspacket", "err", err)
		return
	}
	tCli.Log.Info("GSPacket", "gspacket", gspMeta.Map)
}
