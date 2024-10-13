package fetchers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ondrovic/nexus-mods-scraper/internal/httpclient"
	"github.com/ondrovic/nexus-mods-scraper/internal/types"
	"github.com/ondrovic/nexus-mods-scraper/internal/utils/extractors"

	"github.com/PuerkitoBio/goquery"
)

// FetchModInfoConcurrent retrieves mod information and file details concurrently
// for a specified mod ID and game. It validates URLs and uses provided functions
// for concurrent fetching of mod info and file info extraction. The results are populated
// in the Results struct, and an error is returned if any fetching or extraction step fails.
func FetchModInfoConcurrent(baseUrl, game string, modId int64, concurrentFetch func(tasks ...func() error) error, fetchDocument func(targetURL string) (*goquery.Document, error)) (types.Results, error) {
	modUrl := fmt.Sprintf("%s/%s/mods/%d", baseUrl, game, modId)

	// Validate the initial URL
	if _, err := url.Parse(modUrl); err != nil {
		return types.Results{}, err
	}

	var results types.Results

	// Function to handle mod info fetch
	err := concurrentFetch(
		func() error {
			doc, err := fetchDocument(modUrl)
			if err != nil {
				return err
			}

			if extractors.IsAdultContent(doc, modId) {
				return fmt.Errorf("adult content detected, cookies not working")
			}

			results.Mods = extractors.ExtractModInfo(doc)
			results.Mods.ModID = modId
			results.Mods.LastChecked = time.Now()
			return nil
		},
		func() error {
			filesTabURL := fmt.Sprintf("%s?tab=files", modUrl)

			// Validate files tab URL
			if _, err := url.Parse(filesTabURL); err != nil {
				return err
			}

			filesDoc, err := fetchDocument(filesTabURL)
			if err != nil {
				return err
			}

			results.Mods.Files = extractors.ExtractFileInfo(filesDoc)
			if len(results.Mods.Files) > 0 {
				results.Mods.LatestVersion = results.Mods.Files[0].Version
			}

			return nil
		},
	)

	if err != nil {
		return types.Results{}, err
	}

	return results, nil
}

// FetchDocument sends an HTTP GET request to the target URL, manually attaches cookies
// from the HTTP client's cookie jar, and returns the response as a parsed goquery document.
// It ensures a successful 200 OK status before parsing and returns an error if the request
// or document parsing fails.
func FetchDocument(targetURL string) (*goquery.Document, error) {
	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	// Manually retrieve cookies for the domain
	u, _ := url.Parse(targetURL)
	cookies := httpclient.Client.(*http.Client).Jar.Cookies(u)

	// Build the Cookie header string manually from the cookies
	var cookieHeader []string
	for _, cookie := range cookies {
		cookieHeader = append(cookieHeader, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	req.Header.Set("Cookie", strings.Join(cookieHeader, "; "))

	// Check the request headers (should include the Cookie header now)
	// fmt.Printf("Request Headers: %v\n", req.Header)

	// Use the global httpclient.Client to make the request
	resp, err := httpclient.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("Response Headers: %v\n", resp.Header)

	defer resp.Body.Close()

	// Ensure we received a 200 OK response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch document: %s returned %d", targetURL, resp.StatusCode)
	}

	// Parse the response body into a goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Return the goquery document
	return doc, nil
}
