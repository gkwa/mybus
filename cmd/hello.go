package cmd

import (
	"github.com/gkwa/mybus/core"
	"github.com/spf13/cobra"
)

var (
	showBrowser bool
	site        string
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Scrape content from specified site",
	Long:  `Scrape content from either Hacker News or a specific Dev.to article.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.Info("Running hello command")
		core.Hello(logger, showBrowser, site)
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
	helloCmd.Flags().BoolVar(&showBrowser, "show-browser", false, "Show browser during navigation")
	helloCmd.Flags().StringVar(&site, "site", "hacker-news", "Site to scrape (hacker-news or dev-to)")
}
