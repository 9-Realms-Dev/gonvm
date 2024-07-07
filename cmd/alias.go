package cmd

import (
	"github.com/9-Realms-Dev/go_nvm/internal/nvm"
	"github.com/9-Realms-Dev/go_nvm/internal/util"
	"github.com/spf13/cobra"
)

var aliasCmd = &cobra.Command{
	Use:   "go_nvm alias [name] [version]",
	Short: "Create an alias for a Node.js version",
	Long:  "This command creates an alias for a specific Node.js version",
	Args:  cobra.ExactArgs(2),
	RunE:  CreateAlias,
}

var unaliasCmd = &cobra.Command{
	Use:   "go_nvm unalias [name]",
	Short: "Remove an alias for a Node.js version",
	Long:  "This command removes an existing alias for a Node.js version",
	Args:  cobra.ExactArgs(1),
	RunE:  RemoveAlias,
}

func init() {
	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(unaliasCmd)
}

func CreateAlias(cmd *cobra.Command, args []string) error {
	name := args[0]
	version := args[1]

	util.Logger.Infof("Creating alias %s for %s", name, version)
	err := nvm.SetAliasedVersion(name, version)
	if err != nil {
		return err
	}

	util.Logger.Infof("Alias %s created for version %s", name, version)
	return nil
}

func RemoveAlias(cmd *cobra.Command, args []string) error {
	name := args[0]

	util.Logger.Infof("Removing alias %s", name)
	err := nvm.RemoveAlias(name)
	if err != nil {
		return err
	}

	util.Logger.Infof("Alias %s removed", name)
	return nil
}
