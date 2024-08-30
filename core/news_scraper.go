package core

import (
	"fmt"

	"github.com/go-logr/logr"
)

type NewsScraper struct {
	logger            logr.Logger
	browserManager    BrowserManagerInterface
	hackerNewsScraper HackerNewsScraperInterface
}

func NewNewsScraper(
	logger logr.Logger,
	browserManager BrowserManagerInterface,
	hackerNewsScraper HackerNewsScraperInterface,
) *NewsScraper {
	return &NewsScraper{
		logger:            logger,
		browserManager:    browserManager,
		hackerNewsScraper: hackerNewsScraper,
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

	stories, err := n.hackerNewsScraper.GetTopStories(page)
	if err != nil {
		return err
	}

	for i, story := range stories {
		n.logger.Info(fmt.Sprintf("%d: %s", i+1, story))
	}

	return nil
}
