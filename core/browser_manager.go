package core

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/playwright-community/playwright-go"
)

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
