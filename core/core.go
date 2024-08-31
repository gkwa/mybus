package core

import (
	"github.com/go-logr/logr"
	"github.com/playwright-community/playwright-go"
)

func Hello(logger logr.Logger, showBrowser bool, site string, showLinks bool) {
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
	var siteScraper SiteScraperInterface

	switch site {
	case "hacker-news":
		siteScraper = NewHackerNewsScraper(logger)
	case "dev-to":
		siteScraper = NewDevToScraper(logger, showLinks, showBrowser)
	default:
		logger.Error(nil, "Invalid site specified")
		return
	}

	newsScraper := NewNewsScraper(logger, browserManager, siteScraper)

	if err := newsScraper.ScrapeTopNews(); err != nil {
		logger.Error(err, "Failed to scrape content")
	}

	logger.Info("Content scraped successfully")
	logger.V(1).Info("Debug: Exiting Hello function")
}
