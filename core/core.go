package core

import (
	"fmt"

	"github.com/go-logr/logr"
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

type HackerNewsScraper struct {
	logger logr.Logger
}

func NewHackerNewsScraper(logger logr.Logger) *HackerNewsScraper {
	return &HackerNewsScraper{logger: logger}
}

func (h *HackerNewsScraper) GetTopStories(page playwright.Page) ([]string, error) {
	if _, err := page.Goto("https://news.ycombinator.com"); err != nil {
		return nil, fmt.Errorf("could not goto: %v", err)
	}

	entries, err := page.Locator(".athing").All()
	if err != nil {
		return nil, fmt.Errorf("could not get entries: %v", err)
	}

	var stories []string
	for _, entry := range entries {
		title, err := entry.Locator("td.title > span > a").TextContent()
		if err != nil {
			return nil, fmt.Errorf("could not get text content: %v", err)
		}
		stories = append(stories, title)
	}

	return stories, nil
}

type BrowserManager struct {
	pw     *playwright.Playwright
	logger logr.Logger
}

func NewBrowserManager(pw *playwright.Playwright, logger logr.Logger) *BrowserManager {
	return &BrowserManager{pw: pw, logger: logger}
}

func (b *BrowserManager) LaunchBrowser() (playwright.Browser, error) {
	browser, err := b.pw.Chromium.Launch()
	if err != nil {
		return nil, fmt.Errorf("could not launch browser: %v", err)
	}
	return browser, nil
}

func (b *BrowserManager) CreatePage(browser playwright.Browser) (playwright.Page, error) {
	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %v", err)
	}
	return page, nil
}

func (b *BrowserManager) CloseBrowser(browser playwright.Browser) error {
	if err := browser.Close(); err != nil {
		return fmt.Errorf("could not close browser: %v", err)
	}
	return nil
}

type NewsScraper struct {
	logger            logr.Logger
	browserManager    BrowserManagerInterface
	hackerNewsScraper HackerNewsScraperInterface
}

func NewNewsScraper(logger logr.Logger, browserManager BrowserManagerInterface, hackerNewsScraper HackerNewsScraperInterface) *NewsScraper {
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
	defer n.browserManager.CloseBrowser(browser)

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

func Hello(logger logr.Logger) {
	logger.V(1).Info("Debug: Entering Hello function")

	pw, err := playwright.Run()
	if err != nil {
		logger.Error(err, "could not start playwright")
		return
	}
	defer pw.Stop()

	browserManager := NewBrowserManager(pw, logger)
	hackerNewsScraper := NewHackerNewsScraper(logger)
	newsScraper := NewNewsScraper(logger, browserManager, hackerNewsScraper)

	if err := newsScraper.ScrapeTopNews(); err != nil {
		logger.Error(err, "Failed to scrape top news")
	}

	logger.Info("Top news stories scraped successfully")
	logger.V(1).Info("Debug: Exiting Hello function")
}
