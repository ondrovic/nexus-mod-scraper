package utils

import (
	"os"
	"sync"
)

func ConcurrentFetch(tasks ...func() error) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		go func(t func() error) {
			defer wg.Done()
			if err := t(); err != nil {
				errChan <- err
			}
		}(task)
	}

	wg.Wait()
	close(errChan)

	// Return the first error encountered, if any
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func EnsureDirExists(path string) error {
	// Check if the directory exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		err := os.MkdirAll(path, os.ModePerm) // os.ModePerm ensures the directory is created with the correct permissions
		if err != nil {
			return err
		}

	} else if err != nil {
		return err
	}
	return nil
}
