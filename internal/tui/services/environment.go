package services

import (
	"fmt"
	"os/exec"

	"github.com/9-Realms-Dev/gonvm/internal/nvm"
)

func runCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

type NodeVersion struct {
	NodeVersion string
	NpmVersion  string
}

func GetCurrentDetails() (*NodeVersion, error) {
	npmVersion, err := runCommand("npm", "--version")
	if err != nil {
		return nil, fmt.Errorf("Error checking npm version: %v", err)
	}

	nodeVersion, err := runCommand("node", "--version")
	if err != nil {
		return nil, fmt.Errorf("Error checking node version: %v", err)
	}

	return &NodeVersion{NodeVersion: nodeVersion, NpmVersion: npmVersion}, nil
}

func GetGlobalPackages() (*string, error) {
	output, err := runCommand("npm", "list", "-g", "--depth=0")
	if err != nil {
		return nil, fmt.Errorf("Error checking global packages: %v", err)
	}

	return &output, nil
}

func GetVersions() ([]string, error) {
	versionList, err := nvm.LocalVersions()
	if err != nil {
		return nil, fmt.Errorf("Error checking local versions: %v", err)
	}

	return versionList, nil
}
