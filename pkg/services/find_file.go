package services

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindFile(file string) (string, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Navigate upwards in the directory hierarchy until we find the .env file
	for {
		envPath := filepath.Join(wd, file)
		_, err = os.Stat(envPath)
		if err == nil {
			// .env file found
			return envPath, nil
		}

		// Move to the parent directory
		wd = filepath.Dir(wd)

		// If we've reached the root directory, stop searching
		if wd == filepath.Dir(wd) {
			break
		}
	}

	// .env file not found
	return "", fmt.Errorf(".env file not found")
}
