package cli

import (
	"os"
	"path"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func InitConfig() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}

	setDefaults(viper.GetViper())

	configDir := path.Join(cacheDir, "tools-cli")
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := os.MkdirAll(configDir, 0755)
			if err != nil {
				log.Error(err)
			}
			if err := viper.WriteConfigAs(path.Join(configDir, "config.yaml")); err != nil {
				log.Error(err)
			}
		} else {
			log.Error("Error reading config file: %v", err)
		}
	}
}

func setDefaults(vpr *viper.Viper) {
	vpr.SetDefault("appId", "Robotics.StolenMowers.Service@husqvarnagroup.com")
}
