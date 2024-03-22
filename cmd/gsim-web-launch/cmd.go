package gsim_web_launch

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Tifufu/tools-cli/cmd/cli"
	tifconsole "github.com/Tifufu/tools-cli/internal/tif-console"
	"github.com/Tifufu/tools-cli/internal/winmower"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type gsimWebLaunchOptions struct {
	serialNumber   uint
	platform       pkg.Platform
	tifConsolePath string
	detach         bool
}

func NewGsimWebLaunchCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &gsimWebLaunchOptions{}
	cmd := &cobra.Command{
		Use:   "gsim-web-launch",
		Short: "Launch the gsim web server",
		PreRun: func(cmd *cobra.Command, args []string) {
			if opts.tifConsolePath != "" {
				return
			}

			tCli.Log.Debug("No tifConsolePath provided, looking for default path...")
			cacheDir, err := os.UserCacheDir()
			if err != nil {
				tCli.Log.Error("Error getting user cache dir", "err", err)
				return
			}
			opts.tifConsolePath = filepath.Join(cacheDir, "TifApp", "TifConsole.Auto.exe")
			tCli.Log.Debug("Using default tifConsolePath", "tifConsolePath", opts.tifConsolePath)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if opts.detach {
				exePath, err := os.Executable()
				if err != nil {
					tCli.Log.Error("Error getting executable path", "err", err)
					return
				}

				err = exec.Command(
					"cmd", "/c", "start",
					exePath, "gsim-web-launch", args[0],
					"--tif-console", opts.tifConsolePath,
					"--debug").
					Run()
				if err != nil {
					tCli.Log.Error("Error detaching process", "err", err)
					return
				}
				return
			}

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

	cmd.Flags().StringVar(&opts.tifConsolePath, "tif-console", "", "Path to TifConsole.Auto.exe")
	cmd.MarkFlagFilename("tif-console", "exe")

	cmd.Flags().BoolVar(&opts.detach, "detach", false, "Detach the process")

	cmd.AddCommand(newUpdateRegistryCommand(tCli))

	return cmd
}

func runGsimWebLaunch(tCli *cli.ToolsCli, opts gsimWebLaunchOptions) {
	tCli.Log.Debug("Args", "serialNumber", opts.serialNumber, "platform", opts.platform)

	wm, err := tCli.WinMowerRegistry.DownloadWinMower(opts.platform, context.TODO())
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
	gspMeta, err := gspRegistry.DownloadGSPacket(opts.serialNumber, opts.platform, context.TODO())
	if err != nil {
		tCli.Log.Error("Error getting gspacket", "err", err)
		return
	}
	tCli.Log.Info("GSPacket", "gspacket", gspMeta.Map)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
		Prefix:          opts.platform.String(),
	})
	logger.SetStyles(cli.SubProcessLogStyle("#8b5cf6"))
	formatter := winmower.NewLogFormatter(logger)
	runner := winmower.RunnerContext(ctx, wm.Path)
	runner.SetStdout(formatter)
	runner.SetStderr(formatter)
	runner.SetWorkDir(filepath.Dir(wm.Path))
	err = runner.Start()
	if err != nil {
		tCli.Log.Error("Error starting winmower", "err", err)
		return
	}
	defer func() {
		err := runner.Stop()
		if err != nil {
			tCli.Log.Error("Error waiting for winmower", "err", err)
		}
	}()

	time.Sleep(3 * time.Second) // TODO: Wait for winmower to start by spying on its logs
	tifLogger := logger.WithPrefix("GSPacket bundle")
	tifLogger.SetStyles(cli.SubProcessLogStyle("#3b82f6"))
	tifLogFormatter := tifconsole.NewLogFormatter(tifLogger)
	tifConsole := &tifconsole.TifConsole{
		Path:   opts.tifConsolePath,
		Stdout: tifLogFormatter,
		Stderr: tifLogFormatter,
	}
	err = tifConsole.RunTestBundle(ctx, gspMeta.TestBundle, "-tcpAddress", "127.0.0.1:4250")
	if err != nil {
		tCli.Log.Error("Error running test bundle", "err", err)
		return
	}

	args := []string{
		"-config", gspMeta.Map,
		"-log", "false",
		"-time-scale", "1",
		"-screen-width", "1280",
		"-screen-height", "720",
		"-quality-level", "6",
	}
	cmd := exec.CommandContext(ctx, simMeta.Path, args...)
	err = cmd.Start()
	if err != nil {
		tCli.Log.Error("Error running simulator", "err", err)
		return
	}
	defer cmd.Wait() // Wait does more than just wait, it also cleans up the process
	defer func() {
		err := cmd.Process.Kill()
		if err != nil {
			tCli.Log.Error("Error killing simulator", "err", err)
		}
	}()

	// TODO: Different mowers have different testscripts
	logger = tifLogger.WithPrefix("Start script")
	tifLogFormatter = tifconsole.NewLogFormatter(logger)
	tifConsole.Stdout = tifLogFormatter
	tifConsole.Stderr = tifLogFormatter
	err = tifConsole.RunTestBundle(context.Background(), `D:\Projects\_work\_pocs\tools-cli\assets\testscript.zip`, "-tcpAddress", "127.0.0.1:4250")
	if err != nil {
		tCli.Log.Error("Error running test bundle", "err", err)
	}

	var input string
	fmt.Println("Press enter to exit...")
	fmt.Scanln(&input)
}
