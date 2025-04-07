package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/ncecere/reader-go/internal/api/handlers"
	"github.com/ncecere/reader-go/internal/api/middleware"
	"github.com/ncecere/reader-go/internal/core/ai"
	"github.com/ncecere/reader-go/internal/core/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// Server represents the HTTP server
type Server struct {
	app    *fiber.App
	config *Config
}

// Config holds server configuration
type Config struct {
	Port int
}

// New creates a new server instance
func New(config *Config, browserService *service.Service, aiService *ai.Service) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Add middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.NewMetricsMiddleware())

	// Create handlers
	readerHandler := handlers.NewReaderHandler(browserService)
	summaryHandler := handlers.NewSummaryHandler(browserService, aiService)

	// Setup routes
	app.Get("/metrics", MetricsHandler())
	app.Get("/summary/*", summaryHandler.HandleRequest)
	app.Get("/*", readerHandler.HandleRequest)

	return &Server{
		app:    app,
		config: config,
	}
}

// Start starts the server
func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf(":%d", s.config.Port))
}

// MetricsHandler returns a handler for Prometheus metrics
func MetricsHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		handler := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
		handler(c.Context())
		return nil
	}
}
