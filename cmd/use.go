package cmd

import (
	"fmt"

	"github.com/9-Realms-Dev/go_nvm/internal/nvm"
	tui "github.com/9-Realms-Dev/go_nvm/internal/tui/components"
	"github.com/9-Realms-Dev/go_nvm/internal/util"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "This will set active node version",
	Long:  "Use will take over the symlink to the active node version",
	Args:  cobra.MinimumNArgs(1),
	RunE:  UseNvm,
}

func init() {
	useCmd.Flags().BoolVarP(&latestFlag, "latest", "l", false, "get the latest version")
	useCmd.Flags().BoolVarP(&acceptAllFlag, "yes", "y", false, "accept all prompts")
	rootCmd.AddCommand(useCmd)
}

func UseNvm(cmd *cobra.Command, args []string) error {
	checkLatest, err := cmd.Flags().GetBool("latest")
	acceptAll, err := cmd.Flags().GetBool("yes")
	if err != nil {
		return err
	}

	nodeVersion, err := nvm.GetVersion(args[0], checkLatest, acceptAll)
	if err != nil {
		return err
	}
	versionPath, err := nvm.GetInstallPath(nodeVersion)
	if err != nil {
		return err
	}

	util.Logger.Infof("checking if %s is set....", nodeVersion)
	if !nvm.IsNodeVersionInstalled(versionPath) {
		confirmation, err := tui.ConfirmPrompt(fmt.Sprintf("node %s is not installed. Would you like to install it?", nodeVersion))
		if err != nil {
			return err
		}

		if confirmation {
			InstallNode(cmd, args)
		} else {
			util.Logger.Warnf("node %s was not installed", nodeVersion)
			return nil
		}
	}

	return nvm.SetCurrentVersion(versionPath)
}
