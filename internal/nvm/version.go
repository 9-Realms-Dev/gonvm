package nvm

import (
	"fmt"
	"github.com/9-Realms-Dev/gonvm/internal/tui/styles"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	tui "github.com/9-Realms-Dev/gonvm/internal/tui/components"
	"github.com/9-Realms-Dev/gonvm/internal/util"
	"github.com/PuerkitoBio/goquery"
)

func GetVersion(version string, checkLatest, acceptAll bool) (string, error) {
	if checkLatest {
		if version == "latest" {
			fmt.Println(styles.InfoStyle.Render("Checking for latest version..."))
			return GetRemoteLatestVersion(acceptAll)
		}

		if version == "lts" {
			fmt.Println(styles.InfoStyle.Render("Checking for latest LTS version..."))
			return GetRemoteLTSVersion(acceptAll)
		}

		if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		fmt.Println(styles.InfoStyle.Render(fmt.Sprintf("Checking for version %s...", version)))
		return GetRemoteVersion(version)
	} else {
		if CheckValidVersionPattern(version) {
			return GetLocalVersion(version, checkLatest, acceptAll)
		} else {
			return GetAliasVersion(version, acceptAll)
		}
	}
}

func GetAliasVersion(alias string, acceptAll bool) (string, error) {
	// First, try to get the aliased version
	aliasVersion, err := GetAliasedVersion(alias)
	if err != nil {
		return "", fmt.Errorf("error getting aliased version: %w", err)
	}

	if aliasVersion != "" {
		fmt.Println(styles.PromptStyle.Render(fmt.Sprintf("Found alias version: %s", aliasVersion)))
		return aliasVersion, nil
	}

	// If alias is 'latest' or 'lts' and not found, prompt user to fetch remote version
	if alias == "latest" || alias == "lts" {
		if acceptAll {
			// TODO: Refactor into a function
			fmt.Println(styles.PromptStyle.Render("Fetching remote version..."))
			// Assuming we have a getVersion function that can fetch the latest version
			version, err := GetVersion(alias, true, acceptAll)
			if err != nil {
				return "", fmt.Errorf("error fetching remote version: %w", err)
			}
			fmt.Println(styles.PromptStyle.Render(fmt.Sprintf("Found remote version: %s", version)))
			return version, nil
		}

		question := fmt.Sprintf("%s not found. Would you like to get remote %s version?", alias, alias)
		confirm, err := tui.ConfirmPrompt(question)
		if err != nil {
			return "", fmt.Errorf("error during user prompt: %w", err)
		}

		if confirm {
			// TODO: Refactor into a function
			fmt.Println(styles.PromptStyle.Render("Fetching remote version..."))
			// Assuming we have a getVersion function that can fetch the latest version
			version, err := GetVersion(alias, true, acceptAll)
			if err != nil {
				return "", fmt.Errorf("error fetching remote version: %w", err)
			}
			fmt.Println(styles.PromptStyle.Render(fmt.Sprintf("Found remote version: %s", version)))
			return version, nil
		}
	}

	fmt.Println(styles.ErrorStyle.Render(fmt.Sprintf("No version found for alias: %s", alias)))
	return "", nil
}

func GetLocalVersion(version string, checklatest, acceptAll bool) (string, error) {
	versions, err := LocalVersions()
	if err != nil {
		return "", fmt.Errorf("error getting local versions: %w", err)
	}

	if len(versions) > 0 {
		for _, v := range versions {
			if v == version {
				fmt.Println(styles.PromptStyle.Render(fmt.Sprintf("Found local version: %s", version)))
				return version, nil
			}
		}

		// Version not found, prompt to install
		question := fmt.Sprintf("Version %s not found. Would you like to install it?", version)
		confirm, err := tui.ConfirmPrompt(question)
		if err != nil {
			return "", fmt.Errorf("error during user prompt: %w", err)
		}

		if confirm {
			fmt.Println(styles.PromptStyle.Render("Fetching and installing version..."))
			// Assuming we have a getVersion function that can fetch and install the version
			installedVersion, err := GetVersion(version, checklatest, acceptAll)
			if err != nil {
				return "", fmt.Errorf("error fetching and installing version: %w", err)
			}
			fmt.Println(styles.PromptStyle.Render(fmt.Sprintf("Installed version: %s", installedVersion)))
			return installedVersion, nil
		}
	} else {
		fmt.Println(styles.ErrorStyle.Render("You have no local versions..."))
	}

	return "", nil
}

