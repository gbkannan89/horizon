package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/horizon/core/services/infra/api-gateway/internal/auth"
	"github.com/horizon/core/services/infra/api-gateway/internal/config"
	"github.com/horizon/core/services/infra/api-gateway/internal/docs"
	"github.com/horizon/core/services/infra/api-gateway/internal/handler"
	"github.com/horizon/core/services/infra/api-gateway/internal/metrics"
	"github.com/horizon/core/services/infra/api-gateway/internal/middleware"
	"github.com/horizon/core/services/infra/api-gateway/internal/router"
)

func main() {
	cfg := config.Load()
	gin.SetMode(gin.ReleaseMode)

	// Authentication
	authCfg := auth.DefaultConfig()
	if s := os.Getenv("JWT_SECRET"); s != "" {
		authCfg.Secret = s
	}
	authService := auth.NewService(authCfg)
	middleware.InitAuth(authService)

	// Service wiring
	handler.InitServices(&handler.ServiceRegistry{})

	// Gin engine with full middleware chain
	r := gin.New()
	r.Use(
		middleware.CORS(cfg.AllowedOrigins),
		middleware.SecurityHeaders(),
		middleware.RequestSize(5<<20), // 5MB max body
		middleware.CorrelationID(),
		middleware.RequestLog(),
		middleware.SlowRequest(500*time.Millisecond),
		middleware.Recovery(),
		middleware.RateLimiter(1000, 1*time.Minute),
	)

	// Health endpoints (always public, no rate limit)
	r.GET("/health/live", func(c *gin.Context) { c.JSON(200, gin.H{"status": "live"}) })
	r.GET("/health/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })
	r.GET("/health/depth", func(c *gin.Context) { c.JSON(200, gin.H{"status": "healthy"}) })

	// Metrics (public, monitoring system)
	r.GET("/metrics", func(c *gin.Context) { metrics.Handler()(c.Writer, c.Request) })

	// Build info
	r.GET("/build", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version":   "0.1.0",
			"buildTime": time.Now().UTC().Format(time.RFC3339),
			"service":   "horizon-api-gateway",
		})
	})

	h := handler.NewHandler(authService)
	v1 := r.Group("/api/v1")
	docs.RegisterOpenAPI(v1)
	router.RegisterRoutes(v1, h)

	// Server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("API Gateway starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Graceful shutdown with timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Print("Server stopped")
}
