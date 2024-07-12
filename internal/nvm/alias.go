package nvm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/9-Realms-Dev/go_nvm/internal/util"
	"github.com/spf13/viper"
)

func GetAliasedVersion(name string) (string, error) {
	// Check if GO_NVM_DIR is set
	nvmDir, err := util.GetNvmDirectory()
	if nvmDir == "" {
		return "", fmt.Errorf("GO_NVM_DIR environment variable is not set")
	}

	// Load aliases
	aliases, err := loadAliases()
	if err != nil {
		return "", fmt.Errorf("error loading aliases: %w", err)
	}

	// If aliases is nil, it means the file doesn't exist or is empty
	if aliases == nil {
		return "", nil
	}

	// Return the aliased version if it exists
	if version, ok := aliases[name]; ok {
		return version, nil
	}

	// If the alias doesn't exist, return an empty string and nil error
	return "", nil
}

func SetAliasedVersion(name, version string) error {
	// Check if the nvm directory is set
	_, err := util.GetNvmDirectory()
	if err != nil {
		return fmt.Errorf("error getting NVM directory: %w", err)
	}

	// Load existing aliases
	aliases, err := loadAliases()
	if err != nil {
		return fmt.Errorf("error loading aliases: %w", err)
	}

	// If aliases is nil, initialize it
	if aliases == nil {
		aliases = make(map[string]string)
	}

	// Set the new alias
	aliases[name] = version

	// Save the updated aliases
	if err := saveAliases(aliases); err != nil {
		return fmt.Errorf("error saving aliases: %w", err)
	}

	util.Logger.Infof("Created alias %s for %s\n", name, version)
	return nil
}

func RemoveAlias(name string) error {
	// Get the NVM directory
	_, err := util.GetNvmDirectory()
	if err != nil {
		return fmt.Errorf("error getting NVM directory: %w", err)
	}

	// Load existing aliases
	aliases, err := loadAliases()
	if err != nil {
		return fmt.Errorf("error loading aliases: %w", err)
	}

	// If aliases is nil, it means no aliases were found
	if aliases == nil {
		fmt.Println("No aliases found")
		return nil
	}

	// Check if the alias exists
	if _, exists := aliases[name]; exists {
		// Remove the alias
		delete(aliases, name)

		// Save the updated aliases
		if err := saveAliases(aliases); err != nil {
			return fmt.Errorf("error saving aliases after removal: %w", err)
		}

		fmt.Printf("Removed alias %s\n", name)
	} else {
		fmt.Printf("Alias %s not found\n", name)
	}

	return nil
}

func loadAliases() (map[string]string, error) {
	// Get the GO_NVM_DIR from environment variables
	nvmDir, err := util.GetNvmDirectory()
	if err != nil {
		return nil, err
	}

	// Construct the path to alias.toml
	filePath := filepath.Join(nvmDir, "alias.toml")

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil // Return nil map and nil error if file doesn't exist
	}

	// Create a new Viper instance
	v := viper.New()

	// Set the configuration file path
	v.SetConfigName("alias")
	v.SetConfigType("toml")
	v.AddConfigPath(nvmDir)

	// Read the configuration file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; return nil map and nil error
			return nil, nil
		}
		// Config file was found but another error was produced
		return nil, err
	}

	// Get all settings as a map
	allSettings := v.AllSettings()

	// Convert the map to map[string]string
	aliases := make(map[string]string)
	for key, value := range allSettings {
		if strValue, ok := value.(string); ok {
			aliases[key] = strValue
		}
	}

	return aliases, nil
}

func saveAliases(newAliases map[string]string) error {
	// Get the PY_NVM_DIR from environment variables
	nvmDir, err := util.GetNvmDirectory()
	if err != nil {
		return err
	}

	// Create a new Viper instance
	v := viper.New()

	// Set the configuration file path
	v.SetConfigName("alias")
	v.SetConfigType("toml")
	v.AddConfigPath(nvmDir)

	// Read the existing configuration file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Return error if it's not a "file not found" error
			return fmt.Errorf("error reading config file: %w", err)
		}
		// If file not found, we'll create a new one
	}

	// Update the aliases
	for key, value := range newAliases {
		v.Set(key, value)
	}

	// Ensure the directory exists
	if err := os.MkdirAll(nvmDir, 0755); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	// Save the updated configuration
	if err := v.WriteConfigAs(filepath.Join(nvmDir, "alias.toml")); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}
