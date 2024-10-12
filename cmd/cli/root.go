package cli

import (
	"github.com/spf13/cobra"
)

// RootCmd is the main Cobra command for the scraper CLI tool, providing a short
// description and setting up the command's usage for scraping Nexus Mods and returning
// the information in JSON format.
var RootCmd = &cobra.Command{
	Use:   "scraper",
	Short: "A CLI tool to scrape https://nexusmods.com mods and return the information in JSON format",
}

// Execute runs the RootCmd command, handling any errors that occur during its execution.
// Returns an error if the command fails to execute.
func Execute() error {

	if err := RootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
