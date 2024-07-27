package cmd

import (
	"github.com/9-Realms-Dev/gonvm/internal/nvm"
	"github.com/9-Realms-Dev/gonvm/internal/util"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall a specific Node.js version",
	Long:  "This will remove the version directory from your GO_NVM_DIR",
	RunE:  UninstallNode,
}

func UninstallNode(cmd *cobra.Command, args []string) error {
	reqVersion := args[0]
	// TODO: Replace with get local versions
	nodeVersion, err := nvm.GetVersion(reqVersion, false, false)
	if err != nil {
		return err
	}

	installPath, err := nvm.GetInstallPath(nodeVersion)
	if err != nil {
		return err
	}

	if nvm.IsNodeVersionInstalled(installPath) {
		util.Logger.Infof("Uninstalling node %s...", reqVersion)
		err := nvm.RemoveVersionByPath(installPath)
		if err != nil {
			return err
		}
		util.Logger.Infof("Node %s has been uninstalled", nodeVersion)
	} else {
		util.Logger.Infof("Node %s is not installed", nodeVersion)
	}

	return nil
}
