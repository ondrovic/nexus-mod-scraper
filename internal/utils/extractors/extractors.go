package extractors

import (
	"errors"

	"fmt"
	"strings"

	"nexus-mods-scraper/internal/types"
	"nexus-mods-scraper/internal/utils/formatters"

	"github.com/PuerkitoBio/goquery"
	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all"
)

// IsAdultContent checks if the mod identified by modId is marked as "Adult content"
// in the goquery document. It looks for an h3 tag with the corresponding modId
// and returns true if the text matches "Adult content".
// AKA your not logged in.
func IsAdultContent(doc *goquery.Document, modId int64) bool {
	// Format the ID of the h3 tag based on the modId
	titleId := fmt.Sprintf("#%d-title", modId)

	// Search for the h3 tag with the constructed ID
	titleTag := doc.Find(titleId)

	// Check if the tag exists and has the correct text
	if titleTag.Length() > 0 {
		titleText := titleTag.Text()
		return titleText == "Adult content"
	}

	return false
}

// CookieExtractor extracts cookies for a given domain from all available browser
// cookie stores. It filters cookies based on the provided validCookies list and
// returns a map of matching cookies. Returns an error if no cookies or stores
// are found, or if reading cookies from a store fails.
func CookieExtractor(domain string, validCookies []string) (map[string]string, error) {
	// Declare a map to store cookies
	cookies := make(map[string]string)

	// Find all available cookie stores (for all browsers)
	cookieStores := kooky.FindAllCookieStores()
	if len(cookieStores) == 0 {
		return nil, errors.New("no cookie stores found")
	}

	// Iterate over each cookie store
	for _, store := range cookieStores {
		defer store.Close()

		// Define filters for valid cookies and specific domain
		var filters = []kooky.Filter{
			kooky.Valid,
			kooky.DomainContains(domain),
		}

		// Read cookies based on the filters
		storeCookies, err := store.ReadCookies(filters...)
		if err != nil {
			// Log the error and continue to the next store
			// log.Printf("Failed to read cookies from store: %v, error: %v", store, err)
			continue
		}

		// Filter and store valid cookies in the map
		for _, cookie := range storeCookies {
			for _, valid := range validCookies {
				if cookie.Name == valid {
					cookies[cookie.Name] = cookie.Value
				}
			}
		}

		// Close the store explicitly after reading its cookies
		store.Close()
	}

	// Check if any cookies were found
	if len(cookies) == 0 {
		return nil, errors.New("no matching cookies found")
	}

	// Return the map of cookies
	return cookies, nil
}

// extractChangeLogs parses a goquery document to extract versioned change logs.
// It looks for specific elements containing version and log notes, and returns
// a slice of ChangeLog objects with the version and corresponding notes.
func extractChangeLogs(doc *goquery.Document) []types.ChangeLog {
	var changeLogs []types.ChangeLog

	doc.Find("#section > div > div.wrap.flex > div:nth-child(2) > div > div.tabcontent.tabcontent-mod-page > div.container.tab-description > div.accordionitems > dl > dd:nth-child(8) > div > ul > li").Each(func(i int, s *goquery.Selection) {
		version := strings.TrimSpace(s.Find("h3").Text())

		var notes []string
		s.Find("div.log-change > ul > li").Each(func(j int, li *goquery.Selection) {
			note := strings.TrimSpace(li.Text())
			if note != "" {
				notes = append(notes, note)
			}
		})

		if version != "" && len(notes) > 0 {
			changeLogs = append(changeLogs, types.ChangeLog{
				Version: version,
				Notes:   notes,
			})
		}
	})

	return changeLogs
}

// extractElementText retrieves and formats the text content of the first element
// matching the provided CSS selector from a goquery document.
func extractElementText(doc *goquery.Document, selector string) string {
	return formatters.CleanAndFormatText(doc.Find(selector).Text())
}

// extractCleanTextExcludingElementText retrieves the text content of the first
// element matching the selector, removes any text from the excluded sub-elements
// matching elem, and returns the cleaned and formatted text.
func extractCleanTextExcludingElementText(doc *goquery.Document, selector, elem string) string {
	selection := doc.Find(selector).First()
	if selection.Length() == 0 {
		return ""
	}

	selection.Find(elem).Remove()
	text := selection.Text()

	return formatters.CleanAndFormatText(text)
}

