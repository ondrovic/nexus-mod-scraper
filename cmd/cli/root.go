package cli

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the root command of the CLI tool.
// It defines the command name, description, and usage for the tool that scrapes nexusmods.com mods and returns data in JSON format.
var RootCmd = &cobra.Command{
	Use:   "scraper",
	Short: "A CLI tool to scrape https://nexusmods.com mods and return the information in JSON format",
}

// Execute runs the root command and handles any errors that occur during execution.
// Returns: An error if the command execution fails.
func Execute() error {

	if err := RootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
