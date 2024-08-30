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

	baseURL := "https://dev.to"
	script := `
		const makeAbsolute = (baseURL, relativePath) => {
			if (relativePath.startsWith('http://') || relativePath.startsWith('https://')) {
				return relativePath;
			}
			return new URL(relativePath, baseURL).href;
		};

		const paragraphs = Array.from(document.querySelectorAll('div.article-content p'));
		paragraphs.forEach(p => {
			const links = p.querySelectorAll('a');
			links.forEach(link => {
				link.href = makeAbsolute('` + baseURL + `', link.getAttribute('href'));
			});
		});

		paragraphs.map(p => p.innerHTML);
	`

	result, err := page.Evaluate(script)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate script: %v", err)
	}

	paragraphs, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	content := []string{title}
	for _, p := range paragraphs {
		content = append(content, p.(string))
	}

	return content, nil
}
