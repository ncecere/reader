package browser

import (
	"fmt"

	"github.com/chromedp/chromedp"
)

// BrowserOptions contains configuration for the browser service
type BrowserOptions struct {
	PoolSize    int
	ChromePath  string
	Timeout     int    // in seconds
	MaxMemoryMB int    // Maximum memory per instance in MB
	UserAgent   string // Custom user agent
}

// DefaultOptions returns default browser options
func DefaultOptions() *BrowserOptions {
	return &BrowserOptions{
		PoolSize:    3,
		Timeout:     30,
		MaxMemoryMB: 128,
		UserAgent:   "Mozilla/5.0 (X11; Linux x86_64) Chrome/90.0.4430.212",
	}
}

// GetChromeFlags returns optimized Chrome flags for better performance
func GetChromeFlags(opts *BrowserOptions) []chromedp.ExecAllocatorOption {
	// Performance and memory optimization flags
	flags := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.NoSandbox,

		// SSL/Security settings
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("allow-insecure-localhost", true),
		chromedp.Flag("allow-running-insecure-content", true),

		// Basic settings
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-notifications", true),
		chromedp.Flag("disable-popup-blocking", true),

		// Memory optimization
		chromedp.Flag("js-flags", fmt.Sprintf("--max-old-space-size=%d", opts.MaxMemoryMB)),

		// Set custom user agent
		chromedp.Flag("user-agent", opts.UserAgent),
	}

	// Add Chrome path if provided
	if opts.ChromePath != "" {
		flags = append(flags, chromedp.ExecPath(opts.ChromePath))
	}

	return flags
}
