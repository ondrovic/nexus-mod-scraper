//go:build darwin
// +build darwin

package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDataPath_Darwin(t *testing.T) {
	homeDir := os.Getenv("HOME")
	expectedPath := filepath.Join(homeDir, ".nexus-mod-scraper", "data")
	actualPath := GetDataStoragePath()

	assert.Equal(t, expectedPath, actualPath, "The macOS data path is incorrect.")
}
