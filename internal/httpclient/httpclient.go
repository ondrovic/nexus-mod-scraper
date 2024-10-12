package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
)

// HTTPClient is an interface to abstract the Do method for the HTTP client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a global variable holding the HTTP client used for making requests.
var Client HTTPClient

// InitClient initializes an HTTP client with cookie support.
// It creates an HTTP client with a custom CookieJar that manages cookies across requests.
// Returns:
// - An error if any issues occur during client initialization (returns nil in this implementation).
func InitClient(domain, dir, filename string) error {
	// Create a new CookieJar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	// Initialize the HTTP client with the cookie jar
	Client = &http.Client{
		Jar: jar, // Set the CookieJar to manage cookies automatically
	}

	// Call the helper function to set the cookies
	if err := setCookiesFromFile(domain, dir, filename); err != nil {
		return err
	}

	return nil
}

// setCookiesFromFile reads a JSON file containing cookies, combines the directory and filepath,
// and sets the cookies for the specified domain.
func setCookiesFromFile(domain, dir, filename string) error {
	// Combine dir and filename
	cookieFilePath := filepath.Join(dir, filename)

	// Open the JSON file
	file, err := os.Open(cookieFilePath)
	if err != nil {
		return fmt.Errorf("error opening cookie file: %w", err)
	}
	defer file.Close()

	// Create a map to hold cookie key-value pairs
	var cookiesMap map[string]string
	if err := json.NewDecoder(file).Decode(&cookiesMap); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	// Create cookies and set them
	var cookies []*http.Cookie
	for name, value := range cookiesMap {
		cookies = append(cookies, &http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	// Set cookies for the domain
	if jar, ok := Client.(*http.Client).Jar.(*cookiejar.Jar); ok {
		u, err := url.Parse(domain)
		if err != nil {
			return fmt.Errorf("error parsing domain: %w", err)
		}
		jar.SetCookies(u, cookies) // Cookies are set in the jar
	}

	return nil
}
