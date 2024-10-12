package main

import (
	"nexus-mods-scraper/cmd/cli"
	"runtime"

	sCli "github.com/ondrovic/common/utils/cli"
)

func main() {
	if err := sCli.ClearTerminalScreen(runtime.GOOS); err != nil {
		return
	}

	if err := cli.RootCmd.Execute(); err != nil {
		return
	}
}
