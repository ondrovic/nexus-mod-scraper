package main

import (
	"runtime"

	"github.com/ondrovic/nexus-mods-scraper/cmd/cli"

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
