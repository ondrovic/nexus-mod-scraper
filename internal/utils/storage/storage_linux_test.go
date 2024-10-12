//go:build linux
// +build linux

package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDataPath_Linux(t *testing.T) {
	homeDir := os.Getenv("HOME")
	expectedPath := filepath.Join(homeDir, ".nexus-mods-scraper", "data")
	actualPath := GetDataStoragePath()

	assert.Equal(t, expectedPath, actualPath, "The Linux data path is incorrect.")
}
