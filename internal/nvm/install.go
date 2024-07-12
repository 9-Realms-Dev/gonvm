package nvm

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	tui "github.com/9-Realms-Dev/go_nvm/internal/tui/components"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/9-Realms-Dev/go_nvm/internal/util"
)

func DownloadAndSetupNode(url, installPath string) error {
	// Create the install directory if it doesn't exist
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return err
	}

	// Download the Node.js binary
	util.Logger.Debugf("downloading node.js from %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileName := "node.tar.gz"
	if runtime.GOOS == "windows" {
		fileName = "node.exe"
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tui.CopyWithProgress(file, resp.Body, resp.ContentLength, "Downloading..."); err != nil {
		return err
	}

	// Extract the file if needed
	if runtime.GOOS != "windows" {
		if err := extractTarWithProgress(file, installPath); err != nil {
			return err
		}

		// Find the extracted directory
		entries, err := os.ReadDir(installPath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			util.Logger.Infof("entry: %s", entry.Name())
			if entry.IsDir() {
				oldDir := filepath.Join(installPath, entry.Name())
				if err := shiftContent(installPath, oldDir); err != nil {
					return err
				}
				break
			}
		}
	} else {
		// For Windows, just move the exe file
		newPath := filepath.Join(installPath, "node.exe")
		if err := os.Rename(fileName, newPath); err != nil {
			return err
		}
	}

	// Remove the downloaded file
	return os.Remove(fileName)
}

func GetNodeVersionURL(version string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "386":
			return fmt.Sprintf("https://nodejs.org/dist/%s/node-%s-win-x86/node.exe", version, version), nil
		case "amd64":
			return fmt.Sprintf("https://nodejs.org/dist/%s/node-%s-win-x64/node.exe", version, version), nil
		default:
			return "", fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
		}
	case "darwin", "linux":
		return fmt.Sprintf("https://nodejs.org/dist/%s/node-%s-linux-x64.tar.gz", version, version), nil
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func GetInstallPath(version string) (string, error) {
	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		return setDirectoryPath(version)
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func setDirectoryPath(version string) (string, error) {
	nvmPath, exists := os.LookupEnv("GO_NVM_DIR")
	if !exists {
		var err error
		nvmPath, err = util.SetDefaultDirectory()
		if err != nil {
			return "", err
		}
	}

	versionPath := filepath.Join(nvmPath, "versions", version)

	// Ensure the directory exists
	err := os.MkdirAll(versionPath, 0755)
	if err != nil {
		return "", err
	}

	return versionPath, nil
}

func extractTarWithProgress(tarFile *os.File, installPath string) error {
	util.Logger.Infof("extracting tar file to %s", installPath)
	// Seek back to the beginning of the file
	_, err := tarFile.Seek(0, 0)
	if err != nil {
		util.Logger.Errorf("Error seeking to beginning of file: %v", err)
		return err
	}

	fileContents, err := io.ReadAll(tarFile)
	if err != nil {
		util.Logger.Errorf("Error reading entire file: %v", err)
		return err
	}
	util.Logger.Infof("Read %d bytes from file", len(fileContents))

	gzr, err := gzip.NewReader(bytes.NewReader(fileContents))
	if err != nil {
		util.Logger.Errorf("Error creating gzip reader from memory: %v", err)
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	util.Logger.Infof("starting to extract files")
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(installPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}

	return nil
}

func shiftContent(installPath, oldDir string) error {
	entries, err := os.ReadDir(oldDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		oldPath := filepath.Join(oldDir, entry.Name())
		newPath := filepath.Join(installPath, entry.Name())
		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
	}

	return os.Remove(oldDir)
}
