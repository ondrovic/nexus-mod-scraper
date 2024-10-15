//go:build windows
// +build windows

package storage

import (
	"os"
	"path/filepath"
)

// GetDataStoragePath returns the data storage path in the user's home directory,
// specifically for the nexus-mods-scraper application.
func GetDataStoragePath() string {
	userProfileDir := os.Getenv("USERPROFILE")
	return filepath.Join(userProfileDir, ".nexus-mods-scraper", "data")
}