func GetRemoteVersion(version string) (string, error) {
	url := fmt.Sprintf("https://nodejs.org/dist/%s/", version)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error parsing HTML: %w", err)
		}

		var lastHref string
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			lastHref = href
		})

		re := regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)
		matches := re.FindStringSubmatch(lastHref)
		if len(matches) > 0 {
			return matches[0], nil
		}
		return "", fmt.Errorf("no version number found in the last link")
	} else {
		fmt.Println(styles.WarnStyle.Render(fmt.Sprintf("Could not find version %s. Checking for latest version...", version)))
		latestVersions, err := GetRemoteVersions(version)
		if err != nil {
			return "", fmt.Errorf("error getting remote versions: %w", err)
		}
		if len(latestVersions) > 0 {
			return latestVersions[len(latestVersions)-1], nil
		} else {
			fmt.Println(styles.WarnStyle.Render(fmt.Sprintf("Could not find version %s", version[:3])))
			return "", nil
		}
	}
}

func GetRemoteVersions(version string) ([]string, error) {
	if len(version) > 3 {
		version = version[:3]
	}

	fmt.Println(styles.InfoStyle.Render(fmt.Sprintf("Checking for versions starting with %s...", version)))

	url := "https://nodejs.org/dist/"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error parsing HTML: %w", err)
		}

		re := regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)
		versionNumbers := make([]string, 0)

		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			matches := re.FindStringSubmatch(href)
			if len(matches) > 0 {
				versionNumber := matches[0]
				if strings.HasPrefix(versionNumber, version) {
					versionNumbers = append(versionNumbers, versionNumber)
				}
			}
		})

		return versionNumbers, nil
	} else {
		return nil, fmt.Errorf("error: request was not successful. Status code: %d", resp.StatusCode)
	}
}

func GetRemoteLatestVersion(acceptAll bool) (string, error) {
	url := "https://nodejs.org/dist/latest/"
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error parsing HTML: %w", err)
		}

		var latestVersion string
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			ok, version := matchNodeArchive(href)
			if ok {
				latestVersion = version
				return
			}
		})

		if latestVersion == "" {
			return "", fmt.Errorf("could not find latest version")
		}

		question := fmt.Sprintf("Would you like to set %s as the default latest version?", latestVersion)
		err = confirmLatest(question, latestVersion, "latest", acceptAll)

		return latestVersion, nil
	} else {
		return "", fmt.Errorf("request was not successful. Status code: %d", resp.StatusCode)
	}
}

func GetRemoteLTSVersion(acceptAll bool) (string, error) {
	util.Logger.Info("Checking for latest LTS version...")
	url := "https://nodejs.org/en"
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error parsing HTML: %w", err)
		}

		var ltsVersionLink string
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			title := s.Text()
			if regexp.MustCompile(`LTS`).MatchString(title) {
				ltsVersionLink, _ = s.Attr("href")
				return
			}
		})

		ltsVersion := extractVersionNumber(ltsVersionLink)

		if ltsVersion == "" {
			return "", fmt.Errorf("could not find latest LTS version")
		}

		// Check alias.toml to see what lts version is set
		aliasVersion, err := GetAliasedVersion("lts")
		if err != nil {
			return "", fmt.Errorf("error getting aliased version: %w", err)
		}

		if aliasVersion != ltsVersion {
			question := fmt.Sprintf("Would you like to set %s as the default LTS version?", ltsVersion)
			err = confirmLatest(question, ltsVersion, "lts", acceptAll)
			if err != nil {
				return "", err
			}
		} else {
			fmt.Println(styles.InfoStyle.Render(fmt.Sprintf("Latest LTS version is already set to %s", ltsVersion)))
		}

		return ltsVersion, nil
	} else {
		return "", fmt.Errorf("request was not successful. Status code: %d", resp.StatusCode)
	}
}

