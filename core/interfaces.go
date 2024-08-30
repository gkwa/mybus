package core

import (
	"github.com/playwright-community/playwright-go"
)

type NewsScraperInterface interface {
	ScrapeTopNews() error
}

type HackerNewsScraperInterface interface {
	GetTopStories(page playwright.Page) ([]string, error)
}

type BrowserManagerInterface interface {
	LaunchBrowser() (playwright.Browser, error)
	CreatePage(browser playwright.Browser) (playwright.Page, error)
	CloseBrowser(browser playwright.Browser) error
}
