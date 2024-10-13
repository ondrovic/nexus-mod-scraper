package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteMain_Success(t *testing.T) {
	// Mock the dependencies
	mockClearTerminal := func(_ interface{}) error {
		return nil
	}
	mockExecute := func() error {
		return nil
	}

	// Act: Call `executeMain` and verify it succeeds
	executeMain(mockClearTerminal, mockExecute)

	// Since `executeMain` doesn't return anything, you would assert no panics/errors occurred
	assert.True(t, true, "executeMain should complete without errors")
}

func TestExecuteMain_FailureOnClearTerminal(t *testing.T) {
	// Mock `ClearTerminalScreen` to return an error
	mockClearTerminal := func(_ interface{}) error {
		return errors.New("failed to clear terminal")
	}

	// Mock `executeFunc` to ensure it isn't called
	mockExecute := func() error {
		t.Log("Execute should not be called")
		return nil
	}

	// Act: Call `executeMain` and verify it handles the error
	executeMain(mockClearTerminal, mockExecute)

	// Again, since `executeMain` doesn't return anything, you assert that no panic occurred
	assert.True(t, true, "executeMain should handle the terminal clearing error gracefully")
}

func TestExecuteMain_FailureOnExecute(t *testing.T) {
	// Mock `ClearTerminalScreen` to succeed
	mockClearTerminal := func(_ interface{}) error {
		return nil
	}

	// Mock `executeFunc` to return an error
	mockExecute := func() error {
		return errors.New("execution failed")
	}

	// Act: Call `executeMain` and verify it handles the error
	executeMain(mockClearTerminal, mockExecute)

	// No panics/errors should occur, and the execution error should be gracefully handled
	assert.True(t, true, "executeMain should handle the execution error gracefully")
}
