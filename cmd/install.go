package cmd

import (
	"github.com/9-Realms-Dev/go_nvm/internal/nvm"
	"github.com/9-Realms-Dev/go_nvm/internal/util"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "go_nvm install [node version]",
	Short: "This will go get and install the node version",
	Long:  "installing a node version puts it in your GO_NVM_DIR/versions directory to be used for symlinks",
	Args:  cobra.MinimumNArgs(1),
	RunE:  InstallNode,
}

func init() {
	// TODO: Add flags for getting the latest versions
	rootCmd.AddCommand(installCmd)
}

func InstallNode(cmd *cobra.Command, args []string) error {
	util.Logger.Info("checking if version exists....")
	nodeVersion, err := nvm.GetVersion(args[0], true)
	if err != nil {
		return err
	}

	installPath, err := nvm.GetInstallPath(nodeVersion)
	if err != nil {
		return err
	}

	util.Logger.Infof("checking if %s is installed....", nodeVersion)
	if nvm.CheckNodeVersionInstalled(installPath) {
		util.Logger.Warnf("node %s is already installed", nodeVersion)
		return nil
	}

	url, err := nvm.GetNodeVersionURL(nodeVersion)
	if err != nil {
		return err
	}

	util.Logger.Infof("downloading and installing node %s", nodeVersion)
	err = nvm.DownloadAndSetupNode(url, installPath)
	if err != nil {
		return err
	}

	util.Logger.Infof("node %s installed successfully", nodeVersion)
	return nil
}
