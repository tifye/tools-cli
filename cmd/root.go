package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Tifufu/tools-cli/cmd/amprod"
	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/cmd/config"
	gsimLaunch "github.com/Tifufu/tools-cli/cmd/gsim-web-launch"
	"github.com/Tifufu/tools-cli/cmd/profile"
	"github.com/Tifufu/tools-cli/cmd/sites"
	winmower "github.com/Tifufu/tools-cli/cmd/win-mower"
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/Tifufu/tools-cli/pkg/security"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd *cobra.Command
	tCli    *cli.ToolsCli
)

type persistentOptions struct {
	configPath string
	logDebug   bool
}

var opts = &persistentOptions{}

func newRootCommand(_ *cli.ToolsCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Robotics tools CLI",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVar(&opts.configPath, "config", "", "config file (default is $UserCacheDir/tools-cli/confg.yaml)")
	cmd.PersistentFlags().BoolVar(&opts.logDebug, "debug", false, "cnable debug logging")

	return cmd
}

func init() {
	cobra.MousetrapHelpText = ""

	cobra.OnInitialize(
		func() {
			if opts.logDebug {
				tCli.Log.SetLevel(log.DebugLevel)
			}
			cli.SetConfigPath(opts.configPath)
		},
		cli.InitConfig,
		func() {
			user, err := decodeCachedAuth()
			if err != nil {
				tCli.Log.Debug("Error decoding cached auth", "err", err)
				return
			}

			tCli.Client.Transport = security.NewTifAuthTransport(http.DefaultTransport, user.APIKey, user.AccessToken)
		},
		func() {
			tCli.WinMowerRegistry = pkg.NewWinMowerRegistry(filepath.Join(cli.ConfigDir(), "winmowers"), tCli.BundleRegistry, tCli.Log)

			tCli.WinMowerRegistry.WithClient(*tCli.Client)
			tCli.BundleRegistry.WithClient(*tCli.Client)
		},
	)
}

func Execute() {
	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    false,
		ReportTimestamp: false,
	})
	client := &http.Client{}
	bundleRegistry := pkg.NewBundleRegistry("https://hqvrobotics.azure-api.net")
	tCli = &cli.ToolsCli{
		Log:            logger,
		BundleRegistry: bundleRegistry,
		Client:         client,
	}

	rootCmd = newRootCommand(tCli)
	addCommands(rootCmd, tCli)

	err := rootCmd.Execute()
	if err != nil {
		logger.Error("Error executing command", "err", err)
		logger.Info("Press enter to exit")
		var input string
		fmt.Scan(&input)
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
		winmower.NewWinMowerCommand(toolsCli),
		gsimLaunch.NewGsimWebLaunchCommand(toolsCli),
	)
}

func decodeCachedAuth() (*security.UserProfile, error) {
	encryptedUser := viper.GetString("user")
	if encryptedUser == "" {
		return nil, errors.New("user not authenticated. Please run 'tools login'")
	}
	decoded, err := base64.StdEncoding.DecodeString(encryptedUser)
	if err != nil {
		return nil, fmt.Errorf("error decoding user profile: %v", err)
	}
	user, err := security.DecryptUserProfile(decoded)
	if err != nil {
		return nil, fmt.Errorf("error decrypting user profile: %v", err)
	}

	return user, nil
}
