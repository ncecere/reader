package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics collects performance and operational metrics
type Metrics struct {
	startTime time.Time

	// Browser metrics
	totalRequests     uint64
	successfulReqs    uint64
	failedReqs        uint64
	totalProcessingMs uint64

	// Cache metrics
	cacheHits   uint64
	cacheMisses uint64

	// Memory metrics
	peakMemoryMB uint64
	currentMemMB uint64

	// Pool metrics
	poolSize        int32
	activeInstances int32

	mu sync.RWMutex
}

// New creates a new metrics collector
func New() *Metrics {
	return &Metrics{
		startTime: time.Now(),
	}
}

// RecordRequest records a request attempt
func (m *Metrics) RecordRequest(duration time.Duration, success bool) {
	atomic.AddUint64(&m.totalRequests, 1)
	if success {
		atomic.AddUint64(&m.successfulReqs, 1)
	} else {
		atomic.AddUint64(&m.failedReqs, 1)
	}
	atomic.AddUint64(&m.totalProcessingMs, uint64(duration.Milliseconds()))
}

// RecordCacheAccess records a cache hit or miss
func (m *Metrics) RecordCacheAccess(hit bool) {
	if hit {
		atomic.AddUint64(&m.cacheHits, 1)
	} else {
		atomic.AddUint64(&m.cacheMisses, 1)
	}
}

// UpdateMemoryUsage updates memory usage metrics
func (m *Metrics) UpdateMemoryUsage(currentMB uint64) {
	atomic.StoreUint64(&m.currentMemMB, currentMB)
	if currentMB > atomic.LoadUint64(&m.peakMemoryMB) {
		atomic.StoreUint64(&m.peakMemoryMB, currentMB)
	}
}

// UpdatePoolMetrics updates browser pool metrics
func (m *Metrics) UpdatePoolMetrics(size, active int32) {
	atomic.StoreInt32(&m.poolSize, size)
	atomic.StoreInt32(&m.activeInstances, active)
}

// Stats represents collected metrics
type Stats struct {
	UptimeSeconds   int64
	TotalRequests   uint64
	SuccessfulReqs  uint64
	FailedReqs      uint64
	AvgProcessingMs float64
	CacheHitRate    float64
	CurrentMemoryMB uint64
	PeakMemoryMB    uint64
	PoolSize        int32
	ActiveInstances int32
	InstanceUtilPct float64
}

// GetStats returns current metrics
func (m *Metrics) GetStats() Stats {
	var stats Stats

	stats.UptimeSeconds = int64(time.Since(m.startTime).Seconds())
	stats.TotalRequests = atomic.LoadUint64(&m.totalRequests)
	stats.SuccessfulReqs = atomic.LoadUint64(&m.successfulReqs)
	stats.FailedReqs = atomic.LoadUint64(&m.failedReqs)

	totalMs := atomic.LoadUint64(&m.totalProcessingMs)
	if stats.TotalRequests > 0 {
		stats.AvgProcessingMs = float64(totalMs) / float64(stats.TotalRequests)
	}

	hits := atomic.LoadUint64(&m.cacheHits)
	misses := atomic.LoadUint64(&m.cacheMisses)
	total := hits + misses
	if total > 0 {
		stats.CacheHitRate = float64(hits) / float64(total)
	}

	stats.CurrentMemoryMB = atomic.LoadUint64(&m.currentMemMB)
	stats.PeakMemoryMB = atomic.LoadUint64(&m.peakMemoryMB)
	stats.PoolSize = atomic.LoadInt32(&m.poolSize)
	stats.ActiveInstances = atomic.LoadInt32(&m.activeInstances)

	if stats.PoolSize > 0 {
		stats.InstanceUtilPct = float64(stats.ActiveInstances) / float64(stats.PoolSize) * 100
	}

	return stats
}

// Reset resets all metrics except uptime
func (m *Metrics) Reset() {
	atomic.StoreUint64(&m.totalRequests, 0)
	atomic.StoreUint64(&m.successfulReqs, 0)
	atomic.StoreUint64(&m.failedReqs, 0)
	atomic.StoreUint64(&m.totalProcessingMs, 0)
	atomic.StoreUint64(&m.cacheHits, 0)
	atomic.StoreUint64(&m.cacheMisses, 0)
	atomic.StoreUint64(&m.currentMemMB, 0)
	atomic.StoreUint64(&m.peakMemoryMB, 0)
	atomic.StoreInt32(&m.poolSize, 0)
	atomic.StoreInt32(&m.activeInstances, 0)
}
