package testutils

import (
	"fmt"
	"regexp"
)

func ExtractVersion(output string) (string, error) {
	// Regular expression to match vX.X.X pattern
	re := regexp.MustCompile(`v\d+\.\d+\.\d+`)

	// Find all matches in the output
	matches := re.FindAllString(output, -1)

	if len(matches) > 0 {
		// Return the last match (in case there are multiple version numbers)
		return matches[len(matches)-1], nil
	}

	return "", fmt.Errorf("no version number found in the output")
}
