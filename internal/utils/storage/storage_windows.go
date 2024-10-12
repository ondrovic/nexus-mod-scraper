//go:build windows
// +build windows

package storage

import (
	"os"
	"path/filepath"
)

// GetDataStoragePath returns the data storage path in the user's LOCALAPPDATA
// directory, specifically for the nexus-mods-scraper application.
func GetDataStoragePath() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	return filepath.Join(localAppData, "nexus-mods-scraper", "data")
}
