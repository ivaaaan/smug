package main

import (
	"os"
	"strings"
)

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return path
		}

		return strings.Replace(path, "~", userHome, 1)
	}

	return path
}

func Contains(slice []string, s string) bool {
	for _, e := range slice {
		if e == s {
			return true
		}
	}

	return false
}
