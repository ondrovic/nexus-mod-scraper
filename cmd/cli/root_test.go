package cli

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRootCmd_Initialized(t *testing.T) {
	// Ensure that RootCmd is correctly initialized
	assert.Equal(t, "nexus-mods-scraper", RootCmd.Use)
	assert.Equal(t, "A CLI tool to scrape https://nexusmods.com mods and return the information in JSON format", RootCmd.Short)
}

func TestExecute_Success(t *testing.T) {
	// Mock a successful command execution
	mockCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing (successful execution)
		},
	}

	// Replace RootCmd with the mock command for the test
	RootCmd = mockCmd

	// Execute the command and ensure no error is returned
	err := Execute()
	assert.NoError(t, err)
}

func TestExecute_Failure(t *testing.T) {
	// Mock a command that fails
	mockCmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("execution failed")
		},
	}

	// Replace RootCmd with the mock command for the test
	RootCmd = mockCmd

	// Execute the command and ensure the error is returned
	err := Execute()
	assert.Error(t, err)
	assert.Equal(t, "execution failed", err.Error())
}
