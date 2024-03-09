package profile

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/spf13/cobra"
)

func NewProfileCommand(toolsCli *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "User profile subcommands",
		Run: func(cmd *cobra.Command, args []string) {
			toolsCli.Log.Info("Authenticated as", "user", toolsCli.User.Profile.Email)
		},
	}

	cmd.AddCommand(newListCommand(toolsCli))

	return cmd
}
