package spinners

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/theckman/yacspin"
)

// CreateSpinner initializes and returns a yacspin spinner with the provided
// start and stop messages, characters, and failure configurations.
func CreateSpinner(startMessage, stopCharacter, stopMessage, stopFailCharacter, stopFailMessage string) *yacspin.Spinner {
	cfg := yacspin.Config{
		Frequency:         100 * time.Millisecond,
		Colors:            []string{"fgHiBlue"},
		CharSet:           yacspin.CharSets[14],
		Suffix:            " ",
		SuffixAutoColon:   true,
		Message:           startMessage,
		StopCharacter:     stopCharacter,
		StopColors:        []string{"fgHiGreen"},
		StopMessage:       stopMessage,
		StopFailCharacter: stopFailCharacter,
		StopFailColors:    []string{"fgHiRed"},
		StopFailMessage:   stopFailMessage,
	}

	s, err := yacspin.New(cfg)
	if err != nil {
		fmt.Printf("failed to create spinner: %v\n", err)
		os.Exit(1)
	}

	return s
}

// stopOnSignal stops the spinner with a failure message if an interrupt or
// termination signal is received, ensuring proper cleanup before exiting.
func stopOnSignal(spinner *yacspin.Spinner) {
	// ensure we stop the spinner before exiting, otherwise cursor will remain
	// hidden and terminal will require a `reset`
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh

		spinner.StopFailMessage("interrupted")

		// ignoring error intentionally
		_ = spinner.StopFail()

		os.Exit(0)
	}()
}
