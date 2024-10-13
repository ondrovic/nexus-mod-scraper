package main

import (
	"fmt"
	"runtime"

	sCli "github.com/ondrovic/common/utils/cli"
	"github.com/ondrovic/nexus-mods-scraper/cmd/cli"
)

type clearScreenFunc func(interface{}) error

func run(clearScreen clearScreenFunc, executeFunc func() error) error {
	if err := clearScreen(runtime.GOOS); err != nil {
		return fmt.Errorf("error clearing terminal: %w", err)
	}

	if err := executeFunc(); err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}

	return nil
}

func executeMain(clearScreen clearScreenFunc, executeFunc func() error) {
	if err := run(clearScreen, executeFunc); err != nil {
		return
	}
}

func main() {
	executeMain(sCli.ClearTerminalScreen, cli.RootCmd.Execute)
}
