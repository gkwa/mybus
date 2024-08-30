package core

import (
	"github.com/go-logr/logr"
	"github.com/playwright-community/playwright-go"
)

func Hello(logger logr.Logger) {
	logger.V(1).Info("Debug: Entering Hello function")

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

	browserManager := NewBrowserManager(pw, logger)
	hackerNewsScraper := NewHackerNewsScraper(logger)
	newsScraper := NewNewsScraper(logger, browserManager, hackerNewsScraper)

	if err := newsScraper.ScrapeTopNews(); err != nil {
		logger.Error(err, "Failed to scrape top news")
	}

	logger.Info("Top news stories scraped successfully")
	logger.V(1).Info("Debug: Exiting Hello function")
}
