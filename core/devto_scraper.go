package core

import (
	"fmt"
	"os"
	"strings"

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
	d.logger.V(1).Info("Starting GetContent function")

	d.logger.V(1).Info("Navigating to page")
	if _, err := page.Goto("https://dev.to/lucasnevespereira/build-your-own-linktree-with-go-and-github-pages-3fha", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(60000),
	}); err != nil {
		d.logger.V(1).Error(err, "Failed to navigate to page")
		return nil, fmt.Errorf("could not goto: %v", err)
	}
	d.logger.V(1).Info("Page loaded")

	d.logger.V(1).Info("Getting title")
	title, err := page.Locator("h1").TextContent()
	if err != nil {
		d.logger.V(1).Error(err, "Failed to get title")
		return nil, fmt.Errorf("could not get title: %v", err)
	}
	d.logger.V(1).Info("Title retrieved", "title", title)

	d.logger.V(1).Info("Getting original content")
	bodyLocator := page.Locator("body")
	originalContent, err := bodyLocator.InnerHTML()
	if err != nil {
		d.logger.V(1).Error(err, "Failed to get original content")
		return nil, fmt.Errorf("could not get original content: %v", err)
	}
	d.logger.V(1).Info("Original content retrieved", "length", len(originalContent))

	err = os.WriteFile("original_content.html", []byte(originalContent), 0o644)
	if err != nil {
		d.logger.V(1).Error(err, "Failed to write original content to file")
		return nil, fmt.Errorf("could not write original content: %v", err)
	}
	d.logger.V(1).Info("Original content written to file")

	d.logger.V(1).Info("Updating links")
	baseURL := "https://dev.to"
	updateScript := `
   	const makeAbsolute = (baseURL, relativePath) => {
   		if (!relativePath || relativePath.startsWith('http://') || relativePath.startsWith('https://')) {
   			return relativePath;
   		}
   		return new URL(relativePath, baseURL).href;
   	};

   	const links = document.querySelectorAll('a');
   	links.forEach(link => {
   		link.href = makeAbsolute('` + baseURL + `', link.getAttribute('href'));
   	});

   	document.body.innerHTML;
   `

	updatedContentResult, err := page.Evaluate(updateScript)
	if err != nil {
		d.logger.V(1).Error(err, "Failed to evaluate update script")
		return nil, fmt.Errorf("could not evaluate update script: %v", err)
	}
	d.logger.V(1).Info("Links updated")

	updatedContent, ok := updatedContentResult.(string)
	if !ok {
		d.logger.V(1).Error(nil, "Unexpected updated result type")
		return nil, fmt.Errorf("unexpected updated result type")
	}
	d.logger.V(1).Info("Updated content retrieved", "length", len(updatedContent))

	err = os.WriteFile("updated_content.html", []byte(updatedContent), 0o644)
	if err != nil {
		d.logger.V(1).Error(err, "Failed to write updated content to file")
		return nil, fmt.Errorf("could not write updated content: %v", err)
	}
	d.logger.V(1).Info("Updated content written to file")

	d.logger.V(1).Info("Splitting content into paragraphs")
	paragraphs := strings.Split(updatedContent, "</p>")
	content := []string{title}
	for _, p := range paragraphs {
		if strings.TrimSpace(p) != "" {
			content = append(content, strings.TrimSpace(p)+"</p>")
		}
	}
	d.logger.V(1).Info("Content split into paragraphs", "paragraphCount", len(content))

	d.logger.V(1).Info("GetContent function completed successfully")
	return content, nil
}
