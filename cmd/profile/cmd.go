package profile

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func NewProfileCommand(toolsCli *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "User profile subcommands",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("Authenticated as", "user", toolsCli.User.Profile.Email)
			cmd.Help()
		},
	}

	cmd.AddCommand(newListCommand(toolsCli))

	return cmd
}
