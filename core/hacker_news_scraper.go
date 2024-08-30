package core

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/playwright-community/playwright-go"
)

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
