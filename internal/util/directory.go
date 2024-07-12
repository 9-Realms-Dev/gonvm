package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetNvmDirectory() (string, error) {
	nvmDir := os.Getenv("GO_NVM_DIR")
	if nvmDir == "" {
		return SetDefaultDirectory()
	}

	return nvmDir, nil
}

func SetDefaultDirectory() (string, error) {
	nvmDir := os.Getenv("GO_NVM_DIR")
	if nvmDir == "" {
		Logger.Warn("GO_NVM_DIR value is not set. Creating default directory...")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		nvmDir = filepath.Join(homeDir, ".go_nvm")

		// Create the directory if it doesn't exist
		if err := os.MkdirAll(nvmDir, 0755); err != nil {
			return "", err
		}
		os.Setenv("GO_NVM_DIR", nvmDir)
		Logger.Infof("Created directory %s and set GO_NVM_DIR environment variable", nvmDir)
	}

	// Create the alias.json file if it doesn't exist
	aliasFile := filepath.Join(nvmDir, "alias.toml")

	v := viper.New()
	v.SetConfigFile(aliasFile)
	v.SetConfigType("toml")
	v.AddConfigPath(nvmDir)

	if _, err := os.Stat(aliasFile); os.IsNotExist(err) {
		defaultConfig := map[string]string{
			"latest": "",
			"lts":    "",
		}
		for key, value := range defaultConfig {
			v.SetDefault(key, value)
		}
		if err := v.SafeWriteConfig(); err != nil {
			return "", fmt.Errorf("unable to write default config: %v", err)
		}
		Logger.Infof("Created alias file %s", aliasFile)
	}

	return nvmDir, nil
}
