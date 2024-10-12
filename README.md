![License](https://img.shields.io/badge/license-MIT-blue)

# NexusMods Scraper CLI

A powerful command-line tool to scrape mod information from [https://nexusmods.com](https://nexusmods.com) and return the results in JSON format. This tool also supports extracting cookies from the NexusMods site, required to authenticate your session and scrape mod data.

## Requirements

To run the scraper, you need to have a valid `session-cookies.json` file containing your session cookies for NexusMods.

### Example `session-cookies.json` format:
```json
{
    "nexusmods_session": "<value from your session>",
    "nexusmods_session_refresh": "<value from your session>"
}
```
**Important Note:**  
You need to log into your NexusMods account in your browser, open the developer tools, and find the values for `nexusmods_session` and `nexusmods_session_refresh` cookies. Insert these values into the `session-cookies.json` file for the scraper to work correctly.

## Installation

To install and run this CLI tool, clone the repository and build the project:
```bash
git clone git@github.com:ondrovic/nexus-mods-scraper.git
cd nexus-mods-scraper
go build -o scraper
```
## Usage

### Scrape Command

The `scrape` command fetches mod information for a specific game and mod ID from NexusMods and outputs the results in JSON format.

```bash
./scraper scrape <game-name> <mod-id> [flags]
```
#### Flags:

- `-u, --base-url` (default: `https://nexusmods.com`): Base URL for NexusMods.
- `-d, --cookie-directory` (default: `data`): Directory where the cookie file is stored.
- `-f, --cookie-filename` (default: `session-cookies.json`): Filename for the session cookies.
- `-r, --display-results` (default: `false`): Display the results in the terminal.
- `-s, --save-results` (default: `false`): Save the results to a JSON file.
- `-o, --output-directory` (default: `data`): Directory where the JSON output will be saved.

#### Example:
```bash
./scraper scrape "skyrim" 12345 --display-results
```
This will fetch mod ID `12345` for the game `Skyrim` and display the results in the terminal.

### Extract Cookies Command

The `extract` command extracts valid cookies for NexusMods and saves them to a JSON file, which is used for authentication in the scraper.

**Important Note:**  
At this time the extract func is disabled, it was extracting expired cookies, so not working properly

```bash
./scraper extract [validCookies...] [flags]
```

#### Flags:

- `-d, --output-directory` (default: `data`): Directory where the output file is saved.
- `-f, --output-filename` (default: `session-cookies.json`): Filename to save the session cookies.

#### Example:
```bash
./scraper extract "nexusmods_session" "nexusmods_session_refresh" --output-filename my-cookies.json
```
This will extract the cookies and save them as `my-cookies.json`.

## Notes

- You must have valid cookies in your `session-cookies.json` file before scraping.
- Ensure your `session-cookies.json` file is placed in the correct directory or specify the path with the `--cookie-directory` flag.
- The extract command isn't enabled

## Todo
See [here](TODO)