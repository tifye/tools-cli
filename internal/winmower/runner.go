package winmower

import (
	"context"
	"io"
	"os/exec"
	"syscall"
)

type winmowerRunner struct {
	cmd *exec.Cmd
}

func RunnerContext(ctx context.Context, wmPath string) *winmowerRunner {
	cmd := exec.CommandContext(ctx, wmPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return &winmowerRunner{
		cmd: cmd,
	}
}

func (w *winmowerRunner) Start() error {
	return w.cmd.Start()
}

func (w *winmowerRunner) Stop() error {
	err := w.cmd.Cancel()
	if err != nil {
		return err
	}

	return w.cmd.Wait()
}

func (w *winmowerRunner) SetWorkDir(dir string) {
	w.cmd.Dir = dir
}

func (w *winmowerRunner) SetStdin(stdin io.Reader) {
	w.cmd.Stdin = stdin
}

func (w *winmowerRunner) SetStdout(stdout io.Writer) {
	w.cmd.Stdout = stdout
}

func (w *winmowerRunner) SetStderr(stderr io.Writer) {
	w.cmd.Stderr = stderr
}
