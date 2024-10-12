package formatters

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"nexus-mods-scraper/internal/types"

	"github.com/PuerkitoBio/goquery"
	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
)

func CleanAndFormatText(input string) string {
	// Remove escape characters and trim quotes
	text := strings.Trim(strings.ReplaceAll(input, "\\n", "\n"), "\"")

	// Split the text into lines and trim each line
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	// Filter out empty lines
	var nonEmptyLines []string
	for _, line := range lines {
		if line != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	// If we have exactly two non-empty lines, join them with a comma and space
	if len(nonEmptyLines) == 2 {
		return strings.Join(nonEmptyLines, ", ")
	}

	// If not, just join all non-empty lines with a space
	return strings.Join(nonEmptyLines, " ")
}

// Helper function to clean and trim text
func CleanTextSelect(s *goquery.Selection) string {
	return strings.TrimSpace(s.Text())
}

// Helper function to clean and trim text
func CleanTextStr(s string) string {
	return strings.TrimSpace(s)
}

func CookieDomain(url string) string {
	// Remove http:// or https://
	re := regexp.MustCompile(`^https?://(www\.)?`)
	// Strip the protocol and www if present
	url = re.ReplaceAllString(url, "")
	// Remove everything after .com, .org, .net, etc.
	reDomain := regexp.MustCompile(`^([a-zA-Z0-9-]+\.[a-zA-Z]{2,})(/.*)?$`)
	matches := reDomain.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1] // Return only the domain part
	}
	return url // Fallback in case regex doesn't match
}

func FormatResultsAsJson(mods types.ModInfo) (string, error) {
	jsonData, err := json.MarshalIndent(mods, "", "    ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal mod information: %w", err)
	}
	return string(jsonData), nil
}

func PrintJson(data string) {
	fmt.Println(data)
}

func PrintPrettyJson(data string, useAltColors ...bool) error {
	var obj interface{}

	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	f := colorjson.NewFormatter()
	f.Indent = 4

	if len(useAltColors) > 0 && useAltColors[0] {
		f.KeyColor = color.New(color.FgHiCyan)
		f.StringColor = color.New(color.FgHiMagenta)
	}

	s, err := f.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to marshal formatted JSON: %w", err)
	}

	fmt.Println(string(s))
	return nil
}

func RemoveHTTPPrefix(url string) string {
	re := regexp.MustCompile(`^https?://`)
	return re.ReplaceAllString(url, "")
}

func StrToInt(input string) (int64, error) {
	result, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, err
	}

	return result, nil
}
