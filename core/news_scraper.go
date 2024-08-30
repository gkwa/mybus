package core

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
)

type NewsScraper struct {
	logger         logr.Logger
	browserManager BrowserManagerInterface
	siteScraper    SiteScraperInterface
}

func NewNewsScraper(
	logger logr.Logger,
	browserManager BrowserManagerInterface,
	siteScraper SiteScraperInterface,
) *NewsScraper {
	return &NewsScraper{
		logger:         logger,
		browserManager: browserManager,
		siteScraper:    siteScraper,
	}
}

func (n *NewsScraper) ScrapeTopNews() error {
	browser, err := n.browserManager.LaunchBrowser()
	if err != nil {
		return err
	}
	defer func() {
		if err := n.browserManager.CloseBrowser(browser); err != nil {
			n.logger.Error(err, "failed to close browser")
		}
	}()

	page, err := n.browserManager.CreatePage(browser)
	if err != nil {
		return err
	}

	content, err := n.siteScraper.GetContent(page)
	if err != nil {
		return err
	}

	fmt.Println("Scraped content:")
	for _, item := range content {
		fmt.Println(item)
	}

	// Add a delay only if the browser is visible
	if _, ok := n.browserManager.(*BrowserManager); ok && n.browserManager.(*BrowserManager).launchOption.GetLaunchOptions().Headless != nil && !*n.browserManager.(*BrowserManager).launchOption.GetLaunchOptions().Headless {
		time.Sleep(5 * time.Second)
	}

	return nil
}
