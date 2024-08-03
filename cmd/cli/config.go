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

func ConfigDir() string {
	return filepath.Dir(configPath)
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

	err := viper.ReadInConfig()
	if err == nil {
		return
	}

	_, ok := err.(viper.ConfigFileNotFoundError)
	if !ok {
		log.Error("Error reading config file: %v", err)
		return
	}

	err = os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		log.Error("Error creating directory paths for config", "err", err)
	}

	file, err := os.Create(configPath)
	if err != nil {
		log.Error("Error creating config file", "err", err)
		return
	}

	_ = file.Close()

	if err := viper.WriteConfig(); err != nil {
		log.Error("Error writing to config", "err", err)
		return
	}
}

func setDefaults(vpr *viper.Viper) {
	vpr.SetDefault("appId", "Robotics.StolenMowers.Service@husqvarnagroup.com")
	// Todo: Add sites defaults
}
