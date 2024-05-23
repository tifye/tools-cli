package device

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/spf13/cobra"
)

func NewDeviceCommand(tCli *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "device",
		Short: "Device commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(newOpenCommand(tCli))

	return cmd
}
