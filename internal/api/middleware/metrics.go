package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "reader_http_request_duration_seconds",
			Help:    "Request latency distribution",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"endpoint"},
	)

	requestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "reader_http_request_size_bytes",
			Help:    "Request size distribution",
			Buckets: []float64{100, 1000, 10000, 100000, 1e6, 1e7, 1e8, 1e9},
		},
		[]string{"endpoint"},
	)

	responseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "reader_http_response_size_bytes",
			Help:    "Response size distribution",
			Buckets: []float64{100, 1000, 10000, 100000, 1e6, 1e7, 1e8, 1e9},
		},
		[]string{"endpoint"},
	)

	inFlightRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "reader_http_in_flight_requests",
			Help: "Current number of in-flight requests",
		},
	)

	totalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reader_http_requests_total",
			Help: "Total number of HTTP requests by endpoint and status code",
		},
		[]string{"endpoint", "status"},
	)
)

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		path := "/url" // Normalize all URL paths
		inFlightRequests.Inc()

		// Record request size
		requestSize.WithLabelValues(path).Observe(float64(len(c.Request().Body())))

		// Process request
		err := c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())

		requestDuration.WithLabelValues(path).Observe(duration)
		responseSize.WithLabelValues(path).Observe(float64(len(c.Response().Body())))
		totalRequests.WithLabelValues(path, status).Inc()
		inFlightRequests.Dec()

		return err
	}
}
