package utils

import (
	"errors"
	"os"
	"testing"
)

func TestConcurrentFetch_FirstTaskFails(t *testing.T) {
	// Arrange
	expectedErr := errors.New("task1 failed")
	task1 := func() error { return expectedErr }
	task2 := func() error { return nil }

	// Act
	err := ConcurrentFetch(task1, task2)

	// Assert
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestConcurrentFetch_SecondTaskFails(t *testing.T) {
	// Arrange
	expectedErr := errors.New("task2 failed")
	task1 := func() error { return nil }
	task2 := func() error { return expectedErr }

	// Act
	err := ConcurrentFetch(task1, task2)

	// Assert
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestConcurrentFetch_MultipleTasksFail(t *testing.T) {
	// Arrange
	task1Err := errors.New("task1 failed")
	task2Err := errors.New("task2 failed")
	task1 := func() error { return task1Err }
	task2 := func() error { return task2Err }

	// Act
	err := ConcurrentFetch(task1, task2)

	// Assert
	if err != task1Err && err != task2Err {
		t.Errorf("Expected error to be either %v or %v, got %v", task1Err, task2Err, err)
	}
}

func TestEnsureDirExists_DirAlreadyExists(t *testing.T) {
	// Arrange
	existingDir := "existingDir"
	os.Mkdir(existingDir, os.ModePerm) // Create the directory for the test
	defer os.Remove(existingDir)       // Clean up after test

	// Act
	err := EnsureDirExists(existingDir)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestEnsureDirExists_DirDoesNotExist(t *testing.T) {
	// Arrange
	newDir := "newDir"
	defer os.Remove(newDir) // Clean up after test

	// Act
	err := EnsureDirExists(newDir)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the directory was created
	_, statErr := os.Stat(newDir)
	if os.IsNotExist(statErr) {
		t.Errorf("Expected directory to be created, but it was not")
	}
}

func TestEnsureDirExists_CannotCreateDir(t *testing.T) {
	// Arrange
	invalidDir := "" // Empty directory name should cause an error

	// Act
	err := EnsureDirExists(invalidDir)

	// Assert
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
}
