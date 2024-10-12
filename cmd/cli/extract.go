package cli

import (
	"nexus-mods-scraper/internal/utils/cli"
	"nexus-mods-scraper/internal/utils/exporters"
	"nexus-mods-scraper/internal/utils/extractors"
	"nexus-mods-scraper/internal/utils/formatters"
	"nexus-mods-scraper/internal/utils/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	extractCmd     = &cobra.Command{}
	outputFilename string
)

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

func initExtractFlags(cmd *cobra.Command) {
	cli.RegisterFlag(cmd, "output-directory", "d", storage.GetDataStoragePath(), "Output directory to save the file in", &options.OutputDirectory)
	cli.RegisterFlag(cmd, "output-filename", "f", "session-cookies.json", "Filename to save the session cookies to", &outputFilename)
	cli.RegisterFlag(cmd, "valid-cookie-names", "c", []string{"nexusmods_session", "nexusmods_session_refresh"}, "Names of the cookies to extract", &options.ValidCookies)
}

func ExtractCookies(cmd *cobra.Command, args []string) error {
	domain := formatters.CookieDomain(options.BaseUrl)
	sessionCookies := options.ValidCookies

	extractedCookies, err := extractors.CookieExtractor(domain, sessionCookies)
	if err != nil {
		return err
	}

	if err := exporters.SaveCookiesToJson(options.OutputDirectory, outputFilename, extractedCookies); err != nil {
		return err
	}

	return nil
}
