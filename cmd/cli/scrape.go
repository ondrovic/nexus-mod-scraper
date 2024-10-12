package cli

import (
	"fmt"

	"github.com/savioxavier/termlink"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"nexus-mods-scraper/internal/fetchers"
	"nexus-mods-scraper/internal/httpclient"
	"nexus-mods-scraper/internal/types"
	"nexus-mods-scraper/internal/utils"
	"nexus-mods-scraper/internal/utils/cli"
	"nexus-mods-scraper/internal/utils/exporters"
	"nexus-mods-scraper/internal/utils/formatters"
	"nexus-mods-scraper/internal/utils/spinners"
	"nexus-mods-scraper/internal/utils/storage"

	"path/filepath"
	"strings"
)

var (
	// options holds the command-line flag values using the CliFlags struct.
	options = types.CliFlags{}
	// scrapeCmd is a Cobra command used for scraping operations in the application.
	scrapeCmd = &cobra.Command{}
)

// init initializes the scrape command with usage, description, and argument validation.
// It binds flags using Viper and adds the command to the root command for execution.
func init() {
	scrapeCmd = &cobra.Command{
		Use:   "scrape <game name> <mod id> [flags]",
		Short: "Scape mod",
		Long:  "Scrape mod for game and returns a JSON output",
		Args:  cobra.ExactArgs(2),
		RunE:  run,
	}

	initScrapeFlags(scrapeCmd)
	viper.BindPFlags(scrapeCmd.Flags())
	RootCmd.AddCommand(scrapeCmd)
}

// initScrapeFlags registers the command-line flags for the scrape command, including
// options for the base URL, cookie directory, cookie filename, result display and save
// options, output directory, and valid cookie names. It binds these flags to the
// corresponding fields in the CliFlags struct.
func initScrapeFlags(cmd *cobra.Command) {
	cli.RegisterFlag(cmd, "base-url", "u", "https://nexusmods.com", "Base url for the mods", &options.BaseUrl)
	cli.RegisterFlag(cmd, "cookie-directory", "d", storage.GetDataStoragePath(), "Directory your cookie file is stored in", &options.CookieDirectory)
	cli.RegisterFlag(cmd, "cookie-filename", "f", "session-cookies.json", "Filename where the cookies are stored", &options.CookieFile)
	cli.RegisterFlag(cmd, "display-results", "r", false, "Do you want to display the results in the terminal?", &options.DisplayResults)
	cli.RegisterFlag(cmd, "save-results", "s", false, "Do you want to save the results to a JSON file?", &options.SaveResults)
	cli.RegisterFlag(cmd, "output-directory", "o", storage.GetDataStoragePath(), "Output directory to save files", &options.OutputDirectory)
	cli.RegisterFlag(cmd, "valid-cookie-names", "c", []string{"nexusmods_session", "nexusmods_session_refresh"}, "Names of the cookies to extract", &options.ValidCookies)
}

// run executes the scrape command, validating that either display or save results
// options are enabled. It parses the mod ID and game name from the arguments, reads
// the configuration values from Viper, and then calls the scrapeMod function with
// the populated CliFlags.
func run(cmd *cobra.Command, args []string) error {
	if !options.DisplayResults && !options.SaveResults {
		return fmt.Errorf("at least one of --display-results (-r) or --save-results (-s) must be enabled")
	}
	modId, err := formatters.StrToInt(args[1])
	if err != nil {
		return err
	}

	scraper := types.CliFlags{
		BaseUrl:         viper.GetString("base-url"),
		CookieDirectory: viper.GetString("cookie-directory"),
		CookieFile:      viper.GetString("cookie-filename"),
		DisplayResults:  viper.GetBool("display-results"),
		GameName:        args[0],
		ModId:           modId,
		SaveResults:     viper.GetBool("save-results"),
		OutputDirectory: viper.GetString("output-directory"),
		ValidCookies:    viper.GetStringSlice("valid-cookie-names"),
	}

	return scrapeMod(scraper)
}

// scrapeMod orchestrates the process of scraping mod information, including setting up
// the HTTP client, scraping mod info, displaying results, and saving results based on
// the provided command-line flags. It utilizes spinners to indicate progress throughout
// the operations and returns an error if any step fails.
func scrapeMod(sc types.CliFlags) error {
	// Create and start the main spinner for HTTP client setup
	httpSpinner := spinners.CreateSpinner("Setting up HTTP client", "✓", "HTTP client setup complete", "✗", "HTTP client setup failed")
	if err := httpSpinner.Start(); err != nil {
		return fmt.Errorf("failed to start spinner: %w", err)
	}

	// HTTP Client Setup
	if err := httpclient.InitClient(sc.BaseUrl, sc.CookieDirectory, sc.CookieFile); err != nil {
		httpSpinner.StopFailMessage(fmt.Sprintf("Error setting up HTTP client: %v", err))
		httpSpinner.StopFail()
		return err
	}
	httpSpinner.Stop()

	// Create and start the spinner for scraping mod info
	scrapeSpinner := spinners.CreateSpinner(fmt.Sprintf("Scraping modId: %d for game: %s", sc.ModId, sc.GameName), "✓", "Mod scraping complete", "✗", "Mod scraping failed")
	if err := scrapeSpinner.Start(); err != nil {
		return fmt.Errorf("failed to start spinner: %w", err)
	}

	// Scrape Mod Info
	results, err := fetchers.FetchModInfoConcurrent(sc.BaseUrl, sc.GameName, sc.ModId)
	if err != nil {
		scrapeSpinner.StopFailMessage(fmt.Sprintf("Error scraping mod: %v", err))
		scrapeSpinner.StopFail()
		return err
	}
	scrapeSpinner.Stop()

	// Display Results
	if sc.DisplayResults {
		displaySpinner := spinners.CreateSpinner("Displaying results", "✓", "Results displayed", "✗", "Failed to display results")
		if err := displaySpinner.Start(); err != nil {
			return fmt.Errorf("failed to start display spinner: %w", err)
		}
		displaySpinner.Stop() // Temporarily stop spinner for clean output

		// Print the results
		if err := exporters.DisplayResults(sc, results); err != nil {
			fmt.Println("Error displaying results:", err)
			displaySpinner.StopFail()
			return err
		}
		displaySpinner.Stop() // Restart the spinner after results are displayed
	}

	// Save Results
	if sc.SaveResults {
		saveSpinner := spinners.CreateSpinner("Saving results", "✓", "Results saved successfully", "✗", "Failed to save results")
		if err := saveSpinner.Start(); err != nil {
			return fmt.Errorf("failed to start save spinner: %w", err)
		}

		outputGameDirectory := filepath.Join(sc.OutputDirectory, strings.ToLower(sc.GameName))
		if err := utils.EnsureDirExists(outputGameDirectory); err != nil {
			saveSpinner.StopFailMessage(fmt.Sprintf("Error creating directory: %v", err))
			saveSpinner.StopFail()
			return err
		}

		outputFilename := fmt.Sprintf("%s %d", strings.ToLower(results.Mods.Name), results.Mods.ModId)
		if item, err := exporters.SaveModInfoToJson(sc, results, outputGameDirectory, outputFilename); err != nil {
			saveSpinner.StopFailMessage(fmt.Sprintf("Error saving results: %v", err))
			saveSpinner.StopFail()
			return err
		} else {
			// saveSpinner.StopMessage(fmt.Sprintf("Saved successfully to %s", item))
			saveSpinner.StopMessage(fmt.Sprintf("Saved successfully to %s", termlink.ColorLink(item, item, "green")))
		}
		saveSpinner.Stop()
	}

	return nil
}
