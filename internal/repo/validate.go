package repo

import (
	"fmt"
	"os"
	"strings"
)

var allowedURLPrefixes = []string{
	"https://",
	"http://",
	"git://",
	"ssh://",
	"git@",
}

// ValidateRemoteURL checks that a remote URL is non-empty and uses an allowed protocol.
func ValidateRemoteURL(url string) error {
	if url == "" {
		return fmt.Errorf("remote URL must not be empty")
	}
	if strings.HasPrefix(url, "-") {
		return fmt.Errorf("remote URL must not start with a dash")
	}
	for _, prefix := range allowedURLPrefixes {
		if strings.HasPrefix(url, prefix) {
			return nil
		}
	}
	return fmt.Errorf("remote URL must start with one of: https://, http://, git://, ssh://, git@")
}

// ValidateLocalPath checks that a local path is a valid, existing directory.
func ValidateLocalPath(path string, strict bool) error {
	if path == "" {
		return fmt.Errorf("path must not be empty")
	}
	if strict {
		if strings.Contains(path, "..") {
			return fmt.Errorf("path must not contain '..'")
		}
		if path == "/" {
			return fmt.Errorf("root path '/' is not allowed")
		}
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("cannot stat path %s: %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}
	return nil
}
