package cmd

import (
	"encoding/base64"
	"os"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/cmd/profile"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd *cobra.Command

const encryptionKey string = "{7f8d534a-bf20-4e69-bbf8-54f4a9378f23}"

func newRootCommand(toolsCli *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Robotics tools CLI",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			encryptedUser := viper.GetString("user")
			if encryptedUser == "" {
				log.Fatal("User not authenticated. Please run 'tools login'")
			}
			decoded, err := base64.StdEncoding.DecodeString(encryptedUser)
			if err != nil {
				log.Fatal("Error decoding user profile", "err", err)
			}
			user, err := pkg.DecryptUserProfile(decoded)
			if err != nil {
				log.Fatal("Error decrypting user profile", "err", err)
			}
			toolsCli.User = user
		},
	}
	return cmd
}

func Execute() {
	initConfig()

	toolsCli := &cli.ToolsCli{}
	rootCmd = newRootCommand(toolsCli)
	addCommands(rootCmd, toolsCli)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addCommands(cmd *cobra.Command, toolsCli *cli.ToolsCli) {
	cmd.AddCommand(
		newLoginCommand(),
		profile.NewProfileCommand(toolsCli),
	)
}
