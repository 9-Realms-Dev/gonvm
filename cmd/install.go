package cmd

import (
	"github.com/9-Realms-Dev/gonvm/internal/nvm"
	"github.com/9-Realms-Dev/gonvm/internal/util"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "This will go get and install the node version",
	Long:  "installing a node version puts it in your GO_NVM_DIR/versions directory to be used for symlinks",
	Args:  cobra.MinimumNArgs(1),
	RunE:  InstallNode,
}

func init() {
	installCmd.PersistentFlags().BoolVarP(&latestFlag, "latest", "l", false, "get the latest version")
	installCmd.PersistentFlags().BoolVarP(&acceptAllFlag, "yes", "y", false, "accept all prompts")
	rootCmd.AddCommand(installCmd)
}

func InstallNode(cmd *cobra.Command, args []string) error {
	util.Logger.Info("checking if version exists....")
	checkLatest, err := cmd.Flags().GetBool("latest")
	acceptAll, err := cmd.Flags().GetBool("yes")
	if err != nil {
		return err
	}

	nodeVersion, err := nvm.GetVersion(args[0], checkLatest, acceptAll)
	if err != nil {
		return err
	}

	installPath, err := nvm.GetInstallPath(nodeVersion)
	if err != nil {
		return err
	}

	util.Logger.Infof("checking if %s is installed....", nodeVersion)
	if nvm.IsNodeVersionInstalled(installPath) {
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
