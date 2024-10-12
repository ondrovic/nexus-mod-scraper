//go:build windows
// +build windows

package storage

import (
	"os"
	"path/filepath"
)

func GetDataStoragePath() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	return filepath.Join(localAppData, "nexus-mod-scraper", "data")
}
