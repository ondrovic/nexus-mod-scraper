package cli

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/chelnak/ysmrr/pkg/animations"
	"github.com/chelnak/ysmrr/pkg/colors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"nexus-mods-scraper/internal/fetchers"
	"nexus-mods-scraper/internal/httpclient"
	"nexus-mods-scraper/internal/types"
	"nexus-mods-scraper/internal/utils"
	"nexus-mods-scraper/internal/utils/cli"
	"nexus-mods-scraper/internal/utils/exporters"
	"nexus-mods-scraper/internal/utils/formatters"
	"nexus-mods-scraper/internal/utils/storage"

	"path/filepath"
	"strings"
)

var (
	options   = types.CliFlags{}
	scrapeCmd = &cobra.Command{}

	sm ysmrr.SpinnerManager
)

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

func initScrapeFlags(cmd *cobra.Command) {
	cli.RegisterFlag(cmd, "base-url", "u", "https://nexusmods.com", "Base url for the mods", &options.BaseUrl)
	cli.RegisterFlag(cmd, "cookie-directory", "d", "data", "Directory your cookie file is stored in", &options.CookieDirectory)
	cli.RegisterFlag(cmd, "cookie-filename", "f", "session-cookies.json", "Filename where the cookies are stored", &options.CookieFile)
	cli.RegisterFlag(cmd, "display-results", "r", false, "Do you want to display the results in the terminal?", &options.DisplayResults)
	cli.RegisterFlag(cmd, "save-results", "s", false, "Do you want to save the results to a JSON file?", &options.SaveResults)
	cli.RegisterFlag(cmd, "output-directory", "o", storage.GetDataStoragePath(), "Output directory to save files", &options.OutputDirectory)
	cli.RegisterFlag(cmd, "valid-cookie-names", "c", []string{"nexusmods_session", "nexusmods_session_refresh"}, "Names of the cookies to extract", &options.ValidCookies)
}

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

func init() {
	sm = ysmrr.NewSpinnerManager(
		ysmrr.WithAnimation(animations.Dots),
		ysmrr.WithSpinnerColor(colors.FgHiBlue),
	)
}

func scrapeMod(sc types.CliFlags) error {
	sm.Start()
	defer sm.Stop()

	spnMessages := sm.AddSpinner("Setting up httpclient")

	if err := httpclient.InitClient(sc.BaseUrl, sc.CookieDirectory, sc.CookieFile); err != nil {
		errMessage := "Error setting up httpclient"
		if strings.Contains(err.Error(), "cannot find") {
			errMessage += ", session-cookies.json missing"
		}

		spnMessages.ErrorWithMessagef("%s %v", errMessage, err)

		return nil
	}

	spnMessages.UpdateMessagef("Scraping modId: %d for game: %s", sc.ModId, sc.GameName)

	results, err := fetchers.FetchModInfoConcurrent(sc.BaseUrl, sc.GameName, sc.ModId)
	if err != nil {
		spnMessages.ErrorWithMessagef("Error Scraping mod %v", err)
	}

	if sc.DisplayResults {
		spnMessages.UpdateMessage("Displaying results")

		if err := exporters.DisplayResults(sc, results); err != nil {
			spnMessages.ErrorWithMessagef("Error displaying results %v", err)
			return nil
		}
	}

	if sc.SaveResults {
		spnMessages.UpdateMessage("Saving results")

		outputGameDirectory := filepath.Join(sc.OutputDirectory, strings.ToLower(sc.GameName))

		// Check if the directory exists, if not create it
		if err := utils.EnsureDirExists(outputGameDirectory); err != nil {
			return err
		}

		outputFilename := fmt.Sprintf("%s %d", strings.ToLower(results.Mods.Name), results.Mods.ModId)

		if item, err := exporters.SaveModInfoToJson(sc, results, outputGameDirectory, outputFilename); err != nil {
			spnMessages.ErrorWithMessagef("Error saving results %v", err)
		} else {
			spnMessages.CompleteWithMessagef("Saved successfully to %s", item)
		}

	} else {
		spnMessages.CompleteWithMessage("Scraping complete")
	}

	return nil
}
