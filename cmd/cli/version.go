package cli

import (
	"go.szostok.io/version/extension"
)

const (
	// RepoOwner is the owner of the GitHub repository.
	RepoOwner string = "ondrovic"
	// RepoName is the name of the GitHub repository.
	RepoName string = "nexus-mods-scraper" // TODO: update this value
)

// init adds a new version command to RootCmd. The command includes an upgrade notice
// based on the specified GitHub repository owner and name.
func init() {
	RootCmd.AddCommand(
		extension.NewVersionCobraCmd(
			extension.WithUpgradeNotice(RepoOwner, RepoName),
		),
	)
}
