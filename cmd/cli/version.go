package cli

import (
	"go.szostok.io/version/extension"
)

const (
	// RepoOwner is the owner of the GitHub repository.
	RepoOwner string = "ondrovic"
	// RepoName is the name of the GitHub repository.
	RepoName string = "nexus-mods-scraper"
)

// init initializes the command-line interface by adding the version command
// to the root command, including an upgrade notice with the repository owner
// and name.
func init() {
	RootCmd.AddCommand(
		extension.NewVersionCobraCmd(
			extension.WithUpgradeNotice(RepoOwner, RepoName),
		),
	)
}
