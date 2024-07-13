package testutils

import (
	"fmt"
	"github.com/charmbracelet/log"
	"os"
)

func CheckNodeVersionDownloaded(version string) bool {
	// check if the GO_NVM_DIR exists
	dir := os.Getenv("GO_NVM_DIR")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Print("GO_NVM_DIR does not exist")
		return false
	}
	// check if the version exists
	if _, err := os.Stat(fmt.Sprintf("%s/versions/%s/bin", dir, version)); os.IsNotExist(err) {
		log.Print("Version does not exist")
		return false
	}

	return true
}

func CheckNodeVersionInstalled(version string) bool {
	// Check if the version is downloaded
	ok := CheckNodeVersionDownloaded(version)
	if !ok {
		log.Print("Version not downloaded")
		return false
	}
	// Check if the version is set to current
	dir := os.Getenv("GO_NVM_DIR")
	if _, err := os.Readlink(fmt.Sprintf("%s/current", dir)); os.IsNotExist(err) {
		log.Print("Current link does not exist")
		return false
	}

	return true
}
