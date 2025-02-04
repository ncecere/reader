package converter

import (
	"fmt"
	"regexp"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
)

// HTMLToMarkdown converts HTML content to Markdown format
func HTMLToMarkdown(html string) (string, error) {
	converter := md.NewConverter("", true, nil)

	// Add GitHub Flavored Markdown plugin
	converter.Use(plugin.GitHubFlavored())

	// Add table support
	converter.Use(plugin.Table())

	// Convert to markdown
	markdown, err := converter.ConvertString(html)
	if err != nil {
		return "", err
	}

	// Extract title using regex
	titleRegex := regexp.MustCompile(`<title[^>]*>([^<]+)</title>`)
	matches := titleRegex.FindStringSubmatch(html)
	title := "Untitled"
	if len(matches) > 1 {
		title = matches[1]
	}

	// Add metadata
	metadata := "Title: " + title + "\n" +
		"Visited: " + getCurrentTime() + "\n\n" +
		"Markdown Content:\n"

	return metadata + markdown, nil
}

// getCurrentTime returns the current time in a readable format
func getCurrentTime() string {
	return time.Now().Format("01-02-2006 15:04")
}

// TextToMarkdown converts plain text to a formatted Markdown document
func TextToMarkdown(text string) (string, error) {
	// Add metadata
	metadata := fmt.Sprintf("Summary Generated: %s\n\n", getCurrentTime())

	// Format the text as a markdown document
	// We'll keep the text as-is but add proper spacing and formatting
	content := metadata + text

	return content, nil
}