// ExtractFileInfo parses a goquery document to extract file information, such as
// name, version, upload date, file size, unique downloads, total downloads, and
// description. Returns a slice of File objects with the extracted details.
func ExtractFileInfo(doc *goquery.Document) []types.File {
	fileElements := doc.Find(".file-expander-header")
	files := make([]types.File, 0, fileElements.Length())

	fileElements.Each(func(i int, s *goquery.Selection) {
		file := types.File{
			Name:        formatters.CleanTextSelect(s.Find("p")),
			Version:     formatters.CleanTextSelect(s.Find(".stat-version .stat")),
			UploadDate:  formatters.CleanTextSelect(s.Find(".stat-uploaddate .stat")),
			FileSize:    formatters.CleanTextSelect(s.Find(".stat-filesize .stat")),
			UniqueDLs:   formatters.CleanTextSelect(s.Find(".stat-uniquedls .stat")),
			TotalDLs:    formatters.CleanTextSelect(s.Find(".stat-totaldls .stat")),
			Description: formatters.CleanTextSelect(s.Next().Find(".tabbed-block.files-description")),
		}
		files = append(files, file)
	})

	return files
}

// ExtractModInfo parses a goquery document to extract detailed mod information,
// including name, last updated date, original upload date, creator, changelogs,
// uploader, virus status, short description, full description, tags, dependencies,
// and mods requiring this file. Returns a ModInfo object with the extracted details.
func ExtractModInfo(doc *goquery.Document) types.ModInfo {
	return types.ModInfo{
		Name:             extractElementText(doc, "#pagetitle > h1"),
		LastUpdated:      extractElementText(doc, "#fileinfo > div:nth-child(2) > time"),
		OriginalUpload:   extractElementText(doc, "#fileinfo > div:nth-child(3) > time"),
		Creator:          extractCleanTextExcludingElementText(doc, "#fileinfo > div:nth-child(4)", "h3"),
		ChangeLogs:       extractChangeLogs(doc),
		Uploader:         extractElementText(doc, "#fileinfo > div:nth-child(5) > a"),
		VirusStatus:      extractElementText(doc, "#fileinfo > div:nth-child(6) > div > span"),
		ShortDescription: extractElementText(doc, "#section > div > div.wrap.flex > div:nth-child(2) > div > div.tabcontent.tabcontent-mod-page > div.container.tab-description > p"),
		Description:      extractElementText(doc, "#section > div > div.wrap.flex > div:nth-child(2) > div > div.tabcontent.tabcontent-mod-page > div.container.mod_description_container.condensed"),
		Tags:             extractTags(doc),
		Dependencies:     extractRequirements(doc, "Nexus requirements"),
		ModsUsing:        extractRequirements(doc, "Mods requiring this file"),
	}
}

// extractRequirements parses a goquery document to extract a list of requirements
// from a table with the specified title. It returns a slice of Requirement objects
// containing the name and notes for each requirement. If the table is not found,
// it returns an empty slice.
func extractRequirements(doc *goquery.Document, tableTitle string) []types.Requirement {
	var requirements []types.Requirement

	// Find the correct div.tabbed-block
	block := doc.Find("div.tabbed-block").FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Find("h3").Text() == tableTitle
	}).First()

	if block.Length() == 0 {
		return requirements // Return empty slice if the table is not found
	}

	// Preallocate the slice based on the number of rows
	rowCount := block.Find("table.table.desc-table tbody tr").Length()
	requirements = make([]types.Requirement, 0, rowCount)

	// Extract requirements
	block.Find("table.table.desc-table tbody tr").Each(func(i int, row *goquery.Selection) {
		name := formatters.CleanTextStr(row.Find("td.table-require-name a").Text())
		notes := formatters.CleanTextStr(row.Find("td.table-require-notes").Text())
		requirements = append(requirements, types.Requirement{Name: name, Notes: notes})
	})

	return requirements
}

// extractTags parses a goquery document to extract all tag labels from the tag
// elements on the page. It returns a slice of strings representing the tags.
func extractTags(doc *goquery.Document) []string {
	// Find all tag elements
	elements := doc.Find(".sideitems.side-tags .tags li a span.flex-label")

	// Preallocate the slice
	tags := make([]string, 0, elements.Length())

	// Extract tags
	elements.Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Text())
		if label != "" {
			tags = append(tags, label)
		}
	})

	return tags
}
