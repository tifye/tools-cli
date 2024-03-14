package cmd

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/Tifufu/tools-cli/pkg/security"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "login",
		Short:            "Authenticate user",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {}, // Ignore authentication check from root command
		RunE:             runLogin,
	}
	return cmd
}

func runLogin(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	user, err := security.AuthenticateUser(ctx, viper.GetString("appId"))
	if err != nil {
		return err
	}

	encrypted, err := security.EncryptUserProfile(user)
	if err != nil {
		return err
	}
	encryptedBase64 := base64.StdEncoding.EncodeToString(encrypted)
	viper.Set("user", encryptedBase64)

	log.Info("Authenticated as", "user", user.Profile.Email)
	return viper.WriteConfig()
}
