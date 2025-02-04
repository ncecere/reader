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
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-notifications", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-component-update", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("ignore-certificate-errors", true),

		// Memory optimization
		chromedp.Flag("js-flags", fmt.Sprintf("--max-old-space-size=%d", opts.MaxMemoryMB)),
		chromedp.Flag("memory-pressure-off", true),

		// Performance optimization
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("disable-gpu-compositing", true),
		chromedp.Flag("disable-gpu-rasterization", true),
		chromedp.Flag("disable-gpu-vsync", true),

		// Set custom user agent
		chromedp.Flag("user-agent", opts.UserAgent),
	}

	// Add Chrome path if provided
	if opts.ChromePath != "" {
		flags = append(flags, chromedp.ExecPath(opts.ChromePath))
	}

	return flags
}
