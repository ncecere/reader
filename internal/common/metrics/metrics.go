package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ContentProcessingDuration tracks content processing time by type
	ContentProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "reader_content_processing_duration_seconds",
			Help:    "Content processing time by type",
			Buckets: []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10, 20, 30},
		},
		[]string{"type"},
	)

	// ContentSize tracks processed content size by type
	ContentSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "reader_content_size_bytes",
			Help:    "Processed content size by type",
			Buckets: []float64{1000, 10000, 100000, 1e6, 1e7, 1e8, 1e9, 1e10},
		},
		[]string{"type"},
	)

	// ContentProcessingErrors tracks processing errors by type
	ContentProcessingErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reader_content_processing_errors_total",
			Help: "Processing errors by type",
		},
		[]string{"type", "error"},
	)

	// URLProcessing tracks URLs processed by domain
	URLProcessing = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reader_url_processing_total",
			Help: "URLs processed by domain",
		},
		[]string{"domain"},
	)

	// URLContentTypes tracks content types encountered
	URLContentTypes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reader_url_content_types_total",
			Help: "Content types encountered",
		},
		[]string{"content_type"},
	)

	// URLSizes tracks URL content size distribution
	URLSizes = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "reader_url_sizes_bytes",
			Help:    "URL content size distribution",
			Buckets: []float64{1000, 10000, 100000, 1e6, 1e7, 1e8, 1e9, 1e10},
		},
		[]string{"domain"},
	)
)
