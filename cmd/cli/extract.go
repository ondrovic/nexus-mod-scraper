package cli

import (
	"os"

	"github.com/browserutils/kooky"
	"github.com/ondrovic/nexus-mods-scraper/internal/utils"
	"github.com/ondrovic/nexus-mods-scraper/internal/utils/cli"
	"github.com/ondrovic/nexus-mods-scraper/internal/utils/exporters"
	"github.com/ondrovic/nexus-mods-scraper/internal/utils/extractors"
	"github.com/ondrovic/nexus-mods-scraper/internal/utils/formatters"
	"github.com/ondrovic/nexus-mods-scraper/internal/utils/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// extractCmd is a Cobra command used for extracting information within the application.
	extractCmd = &cobra.Command{}
	// outputFilename is a string variable that stores the name of the file to which
	// output will be saved.
	outputFilename string
)

// init initializes the extract command, setting its usage, description, and argument validation.
// It binds flags using Viper and adds the extract command to the root command for extracting
// cookies and saving them to a JSON file.
func init() {
	extractCmd = &cobra.Command{
		Use:   "extract",
		Short: "Extract cookies",
		Long:  "Extract cookies for https://nexusmods.com to use with the scraper, will save to json file",
		Args:  cobra.NoArgs,
		RunE:  ExtractCookies,
	}

	initExtractFlags(extractCmd)
	viper.BindPFlags(extractCmd.Flags())
	RootCmd.AddCommand(extractCmd)
}

// initExtractFlags registers the command-line flags for the extract command, including
// options for the output directory, output filename, and valid cookie names to extract.
// These flags are bound to the corresponding variables and fields in CliFlags.
func initExtractFlags(cmd *cobra.Command) {
	cli.RegisterFlag(cmd, "output-directory", "d", storage.GetDataStoragePath(), "Output directory to save the file in", &options.OutputDirectory)
	cli.RegisterFlag(cmd, "output-filename", "f", "session-cookies.json", "Filename to save the session cookies to", &outputFilename)
	cli.RegisterFlag(cmd, "valid-cookie-names", "c", []string{"nexusmods_session", "nexusmods_session_refresh"}, "Names of the cookies to extract", &options.ValidCookies)
}

// ExtractCookies extracts cookies from the specified domain using the valid cookie names,
// then saves them as a JSON file in the designated output directory. Returns an error
// if cookie extraction or saving fails.
func ExtractCookies(cmd *cobra.Command, args []string) error {
	domain := formatters.CookieDomain(options.BaseUrl)
	sessionCookies := options.ValidCookies

	extractedCookies, err := extractors.CookieExtractor(domain, sessionCookies, kooky.FindAllCookieStores)
	if err != nil {
		return err
	}

	if err := exporters.SaveCookiesToJson(options.OutputDirectory, outputFilename, extractedCookies, os.OpenFile, utils.EnsureDirExists); err != nil {
		return err
	}

	return nil
}
