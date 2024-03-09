package cmd

import (
	"encoding/base64"
	"os"

	"github.com/Tifufu/tools-cli/cmd/amprod"
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/cmd/config"
	"github.com/Tifufu/tools-cli/cmd/profile"
	"github.com/Tifufu/tools-cli/cmd/sites"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd *cobra.Command
	tCli    *cli.ToolsCli
)

const encryptionKey string = "{7f8d534a-bf20-4e69-bbf8-54f4a9378f23}"

type persistentOptions struct {
	configPath string
	logDebug   bool
}

var opts = &persistentOptions{}

func newRootCommand(tCli *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Robotics tools CLI",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			encryptedUser := viper.GetString("user")
			if encryptedUser == "" {
				tCli.Log.Fatal("User not authenticated. Please run 'tools login'")
			}
			decoded, err := base64.StdEncoding.DecodeString(encryptedUser)
			if err != nil {
				tCli.Log.Fatal("Error decoding user profile", "err", err)
			}
			user, err := pkg.DecryptUserProfile(decoded)
			if err != nil {
				tCli.Log.Fatal("Error decrypting user profile", "err", err)
			}
			tCli.User = user
		},
	}

	cmd.PersistentFlags().StringVar(&opts.configPath, "config", "", "config file (default is $UserCacheDir/tools-cli/confg.yaml)")
	cmd.PersistentFlags().BoolVar(&opts.logDebug, "debug", false, "cnable debug logging")

	return cmd
}

func init() {
	cobra.OnInitialize(
		func() {
			if opts.logDebug {
				tCli.Log.SetLevel(log.DebugLevel)
			}
			cli.SetConfigPath(opts.configPath)
		},
		cli.InitConfig,
	)
}

func Execute() {
	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    false,
		ReportTimestamp: false,
	})
	tCli = &cli.ToolsCli{
		Log: logger,
	}
	rootCmd = newRootCommand(tCli)
	addCommands(rootCmd, tCli)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addCommands(cmd *cobra.Command, toolsCli *cli.ToolsCli) {
	cmd.AddCommand(
		newLoginCommand(),
		sites.NewSitesCommand(toolsCli),
		profile.NewProfileCommand(toolsCli),
		amprod.NewAmProdCommand(toolsCli),
		config.NewConfigCommand(toolsCli),
	)
}
