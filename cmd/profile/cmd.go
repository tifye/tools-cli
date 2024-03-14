package profile

import (
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/spf13/cobra"
)

func NewProfileCommand(toolsCli *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "User profile subcommands",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if toolsCli.User == nil {
				toolsCli.Log.Fatal("Must be authenticated to use profile commands. Run `tools auth login` to authenticate.")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			toolsCli.Log.Info("Authenticated as", "user", toolsCli.User.Profile.Email)
		},
	}

	cmd.AddCommand(newListCommand(toolsCli))

	return cmd
}
