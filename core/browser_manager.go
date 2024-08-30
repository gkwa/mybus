package core

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/playwright-community/playwright-go"
)

type BrowserManager struct {
	pw           *playwright.Playwright
	logger       logr.Logger
	launchOption BrowserLaunchOption
}

type BrowserLaunchOption interface {
	GetLaunchOptions() playwright.BrowserTypeLaunchOptions
}

type HeadlessBrowserLaunch struct{}

func (h HeadlessBrowserLaunch) GetLaunchOptions() playwright.BrowserTypeLaunchOptions {
	return playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	}
}

type VisibleBrowserLaunch struct{}

func (v VisibleBrowserLaunch) GetLaunchOptions() playwright.BrowserTypeLaunchOptions {
	return playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	}
}

func NewBrowserManager(pw *playwright.Playwright, logger logr.Logger, showBrowser bool) *BrowserManager {
	var launchOption BrowserLaunchOption
	if showBrowser {
		launchOption = VisibleBrowserLaunch{}
	} else {
		launchOption = HeadlessBrowserLaunch{}
	}
	return &BrowserManager{pw: pw, logger: logger, launchOption: launchOption}
}

func (b *BrowserManager) LaunchBrowser() (playwright.Browser, error) {
	browser, err := b.pw.Chromium.Launch(b.launchOption.GetLaunchOptions())
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
