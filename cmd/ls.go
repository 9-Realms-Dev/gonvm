package cmd

import (
	"fmt"

	"github.com/9-Realms-Dev/gonvm/internal/nvm"
	"github.com/9-Realms-Dev/gonvm/internal/util"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List installed Node.js versions",
	Long:  "This command lists all Node.js versions installed locally",
	Args:  cobra.NoArgs,
	RunE:  ListInstalledVersions,
}

var lsRemoteCmd = &cobra.Command{
	Use:   "ls-remote",
	Short: "List available Node.js versions",
	Long:  "This command lists all Node.js versions available for installation",
	Args:  cobra.NoArgs,
	RunE:  ListRemoteVersions,
}

func init() {
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(lsRemoteCmd)
}

func ListInstalledVersions(cmd *cobra.Command, args []string) error {
	util.Logger.Info("Listing installed versions...")
	versionList, err := nvm.LocalVersions()
	if err != nil {
		return err
	}

	// TODO: Replace with tui list
	for _, v := range versionList {
		fmt.Println(v)
	}

	return nil
}

func ListRemoteVersions(cmd *cobra.Command, args []string) error {
	util.Logger.Info("Listing remote versions...")
	versionList, err := nvm.RemoteVersions()
	if err != nil {
		return err
	}

	// TODO: Replace with tui list
	for _, v := range versionList {
		fmt.Println(v)
	}

	return nil
}
