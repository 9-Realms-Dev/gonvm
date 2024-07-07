package node

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func removeVersionByPath(versionPath string) error {
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("red"))

	err := os.RemoveAll(versionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(errorStyle.Render("Error: Could not find %s"), versionPath)
		} else if os.IsPermission(err) {
			return fmt.Errorf(errorStyle.Render("Error: Permission denied for %s"), versionPath)
		} else {
			return fmt.Errorf(errorStyle.Render("Error: %v"), err)
		}
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("green"))
	fmt.Println(successStyle.Render(fmt.Sprintf("Successfully removed %s", versionPath)))

	return nil
}

// TODO: Remove all LTS but the latest LTS version, find largest even number remove all others.

// TODO: Uninstall all versions
