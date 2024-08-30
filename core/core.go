package core

import (
	"github.com/go-logr/logr"
	"github.com/playwright-community/playwright-go"
)

func Hello(logger logr.Logger, showBrowser bool) {
	logger.V(1).Info("Debug: Entering Hello function")

	if err := playwright.Install(&playwright.RunOptions{Verbose: true}); err != nil {
		logger.Error(err, "Could not install playwright")
		return
	}

	pw, err := playwright.Run()
	if err != nil {
		logger.Error(err, "could not start playwright")
		return
	}
	defer func() {
		if err := pw.Stop(); err != nil {
			logger.Error(err, "failed to stop playwright")
		}
	}()

	browserManager := NewBrowserManager(pw, logger, showBrowser)
	hackerNewsScraper := NewHackerNewsScraper(logger)
	newsScraper := NewNewsScraper(logger, browserManager, hackerNewsScraper)

	if err := newsScraper.ScrapeTopNews(); err != nil {
		logger.Error(err, "Failed to scrape top news")
	}

	logger.Info("Top news stories scraped successfully")
	logger.V(1).Info("Debug: Exiting Hello function")
}
