//go:build darwin
// +build darwin

package storage

import (
	"os"
	"path/filepath"
)

func GetDataStoragePath() string {
	homeDir := os.Getenv("HOME")
	return filepath.Join(homeDir, ".nexus-mod-scraper", "data")
}
