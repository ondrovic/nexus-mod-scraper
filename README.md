![License](https://img.shields.io/badge/license-MIT-blue)
[![releaser](https://github.com/ondrovic/nexus-mods-scraper/actions/workflows/releaser.yml/badge.svg)](https://github.com/ondrovic/nexus-mods-scraper/actions/workflows/releaser.yml)
[![testing](https://github.com/ondrovic/nexus-mods-scraper/actions/workflows/testing.yml/badge.svg)](https://github.com/ondrovic/nexus-mods-scraper/actions/workflows/testing.yml)
[![codecov](https://codecov.io/gh/ondrovic/nexus-mods-scraper/graph/badge.svg?token=RxpgtxqYis)](https://codecov.io/gh/ondrovic/nexus-mods-scraper)
[![Go Report Card](https://goreportcard.com/badge/github.com/ondrovic/nexus-mods-scraper)](https://goreportcard.com/report/github.com/ondrovic/nexus-mods-scraper)
# NexusMods Scraper CLI

A powerful command-line tool to scrape mod information from [https://nexusmods.com](https://nexusmods.com) and return the results in JSON format. This tool also supports extracting cookies from the NexusMods site, which is require to properly scrape the data.

## Requirements

To run the scraper, you need to have a valid `session-cookies.json` file containing your session cookies for NexusMods. You can run the `extract` command and it should grab and save them as long as they exist.

### Example `session-cookies.json` format:

```json
{
  "nexusmods_session": "<value from your session>",
  "nexusmods_session_refresh": "<value from your session>"
}
```

**Important Note:**  
You need to log into your NexusMods account in your browser.

## Installation

To install and run this CLI tool, clone the repository and build the project:

```bash
git clone git@github.com:ondrovic/nexus-mods-scraper.git
cd nexus-mods-scraper
go build -o scraper

-or-

make build
```

## Just run it

To just run the scraper without installing it:

```bash
git clone git@github.com:ondrovic/nexus-mods-scraper.git
cd nexus-mods-scraper
go run nexus-mods-scraper.go
```

## Usage

### Scrape Command

The `scrape` command fetches mod information for a specific game and mod ID from NexusMods and outputs the results in JSON format.

```bash
./nexus-mods-scraper scrape <game-name> <mod-id> [flags]
```

#### Flags:

- `-u, --base-url` (default: `https://nexusmods.com`): Base URL for NexusMods.
- `-d, --cookie-directory` (default: `~/.nexus-mods-scraper/data`): Directory where the cookie file is stored.
- `-f, --cookie-filename` (default: `session-cookies.json`): Filename for the session cookies.
- `-r, --display-results` (default: `false`): Display the results in the terminal.
- `-s, --save-results` (default: `false`): Save the results to a JSON file.
- `-o, --output-directory` (default: `~/.nexus-mods-scraper/data`): Directory where the JSON output will be saved.
- `-c, --valid-cookie-names` (default: `[]string{"nexusmods_session", "nexusmods_session_refresh"}`): Names of the cookies you wish to extract and use.

#### Example:

```bash
./nexus-mods-scraper scrape "skyrim" 12345 --display-results
```

This will fetch mod ID `12345` for the game `Skyrim` and display the results in the terminal.

### Extract Cookies Command

The `extract` command extracts valid cookies for NexusMods and saves them to a JSON file, which is used for authentication in the scraper.

#### Examples:

```bash
./nexus-mods-scraper extract [flags]
```

This will extract the default cookies and save them to the default location.

```bash
./nexus-mods-scraper extract -c "cookie_name","another_cookie","cookie"
```

This will attempt to extract the cookies specified, if found they will be saved in the default location.

#### Flags:

- `-d, --output-directory` (default: `~/.nexus-mods-scraper/data`): Directory where the output file is saved.
- `-f, --output-filename` (default: `session-cookies.json`): Filename to save the session cookies.
- `-c, --valid-cookie-names` (default: `[]string{"nexusmods_session", "nexusmods_session_refresh"}`): Names of the cookies you wish to extract and use.

#### Example:

```bash
./nexus-mods-scraper extract --output-filename my-cookies.json
```

This will extract the cookies and save them as `my-cookies.json`.

## Notes

- You must have valid cookies in your `session-cookies.json` file before scraping.
- Ensure your `session-cookies.json` file is placed in the correct directory or specify the path with the `--cookie-directory` flag.
- Written using [go v1.23.2](https://go.dev/dl/)

## Todo

See [here](TODO)

## Main Packages used

- [goquery](github.com/PuerkitoBio/goquery) - handles the heavy lifting for the scaping
- [colorjson](github.com/TylerBrock/colorjson) - handles making things pretty
- [kooky](github.com/browserutils/kooky) - handles the cookie extraction
- [yacspin](github.com/theckman/yacspin) - spinners
- [cobra](github.com/spf13/cobra) - cli
- [version](go.szostok.io/version) - version command
- [termlink](github.com/savioxavier/termlink) - handles ctrl+click on files
