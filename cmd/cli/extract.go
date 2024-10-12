package cli

import (
	"nexus-mods-scraper/internal/utils/cli"
	"nexus-mods-scraper/internal/utils/exporters"
	"nexus-mods-scraper/internal/utils/extractors"
	"nexus-mods-scraper/internal/utils/formatters"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	extractCmd     = &cobra.Command{}
	outputFilename string
)

func init() {
	extractCmd = &cobra.Command{
		Use:   "extract [validCookies...]",
		Short: "Extract cookies",
		Long:  "Extract cookies for https://nexusmods.com to use with the scraper, will save to json file",
		Args:  cobra.MinimumNArgs(1),
		RunE:  ExtractCookies,
	}

	initExtractFlags(extractCmd)
	viper.BindPFlags(extractCmd.Flags())
	// RootCmd.AddCommand(extractCmd)
}

/*
	TODO: Fix this
	There is a bug where it's not pulling the latest cookies
*/

func initExtractFlags(cmd *cobra.Command) {
	cli.RegisterFlag(cmd, "output-directory", "d", "data", "Output directory to save the file in", &options.OutputDirectory)
	cli.RegisterFlag(cmd, "output-filename", "f", "session-cookies.json", "Filename to save the session cookies to", &outputFilename)
}

func ExtractCookies(cmd *cobra.Command, args []string) error {
	domain := formatters.CookieDomain(options.BaseUrl)
	sessionCookies := args

	extractedCookies, err := extractors.CookieExtractor(domain, sessionCookies)
	if err != nil {
		return err
	}

	if err := exporters.SaveCookiesToJson(options.OutputDirectory, outputFilename, extractedCookies); err != nil {
		return err
	}

	return nil
}
