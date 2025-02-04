package browser

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/ncecere/reader-go/internal/common/logger"
	"github.com/ncecere/reader-go/internal/core/metrics"
	"go.uber.org/zap"
)

// Pool manages a pool of browser instances
type Pool struct {
	instances []*Instance
	opts      *BrowserOptions
	mu        sync.RWMutex // Changed to RWMutex for better concurrency
	ready     chan *Instance
	metrics   *metrics.Metrics
}

// Instance represents a browser instance
type Instance struct {
	ctx    context.Context
	cancel context.CancelFunc
	inUse  bool
}

// NewPool creates a new browser pool
func NewPool(opts *BrowserOptions, metrics *metrics.Metrics) (*Pool, error) {
	pool := &Pool{
		opts:      opts,
		instances: make([]*Instance, opts.PoolSize),
		ready:     make(chan *Instance, opts.PoolSize),
		metrics:   metrics,
	}

	logger.Log.Info("Initializing Chrome instances")

	// Initialize instances
	for i := 0; i < opts.PoolSize; i++ {
		instance, err := pool.createInstance()
		if err != nil {
			pool.Close()
			return nil, fmt.Errorf("failed to create instance %d: %w", i, err)
		}
		pool.instances[i] = instance
		pool.ready <- instance
	}

	// Warm up instances
	if err := pool.Warmup(context.Background()); err != nil {
		logger.Log.Warn("Failed to warm up browser instances", zap.Error(err))
	}

	logger.Log.Info("Browser pool initialized successfully",
		zap.Int("pool_size", opts.PoolSize))

	return pool, nil
}

// Execute runs a function with a browser instance
func (p *Pool) Execute(ctx context.Context, fn func(context.Context) error) error {
	instance, err := p.getInstance()
	if err != nil {
		return fmt.Errorf("failed to get browser instance: %w", err)
	}

	// Create a new context with timeout from the pool options
	timeoutCtx, cancel := context.WithTimeout(instance.ctx, time.Duration(p.opts.Timeout)*time.Second)
	defer cancel()

	// Create a merged context that will be canceled if either the parent context
	// or the timeout context is canceled
	mergedCtx, mergedCancel := context.WithCancel(timeoutCtx)
	defer mergedCancel()

	// Monitor parent context for cancellation
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			mergedCancel()
		case <-timeoutCtx.Done():
			mergedCancel()
		case <-done:
			return
		}
	}()

	// Execute function with merged context
	err = fn(mergedCtx)
	close(done)

	// Mark instance as available
	p.mu.Lock()
	defer p.mu.Unlock()
	instance.inUse = false

	// Handle context cancellation
	if err != nil {
		if err == context.Canceled || err == context.DeadlineExceeded {
			logger.Log.Info("Browser instance timed out, recreating...")
			newInstance, createErr := p.createInstance()
			if createErr == nil {
				for i, inst := range p.instances {
					if inst == instance {
						instance.cancel() // Cancel old instance
						p.instances[i] = newInstance
						break
					}
				}
			}
		}
		return err
	}

	return nil
}

// updateMetrics updates pool and memory metrics
func (p *Pool) updateMetrics() {
	activeCount := int32(0)
	for _, instance := range p.instances {
		if instance.inUse {
			activeCount++
		}
	}
	p.metrics.UpdatePoolMetrics(int32(p.opts.PoolSize), activeCount)

	// Estimate memory usage (rough estimate: 100MB base + MaxMemoryMB per active instance)
	currentMemoryMB := uint64(100 + (p.opts.MaxMemoryMB * int(activeCount)))
	p.metrics.UpdateMemoryUsage(currentMemoryMB)
}

// getInstance gets an available browser instance
func (p *Pool) getInstance() (*Instance, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Update metrics
	p.updateMetrics()

	// Find available instance
	for _, instance := range p.instances {
		if !instance.inUse {
			instance.inUse = true
			p.updateMetrics() // Update metrics after instance allocation
			return instance, nil
		}
	}

	// No available instances, create new one if possible
	if len(p.instances) < p.opts.PoolSize*2 {
		instance, err := p.createInstance()
		if err != nil {
			return nil, err
		}
		instance.inUse = true
		p.instances = append(p.instances, instance)
		p.updateMetrics() // Update metrics after instance creation
		return instance, nil
	}

	return nil, fmt.Errorf("no available browser instances")
}

// Warmup pre-initializes browser instances by navigating to about:blank
func (p *Pool) Warmup(ctx context.Context) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, instance := range p.instances {
		if err := chromedp.Run(instance.ctx, chromedp.Navigate("about:blank")); err != nil {
			return fmt.Errorf("warmup failed: %w", err)
		}
	}
	return nil
}

// createInstance creates a new browser instance
func (p *Pool) createInstance() (*Instance, error) {
	// Get optimized Chrome flags
	opts := GetChromeFlags(p.opts)
	opts = append(opts, chromedp.WindowSize(1920, 1080))

	// Create context with custom allocator options
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// Create browser context with selective logging
	ctx, _ := chromedp.NewContext(allocCtx,
		chromedp.WithLogf(func(format string, args ...interface{}) {
			// Only log non-cookie related messages at debug level
			if !strings.Contains(format, "cookiePart") {
				logger.Log.Debug(fmt.Sprintf(format, args...))
			}
		}),
		chromedp.WithErrorf(func(format string, args ...interface{}) {
			// Filter out known non-critical errors
			msg := fmt.Sprintf(format, args...)
			if !strings.Contains(msg, "cookiePart") &&
				!strings.Contains(msg, "unknown ClientNavigationReason") {
				logger.Log.Error(msg)
			}
		}),
	)

	// Create base context with long timeout
	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)

	// Start the browser
	if err := chromedp.Run(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start browser: %w", err)
	}

	return &Instance{
		ctx:    ctx,
		cancel: cancel,
		inUse:  false,
	}, nil
}

// Close closes all browser instances
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, instance := range p.instances {
		if instance.cancel != nil {
			instance.cancel()
		}
	}
}
