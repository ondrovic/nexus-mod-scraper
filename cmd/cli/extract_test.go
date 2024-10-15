package cli

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all"
	"github.com/ondrovic/nexus-mods-scraper/internal/utils/exporters"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCookieStore struct {
	mock.Mock
}

// Implement http.CookieJar methods (since CookieStore embeds http.CookieJar)
func (m *MockCookieStore) SetCookies(u *url.URL, cookies []*http.Cookie) {
	// Needed for our tests, but you can implement if necessary
}

func (m *MockCookieStore) Cookies(u *url.URL) []*http.Cookie {
	// Needed for our tests, but you can implement if necessary
	return nil
}

// Mock the SubJar method
func (m *MockCookieStore) SubJar(filters ...kooky.Filter) (http.CookieJar, error) {
	args := m.Called(filters)
	return args.Get(0).(http.CookieJar), args.Error(1)
}

// Mock the ReadCookies method
func (m *MockCookieStore) ReadCookies(filters ...kooky.Filter) ([]*kooky.Cookie, error) {
	args := m.Called(filters)
	return args.Get(0).([]*kooky.Cookie), args.Error(1)
}

// Mock the Browser method
func (m *MockCookieStore) Browser() string {
	args := m.Called()
	return args.String(0)
}

// Mock the Profile method
func (m *MockCookieStore) Profile() string {
	args := m.Called()
	return args.String(0)
}

// Mock the IsDefaultProfile method
func (m *MockCookieStore) IsDefaultProfile() bool {
	args := m.Called()
	return args.Bool(0)
}

// Mock the FilePath method
func (m *MockCookieStore) FilePath() string {
	args := m.Called()
	return args.String(0)
}

// Mock the Close method
func (m *MockCookieStore) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCookieStore) CookieExtractor(domain string, validCookies []string, storeFinder func() []kooky.CookieStore) (map[string]string, error) {
	args := m.Called(domain, validCookies, storeFinder)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockCookieStore) SaveCookiesToJson(outputDir, filename string, data map[string]string, openFile func(string, int, os.FileMode) (*os.File, error), ensureDirExists func(string) error) error {
	args := m.Called(outputDir, filename, data, openFile, ensureDirExists)
	return args.Error(0)
}

func TestExtractCookies_Success(t *testing.T) {
	// Arrange: Create a mock cookie store
	mockStore := new(MockCookieStore)

	// Define a mock cookie
	cookie := &kooky.Cookie{
		Cookie: http.Cookie{
			Name:   "session",
			Value:  "1234",
			Domain: "example.com",
		},
		Creation:  time.Now(),
		Container: "MockBrowser",
	}

	// Mock methods that are called by CookieExtractor
	mockStore.On("ReadCookies", mock.Anything).Return([]*kooky.Cookie{cookie}, nil)
	mockStore.On("Close").Return(nil)

	// Create a mock store provider to avoid using live cookie stores
	mockStoreProvider := func() []kooky.CookieStore {
		return []kooky.CookieStore{mockStore}
	}

	// Mock the `openFileFunc` and `ensureDirExistsFunc`
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, "session-cookies.json")

	mockOpenFile := func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return os.OpenFile(tempFilePath, flag, perm)
	}

	mockEnsureDirExists := func(dir string) error {
		return nil // Simulate directory existence or creation
	}

	// Set the options (these can be set globally or adjusted as necessary)
	options.BaseUrl = "http://example.com"
	options.ValidCookies = []string{"session"}
	options.OutputDirectory = tempDir
	outputFilename = "session-cookies.json"

	// Act: Call ExtractCookies using the mockStoreProvider
	cmd := &cobra.Command{}
	args := []string{}
	err := ExtractCookies(cmd, args, mockStoreProvider)

	// Call SaveCookiesToJson with mocked functions
	err = exporters.SaveCookiesToJson(options.OutputDirectory, outputFilename, map[string]string{"session": "1234"}, mockOpenFile, mockEnsureDirExists)

	// Assert: Verify no error and that all expectations on the mocks are met
	assert.NoError(t, err)
	mockStore.AssertExpectations(t)

	// Verify the contents of the temp file
	fileContent, err := os.ReadFile(tempFilePath)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	expectedContent := `{
    "session": "1234"
}`

	assert.JSONEq(t, expectedContent, string(fileContent), "The cookie data written to the file is not as expected")
}

func TestExtractCookies_ErrorInCookieExtractor(t *testing.T) {
	// Arrange: Create a mock cookie store
	mockStore := new(MockCookieStore)

	// Mock store provider to return the mock store
	mockStoreProvider := func() []kooky.CookieStore {
		return []kooky.CookieStore{mockStore}
	}

	// Simulate error in CookieExtractor
	mockStore.On("CookieExtractor", "example.com", []string{"session"}, mock.Anything).Return(nil, errors.New)

	// Mock ReadCookies and Close (since they are called internally)
	mockStore.On("ReadCookies", mock.Anything).Return([]*kooky.Cookie{}, nil) // Return empty slice instead of nil
	mockStore.On("Close").Return(nil)                                         // Simulate successful closing

	// Set the options
	options.BaseUrl = "http://example.com"
	options.ValidCookies = []string{"session"}
	options.OutputDirectory = "/tmp"
	outputFilename = "session-cookies.json"

	// Act: Call ExtractCookies using the mockStoreProvider
	cmd := &cobra.Command{}
	args := []string{}
	err := ExtractCookies(cmd, args, mockStoreProvider)

	// Assert: Verify the error from CookieExtractor is returned
	assert.Error(t, err)
	assert.Equal(t, "no matching cookies found", err.Error())
}
