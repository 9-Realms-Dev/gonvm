package cmd

import (
	"github.com/9-Realms-Dev/gonvm/internal/tui"
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	userLicense string

	// globally flags
	latestFlag    bool = false
	acceptAllFlag bool = false

	rootCmd = &cobra.Command{
		Use:   "gonvm",
		Short: "Go based nvm",
		Long:  "Go base node version manager with tui",
		Run:   ActivateTui,
	}
)

func ActivateTui(cmd *cobra.Command, args []string) {
	tui.Dashboard()
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// TODO: Read from config file and setup global settings
}
