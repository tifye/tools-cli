package tifconsole

import (
	"context"
	"io"
	"os/exec"
)

type TifConsole struct {
	Path   string
	Stdout io.Writer
	Stderr io.Writer
}

func (tc *TifConsole) RunTestBundle(ctx context.Context, bundlePath string, args ...string) error {
	cmdArgs := append([]string{bundlePath}, args...)
	cmd := exec.CommandContext(ctx, tc.Path, cmdArgs...)
	cmd.Stdout = tc.Stdout
	cmd.Stderr = tc.Stderr
	return cmd.Run()
}