func RemoteVersions() ([]string, error) {
	url := "https://nodejs.org/dist/"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error parsing HTML: %w", err)
		}

		versions := []string{}

		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			match := extractVersionNumber(href)
			versions = append(versions, match)
		})

		return versions, nil
	} else {
		return nil, fmt.Errorf("request was not successful. Status code: %d", resp.StatusCode)
	}
}

func LocalVersions() ([]string, error) {
	nvmDir, err := util.GetNvmDirectory()
	if err != nil {
		// If GetNvmDirectory returns an error, it means the directory is not set
		// Let's prompt the user to fix the setup
		question := "GO_NVM_DIR not set. Would you like to fix your setup?"
		fixSetup, err := tui.ConfirmPrompt(question)
		if err != nil {
			return nil, fmt.Errorf("error during user prompt: %w", err)
		}

		if fixSetup {
			nvmDir, err = util.SetDefaultDirectory()
			if err != nil {
				return nil, fmt.Errorf("error setting default directory: %w", err)
			}
			fmt.Println(styles.InfoStyle.Render(fmt.Sprintf("Set GO_NVM_DIR to %s", nvmDir)))
			fmt.Println(styles.InfoStyle.Render("Run `go_nvm install <version>` to install a version of node"))
			return []string{}, nil
		} else {
			fmt.Println(styles.WarnStyle.Render("GO_NVM_DIR not set."))
			return nil, nil
		}
	}

	versionsPath := filepath.Join(nvmDir, "versions")
	files, err := os.ReadDir(versionsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the versions directory doesn't exist, return an empty slice
			return []string{}, nil
		}
		return nil, fmt.Errorf("error reading versions directory: %w", err)
	}

	var versions []string
	for _, file := range files {
		if file.IsDir() {
			versions = append(versions, file.Name())
		}
	}

	return versions, nil
}

func IsNodeVersionInstalled(versionPath string) bool {
	// Check if the version path exists and has a bin directory inside it
	_, err := os.Stat(filepath.Join(versionPath, "bin"))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(styles.ErrorStyle.Render(fmt.Sprintf("Error: Could not find %s", versionPath)))
		} else if os.IsPermission(err) {
			fmt.Println(styles.ErrorStyle.Render(fmt.Sprintf("Error: Permission denied for %s", versionPath)))
		} else {
			fmt.Println(styles.ErrorStyle.Render(fmt.Sprintf("Error checking %s: %v", versionPath, err)))
		}
		return false
	}
	return true
}

func CheckValidVersionPattern(version string) bool {
	versionPattern := `^(v?\d+(\.\d+){0,2})$`
	match, _ := regexp.MatchString(versionPattern, version)
	return match
}

func extractVersionNumber(href string) string {
	re := regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)
	matches := re.FindStringSubmatch(href)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

func matchNodeArchive(filename string) (bool, string) {
	pattern := `^node-(v\d+\.\d+\.\d+)\.tar\.gz$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(filename)
	if matches == nil {
		return false, ""
	}

	return true, matches[1] // matches[1] contains the version number
}

func confirmLatest(question, version, alias string, acceptAll bool) error {
	var setLatest bool
	var err error
	if !acceptAll {
		setLatest, err = tui.ConfirmPrompt(question)
		if err != nil {
			return fmt.Errorf("error during user prompt: %w", err)
		}
	}

	if setLatest || acceptAll {
		err := SetAliasedVersion(alias, version)
		if err != nil {
			return fmt.Errorf("error setting aliased version: %w", err)
		}
		fmt.Println(styles.InfoStyle.Render(fmt.Sprintf("Set %s as the latest %s version", version, alias)))
	}

	return nil
}
