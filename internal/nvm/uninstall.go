package nvm

import (
	"fmt"
	"github.com/9-Realms-Dev/gonvm/internal/tui/styles"
	"os"
)

func RemoveVersionByPath(versionPath string) error {
	err := os.RemoveAll(versionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(styles.ErrorStyle.Render("Error: Could not find %s"), versionPath)
		} else if os.IsPermission(err) {
			return fmt.Errorf(styles.ErrorStyle.Render("Error: Permission denied for %s"), versionPath)
		} else {
			return fmt.Errorf(styles.ErrorStyle.Render("Error: %v"), err)
		}
	}

	fmt.Println(styles.SuccessStyle.Render(fmt.Sprintf("Successfully removed %s", versionPath)))

	return nil
}

// TODO: Remove all LTS but the latest LTS version, find largest even number remove all others.

// TODO: Uninstall all versions
