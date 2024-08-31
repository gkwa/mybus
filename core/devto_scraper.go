package core

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/playwright-community/playwright-go"
)

const POPUP_CLOSER_TIMER = 3

type DevToScraper struct {
	logger      logr.Logger
	showLinks   bool
	showBrowser bool
}

func NewDevToScraper(logger logr.Logger, showLinks, showBrowser bool) *DevToScraper {
	return &DevToScraper{logger: logger, showLinks: showLinks, showBrowser: showBrowser}
}

func (d *DevToScraper) GetContent(page playwright.Page) ([]string, error) {
	d.logger.V(1).Info("Starting GetContent function")

	targetURL := "https://dev.to/lucasnevespereira/build-your-own-linktree-with-go-and-github-pages-3fha"

	d.logger.V(1).Info("Navigating to page")
	if _, err := page.Goto(targetURL, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(60000),
	}); err != nil {
		d.logger.V(1).Error(err, "Failed to navigate to page")
		return nil, fmt.Errorf("could not goto: %v", err)
	}
	d.logger.V(1).Info("Page loaded")

	if d.showBrowser {
		d.logger.V(1).Info("Waiting for POPUP_CLOSER_TIMER to close popup")
		time.Sleep(POPUP_CLOSER_TIMER * time.Second)

		d.logger.V(1).Info("Attempting to close popup")
		closePopupScript := `
			const closeButton = document.querySelector('button[aria-label="Close"]');
			if (closeButton) {
				closeButton.click();
				console.log("Popup closed");
			} else {
				console.log("No popup found");
			}
		`
		_, err := page.Evaluate(closePopupScript)
		if err != nil {
			d.logger.V(1).Error(err, "Failed to execute popup close script")
			return nil, fmt.Errorf("could not close popup: %v", err)
		}
		d.logger.V(1).Info("Popup close script executed")
	}

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
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		d.logger.V(1).Error(err, "Failed to parse URL")
		return nil, fmt.Errorf("could not parse URL: %v", err)
	}
	baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	updateScript := fmt.Sprintf(`
		const makeAbsolute = (baseURL, relativePath) => {
			if (!relativePath || relativePath.startsWith('http://') || relativePath.startsWith('https://')) {
				return relativePath;
			}
			return new URL(relativePath, baseURL).href;
		};

		const links = document.querySelectorAll('a');
		const absoluteLinks = [];
		links.forEach(link => {
			const absoluteHref = makeAbsolute('%s', link.getAttribute('href'));
			link.href = absoluteHref;
			absoluteLinks.push(absoluteHref);
		});

		window.updatedContent = document.body.innerHTML;
		window.absoluteLinks = absoluteLinks;
	`, baseURL)

	_, err = page.Evaluate(updateScript)
	if err != nil {
		d.logger.V(1).Error(err, "Failed to evaluate update script")
		return nil, fmt.Errorf("could not evaluate update script: %v", err)
	}
	d.logger.V(1).Info("Links updated")

	updatedContent, err := page.Evaluate(`window.updatedContent`)
	if err != nil {
		d.logger.V(1).Error(err, "Failed to retrieve updated content")
		return nil, fmt.Errorf("could not retrieve updated content: %v", err)
	}

	updatedContentStr, ok := updatedContent.(string)
	if !ok {
		d.logger.V(1).Error(nil, "Unexpected updated content type")
		return nil, fmt.Errorf("unexpected updated content type")
	}
	d.logger.V(1).Info("Updated content retrieved", "length", len(updatedContentStr))

	if d.showLinks {
		links, err := page.Evaluate(`window.absoluteLinks`)
		if err != nil {
			d.logger.V(1).Error(err, "Failed to retrieve links")
			return nil, fmt.Errorf("could not retrieve links: %v", err)
		}

		linksSlice, ok := links.([]interface{})
		if !ok {
			d.logger.V(1).Error(nil, "Unexpected links type")
			return nil, fmt.Errorf("unexpected links type")
		}

		fmt.Println("Absolute links found on the page:")
		for _, link := range linksSlice {
			fmt.Println(link)
		}
	}

	err = os.WriteFile("updated_content.html", []byte(updatedContentStr), 0o644)
	if err != nil {
		d.logger.V(1).Error(err, "Failed to write updated content to file")
		return nil, fmt.Errorf("could not write updated content: %v", err)
	}
	d.logger.V(1).Info("Updated content written to file")

	d.logger.V(1).Info("Splitting content into paragraphs")
	paragraphs := strings.Split(updatedContentStr, "</p>")
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
