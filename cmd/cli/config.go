package cli

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

var configPath string

func SetConfigPath(path string) {
	configPath = path
}

func ConfigPath() string {
	return configPath
}

func InitConfig() {
	setDefaults(viper.GetViper())

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			log.Fatal(err)
		}
		configDir := filepath.Join(cacheDir, "tools-cli")
		viper.AddConfigPath(configDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		SetConfigPath(filepath.Join(configDir, "config.yaml"))
	}

	configPath = ConfigPath()
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := os.MkdirAll(configPath, 0755)
			if err != nil {
				log.Error(err)
			}
			if err := viper.WriteConfig(); err != nil {
				log.Error(err)
			}
		} else {
			log.Error("Error reading config file: %v", err)
		}
	}
}

func setDefaults(vpr *viper.Viper) {
	vpr.SetDefault("appId", "Robotics.StolenMowers.Service@husqvarnagroup.com")
	// Todo: Add sites defaults
}
