package exporters

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"nexus-mods-scraper/internal/types"
	"nexus-mods-scraper/internal/utils"
	"nexus-mods-scraper/internal/utils/formatters"
)

func DisplayResults(sc types.CliFlags, results types.Results) error {
	jsonResults, err := formatters.FormatResultsAsJson(results.Mods)
	if err != nil {
		return fmt.Errorf("error while attempting to format results: %v", err)
	}

	formatters.PrintPrettyJson(jsonResults)
	return nil
}

func SaveCookiesToJson(dir string, filename string, data interface{}) error {
	// Check if the directory exists, if not create it
	if err := utils.EnsureDirExists(dir); err != nil {
		return err
	}
	// Join the directory and filename using filepath.Join for cross-platform compatibility
	fullPath := filepath.Join(dir, filename)

	// Open the file for writing (create if not exists, truncate if it exists)
	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert the data to a JSON formatted byte slice
	jsonData, err := json.MarshalIndent(data, "", "    ") // Using 4 spaces for indentation
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	fmt.Printf("Extracted cookies saved to %s\n", fullPath)
	return nil
}

func SaveModInfoToJson(sc types.CliFlags, data interface{}, dir, filename string) (string, error) {

	if err := utils.EnsureDirExists(dir); err != nil {
		return "", err
	}

	fullPath := filepath.Join(dir, fmt.Sprintf("%s.json", filename))

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error formatting data: %s - %v", fullPath, err)
	}

	err = os.WriteFile(fullPath, jsonData, 0644)
	if err != nil {
		return "", fmt.Errorf("error saving file: %s - %v", fullPath, err)
	}

	return fullPath, nil
}
