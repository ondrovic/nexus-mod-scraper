package cli

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/ondrovic/nexus-mods-scraper/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock structures for each dependency
type Mocker struct {
	mock.Mock
}

// Mock func for httpclient
func (m *Mocker) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *Mocker) SetCookies(u *url.URL, cookies []*http.Cookie) {
	m.Called(u, cookies)
}

func (m *Mocker) Cookies(u *url.URL) []*http.Cookie {
	args := m.Called(u)
	return args.Get(0).([]*http.Cookie)
}

func (m *Mocker) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

var mockFetchDocument = func(_ string) (*goquery.Document, error) {
	html := `<html><body>Mocked HTML content</body></html>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	return doc, nil
}

var mockFetchModInfoConcurrent = func(baseUrl, game string, modId int64, concurrentFetch func(tasks ...func() error) error, fetchDocument func(targetURL string) (*goquery.Document, error)) (types.Results, error) {
	return types.Results{
		Mods: types.ModInfo{
			Name:  "Mocked Mod",
			ModID: modId,
		},
	}, nil
}

// Spinner mocks
func (m *Mocker) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *Mocker) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *Mocker) StopFail() error {
	args := m.Called()
	return args.Error(0)
}

func (m *Mocker) StopFailMessage(msg string) {
	m.Called(msg)
}

// MockUtils implementation for EnsureDirExists
func (m *Mocker) EnsureDirExists(dir string) error {
	args := m.Called(dir)
	return args.Error(0)
}

func TestRun_NoResultsFlagSet(t *testing.T) {
	mockCmd := &cobra.Command{}
	args := []string{"scrape", "game-name", "1234"}

	err := run(mockCmd, args)

	// Assert the expected error
	assert.EqualError(t, err, "at least one of --display-results (-r) or --save-results (-s) must be enabled")
}

func TestRun_InvalidModID(t *testing.T) {
	// Create a new mock command
	mockCmd := &cobra.Command{
		Use:  "scrape",
		RunE: run, // Point to the real `run` function
	}

	// Initialize the scraper flags (as done in `initScrapeFlags`)
	initScrapeFlags(mockCmd)

	// Set the args as if they were passed via command-line
	args := []string{"game", "toast", "--display-results"}

	// Set the args to the command
	mockCmd.SetArgs(args)

	// Execute the command, which will trigger the flag parsing
	err := mockCmd.Execute()

	// Assert that error was returned
	assert.Error(t, err)
	assert.EqualError(t, err, "strconv.ParseInt: parsing \"toast\": invalid syntax")

	// Optionally, you can also assert the `DisplayResults` is set to true
	assert.True(t, viper.GetBool("display-results"))
}

func TestScrapeMod_WithMockedFunctions(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Create a temporary session-cookies.json file
	tempFilePath := filepath.Join(tempDir, "session-cookies.json")
	err := os.WriteFile(tempFilePath, []byte("{}"), 0644) // Create an empty JSON file
	require.NoError(t, err)                               // Ensure the file was created successfully

	// Create a temporary directory for output
	tempOutputDir := filepath.Join(tempDir, "output")
	err = os.Mkdir(tempOutputDir, 0755) // Ensure the output directory is created
	require.NoError(t, err)

	// Prepare test CliFlags with the temporary directories
	sc := types.CliFlags{
		BaseUrl:         "https://somesite.com",
		CookieDirectory: tempDir,
		CookieFile:      "session-cookies.json", // Just the filename, the directory is provided in CookieDirectory
		DisplayResults:  true,
		GameName:        "game",
		ModID:           1234,
		SaveResults:     true,
		OutputDirectory: tempOutputDir, // Use the temporary output directory
	}

	// Act
	err = scrapeMod(sc, mockFetchModInfoConcurrent, mockFetchDocument)

	// Assert
	assert.NoError(t, err)
}
