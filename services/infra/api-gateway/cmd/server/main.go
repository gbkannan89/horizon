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
	"github.com/horizon/core/services/infra/api-gateway/internal/config"
	"github.com/horizon/core/services/infra/api-gateway/internal/handler"
	"github.com/horizon/core/services/infra/api-gateway/internal/middleware"
	"github.com/horizon/core/services/infra/api-gateway/internal/router"
)

func main() {
	cfg := config.Load()
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(
		middleware.CORS(cfg.AllowedOrigins),
		middleware.CorrelationID(),
		middleware.RequestLog(),
		middleware.Recovery(),
	)

	r.GET("/health/live", func(c *gin.Context) { c.JSON(200, gin.H{"status": "live", "timestamp": time.Now().UTC()}) })
	r.GET("/health/ready", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ready"}) })
	r.GET("/health/depth", func(c *gin.Context) { c.JSON(200, gin.H{"status": "healthy", "checks": gin.H{"api": "live"}}) })

	h := handler.NewHandler()
	v1 := r.Group("/api/v1")
	router.RegisterRoutes(v1, h)

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}
	go func() {
		log.Printf("API Gateway starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Print("Server stopped")
}
