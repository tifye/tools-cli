package cli

import (
	"os"
	"path"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func GetConfigPath() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(cacheDir, "tools-cli")
}

func InitConfig() {
	setDefaults(viper.GetViper())

	configDir := GetConfigPath()
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
