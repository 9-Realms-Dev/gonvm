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
	}

	// Create the alias.json file if it doesn't exist
	aliasFile := filepath.Join(nvmDir, "alias.toml")

	v := viper.New()
	v.SetConfigName("alias")
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
	}

	return nvmDir, nil
}
