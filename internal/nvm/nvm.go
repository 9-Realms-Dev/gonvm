package nvm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/9-Realms-Dev/gonvm/internal/util"
)

func SetCurrentVersion(versionPath string) error {
	version := filepath.Base(versionPath)

	isLinked, err := CreateSymlink(versionPath)
	if err != nil {
		return fmt.Errorf("error creating symlink: %w", err)
	}

	if isLinked {
		fmt.Println(styles.SuccessStyle.Render(fmt.Sprintf("Now using node %s", version)))
	} else {
		return fmt.Errorf("could not set node %s", version)
	}

	return nil
}

func CreateSymlink(versionPath string) (bool, error) {
	nvmDir, err := util.GetNvmDirectory()
	if err != nil {
		return false, fmt.Errorf("error getting NVM directory: %w", err)
	}

	currentSymlink := filepath.Join(nvmDir, "current")

	// Remove existing symlink if it exists
	err = os.Remove(currentSymlink)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("error removing existing symlink: %w", err)
	}

	err = os.Symlink(versionPath, currentSymlink)
	if err != nil {
		return false, fmt.Errorf("error creating symlink: %w", err)
	}

	return true, nil
}
