package extractors

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ondrovic/nexus-mods-scraper/internal/types"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCookieStore struct {
	mock.Mock
}

// Implement http.CookieJar methods (since CookieStore embeds http.CookieJar)
func (m *MockCookieStore) SetCookies(u *url.URL, cookies []*http.Cookie) {
	// Not needed for our tests, but you can implement if necessary
}

func (m *MockCookieStore) Cookies(u *url.URL) []*http.Cookie {
	// Not needed for our tests, but you can implement if necessary
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

func TestIsAdultContent(t *testing.T) {
	html := `<html><h3 id="12345-title">Adult content</h3></html>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	result := IsAdultContent(doc, 12345)
	assert.True(t, result, "Expected true for adult content")
}

func TestCookieExtractor_Success(t *testing.T) {
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

	// Mock methods that are actually called by CookieExtractor
	mockStore.On("ReadCookies", mock.Anything).Return([]*kooky.Cookie{cookie}, nil)
	mockStore.On("Close").Return(nil)

	// Create a mock function that returns the mock store
	mockStoreProvider := func() []kooky.CookieStore {
		return []kooky.CookieStore{mockStore}
	}

	// Act: Call CookieExtractor with the mock store provider
	result, err := CookieExtractor("example.com", []string{"session"}, mockStoreProvider)

	// Assert: Verify the results
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"session": "1234"}, result)
	mockStore.AssertExpectations(t)
}

func TestCookieExtractor_NoCookieStores(t *testing.T) {
	// Arrange: Mock function that returns no cookie stores
	mockStoreProvider := func() []kooky.CookieStore {
		return []kooky.CookieStore{}
	}

	// Act: Call CookieExtractor with the mock store provider
	result, err := CookieExtractor("example.com", []string{"session"}, mockStoreProvider)

	// Assert: Verify that the correct error is returned
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "no cookie stores found", err.Error())
}

func TestCookieExtractor_NoMatchingCookies(t *testing.T) {
	// Arrange: Create a mock cookie store that returns no matching cookies
	mockStore := new(MockCookieStore)

	// No matching cookies
	mockStore.On("ReadCookies", mock.Anything).Return([]*kooky.Cookie{}, nil)
	mockStore.On("Close").Return(nil)

	// Mock function that returns the mock store
	mockStoreProvider := func() []kooky.CookieStore {
		return []kooky.CookieStore{mockStore}
	}

	// Act: Call CookieExtractor with the mock store provider
	result, err := CookieExtractor("example.com", []string{"session"}, mockStoreProvider)

	// Assert: Verify that the correct error is returned
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "no matching cookies found", err.Error())
}

func TestExtractChangeLogs(t *testing.T) {
	html := `
		<div id="section">
			<div>
				<div class="wrap flex">
					<div>
						<div class="tabcontent tabcontent-mod-page">
							<div class="container tab-description">
								<div class="accordionitems">
									<dl>
										<dd>
											<div>
												<ul>
													<li>
														<h3>v1.0</h3>
														<div class="log-change">
															<ul>
																<li>Initial release</li>
															</ul>
														</div>
													</li>
													<li>
														<h3>v1.1</h3>
														<div class="log-change">
															<ul>
																<li>Fixed bug</li>
																<li>Improved performance</li>
															</ul>
														</div>
													</li>
												</ul>
											</div>
										</dd>
									</dl>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	// Act
	changeLogs := extractChangeLogs(doc)

	// Assert
	expectedChangeLogs := []types.ChangeLog{
		{
			Version: "v1.0",
			Notes:   []string{"Initial release"},
		},
		{
			Version: "v1.1",
			Notes:   []string{"Fixed bug", "Improved performance"},
		},
	}

	assert.Equal(t, expectedChangeLogs, changeLogs)
}

func TestExtractElementText(t *testing.T) {
	html := `<div class="element"> Hello World </div>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	result := extractElementText(doc, ".element")
	assert.Equal(t, "Hello World", result)
}

func TestExtractCleanTextExcludingElementText(t *testing.T) {
	html := `<div class="element"> Hello <span>remove this</span> World </div>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	result := extractCleanTextExcludingElementText(doc, ".element", "span")
	assert.Equal(t, "Hello World", result)
}

func TestExtractFileInfo(t *testing.T) {
	html := `<div class="file-expander-header"><p>File1</p><div class="stat-version"><div class="stat">v1.0</div></div><div class="stat-uploaddate"><div class="stat">2024-01-01</div></div></div>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	result := ExtractFileInfo(doc)
	assert.Len(t, result, 1)
	assert.Equal(t, "File1", result[0].Name)
	assert.Equal(t, "v1.0", result[0].Version)
	assert.Equal(t, "2024-01-01", result[0].UploadDate)
}

func TestExtractModInfo(t *testing.T) {
	html := `<div id="pagetitle" class="clearfix">
				<h1>Mod Name</h1>
			</div>
			<div class="wrap flex">
				<div class="col-1-1 info-details">
					<div id="fileinfo" class="sideitems clearfix">
						<h2>File information</h2>
						<div class="sideitem timestamp">
							<h3>Last updated</h3>
							<time datetime="2024-10-13 10:44">
								<span class="date">13 October 2024</span>
								<span class="time">10:44AM</span>
							</time>
						</div>
						<div class="sideitem timestamp">
							<h3>Original upload</h3>
							<time datetime="2024-10-13 10:44">
								<span class="date">13 October 2024</span>
								<span class="time">10:44AM</span>
							</time>
						</div>
						<div class="sideitem">
							<h3>Created by</h3>
							Mod Creator
						</div>
						<div class="sideitem">
							<h3>Uploaded by</h3>
							<a href="https://www.somesite.com/somegame/someuser/1234">Uploader Name</a>
						</div>
						<div class="sideitem">
							<h3>Virus scan</h3>
							<div class="result  inline-flex" style="height: 25px; position: relative; top: 5px;">
								<svg title="" class="icon icon-exclamation">
									<use xlink:href="https://www.somesite.com/assets/images/icons/icons.svg#icon-exclamation">
									</use>
								</svg> <span class="flex-label">
									Some files not scanned </span>
							</div>
						</div>
					</div>
					<div class="sideitems side-tags">
						<h2>Tags for this mod</h2>
						<div class="sideitem clearfix">
							<ul class="tags">
								<span></span><span class="js-hidable-tags hidden"></span>
							</ul>
						</div>
					</div>
				</div>
			</div>`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	// Act
	result := ExtractModInfo(doc)

	// Assert
	expectedModInfo := types.ModInfo{
		Name:           "Mod Name",
		LastUpdated:    "13 October 2024, 10:44AM",
		OriginalUpload: "13 October 2024, 10:44AM",
		Creator:        "Mod Creator",
		Uploader:       "Uploader Name",
		VirusStatus:    "Some files not scanned",
		Tags:           []string{},
	}

	assert.Equal(t, expectedModInfo, result)
}

func TestExtractRequirements(t *testing.T) {
	html := `
		<div class="tabbed-block">
			<h3>Nexus requirements</h3>
			<table class="table desc-table">
				<thead>
					<tr>
						<th class="table-require-name header headerSortUp"><span class="table-header">Mod name</span></th>
						<th class="table-require-notes"><span class="table-header">Notes</span></th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td class="table-require-name">
							<a href="https://www.site.com/mod/1234">Requirement1</a>
						</td>
						<td class="table-require-notes">Note1</td>
					</tr>
				</tbody>
			</table>
		</div>`

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	// Act
	result := extractRequirements(doc, "Nexus requirements")

	// Assert
	assert.Len(t, result, 1, "Expected 1 requirement")
	assert.Equal(t, "Requirement1", result[0].Name)
	assert.Equal(t, "Note1", result[0].Notes)
}

func TestExtractTags(t *testing.T) {
	html := `<div class="sideitems side-tags"><ul class="tags"><li><a><span class="flex-label">Tag1</span></a></li></ul></div>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	result := extractTags(doc)
	assert.Len(t, result, 1)
	assert.Equal(t, "Tag1", result[0])
}
