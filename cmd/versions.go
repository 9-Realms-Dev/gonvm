package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of go_nvm",
	Long:  "Print the version number of go_nvm",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("v2.0.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
