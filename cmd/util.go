package cmd

import (
	"os"
	"strings"
)

func getFolderName(path string) string {
	parts := strings.Split(path, string(os.PathSeparator))

	if len(parts) == 0 || parts[0] == "." {
		return ""
	}

	return parts[len(parts)-1]
}

func checkFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err // Other error
	}
	return true, nil
}
