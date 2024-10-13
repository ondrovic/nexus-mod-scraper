//go:build windows
// +build windows

package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDataPath_Windows(t *testing.T) {
	userProfileDir := os.Getenv("USERPROFILE")
	expectedPath := filepath.Join(userProfileDir, ".nexus-mods-scraper", "data")
	actualPath := GetDataStoragePath()

	assert.Equal(t, expectedPath, actualPath, "The Windows data path is incorrect.")
}
