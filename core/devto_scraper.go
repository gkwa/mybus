package core

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/playwright-community/playwright-go"
)

type DevToScraper struct {
	logger logr.Logger
}

func NewDevToScraper(logger logr.Logger) *DevToScraper {
	return &DevToScraper{logger: logger}
}

func (d *DevToScraper) GetContent(page playwright.Page) ([]string, error) {
	if _, err := page.Goto("https://dev.to/lucasnevespereira/build-your-own-linktree-with-go-and-github-pages-3fha"); err != nil {
		return nil, fmt.Errorf("could not goto: %v", err)
	}

	title, err := page.Locator("h1").TextContent()
	if err != nil {
		return nil, fmt.Errorf("could not get title: %v", err)
	}

	paragraphs, err := page.Locator("div.article-content p").All()
	if err != nil {
		return nil, fmt.Errorf("could not get paragraphs: %v", err)
	}

	content := []string{title}
	for _, p := range paragraphs {
		text, err := p.TextContent()
		if err != nil {
			return nil, fmt.Errorf("could not get paragraph text: %v", err)
		}
		content = append(content, text)
	}

	return content, nil
}
