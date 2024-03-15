package winmower

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type startOptions struct {
	platform pkg.Platform
	detach   bool
	showRaw  bool
	// TODO: Add flag for working directory
}

type winMowerLogger struct {
	logger *log.Logger
}

// TODO: Handle working directories properly
// by default should be in config folder
func newStartCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &startOptions{}
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the WinMower",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStart(tCli, opts, cmd)
		},
	}

	cmd.Flags().VarP(&opts.platform, "platform", "p", "Platform to start")
	cmd.MarkFlagRequired("platform")

	cmd.Flags().BoolVarP(&opts.showRaw, "raw", "r", false, "Show raw output")
	cmd.Flags().BoolVarP(&opts.detach, "detach", "d", false, "Detach from the process")

	return cmd
}

func runStart(tCli *cli.ToolsCli, opts *startOptions, cmd *cobra.Command) error {
	winMower, err := tCli.WinMowerRegistry.GetWinMower(opts.platform, cmd.Context())
	if err != nil {
		tCli.Log.Error("Error getting winmower", "err", err)
		return err
	}

	if opts.detach {
		tCli.Log.Info("Starting winmower in detached mode", "platform", opts.platform.String(), "raw", opts.showRaw)
		if opts.showRaw {
			err = pkg.OpenURL(winMower.Path)
			if err != nil {
				tCli.Log.Error("Error opening winmower", "err", err)
				return err
			}
		} else {
			err := exec.Command(
				"cmd", "/c", "start",
				"tools-cli", "winmower", "start", // TODO: Instead use location of current executable
				"-p", opts.platform.String(),
				"--debug").
				Run()
			if err != nil {
				tCli.Log.Error("Error detaching winmower", "err", err)
				return err
			}
		}
		return nil
	}

	exCmd := exec.CommandContext(cmd.Context(), winMower.Path)
	exCmd.Dir = filepath.Dir(winMower.Path)
	exCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
	exCmd.Stdin = os.Stdin

	if opts.showRaw {
		exCmd.Stdout = os.Stdout
		exCmd.Stderr = os.Stderr
	} else {
		logWriter := constructWinMowerLogger()
		exCmd.Stdout = logWriter
		exCmd.Stderr = logWriter
	}

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
}

func (w *winMowerLogger) Write(data []byte) (int, error) {
	b := bytes.TrimSuffix(data, []byte("\n"))
	lines := bytes.Split(b, []byte("\n"))

	first := lines[0]
	if first[0] == byte('[') {
		switch {
		case bytes.Contains(first, []byte("ERROR")):
			w.logger.Error(string(trimWinMowerLogLine(first)))
		case bytes.Contains(first, []byte("WARNING")):
			w.logger.Warn(string(trimWinMowerLogLine(first)))
		case bytes.Contains(first, []byte("INFO")):
			w.logger.Info(string(trimWinMowerLogLine(first)))
		case bytes.Contains(first, []byte("DEBUG")):
			w.logger.Debug(string(trimWinMowerLogLine(first)))
		default:
			w.logger.Info(string(trimWinMowerLogLine(first)))
		}
	}

	rest := lines[1:]
	for _, line := range rest {
		w.Write(line)
	}

	return len(data), nil
}

func trimWinMowerLogLine(line []byte) []byte {
	const logLevelPrefixLen = 9 // winmower prefixes log level
	return line[logLevelPrefixLen:]
}

func constructWinMowerLogger() *winMowerLogger {
	logWriter := winMowerLogger{
		logger: log.NewWithOptions(os.Stdout, log.Options{
			ReportCaller:    false,
			ReportTimestamp: true,
			TimeFormat:      time.TimeOnly,
		}),
	}
	color := lipgloss.Color("#8b5cf6")
	style := log.DefaultStyles()
	style.Prefix = lipgloss.NewStyle().Foreground(color)
	style.Timestamp = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderLeftBackground(color).
		BorderLeftForeground(color).
		BorderLeft(true).
		PaddingLeft(1)
	logWriter.logger.SetStyles(style)
	return &logWriter
}
