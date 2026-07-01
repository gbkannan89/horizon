package main

import (
	"log"
	"net/http"

	"github.com/horizon/core/internal/bootstrap"

	authReg "github.com/horizon/core/services/experiences/auth/register"
	dashReg "github.com/horizon/core/services/experiences/dashboard/register"
	acctReg "github.com/horizon/core/services/experiences/accounts/register"
	pfReg "github.com/horizon/core/services/experiences/portfolio/register"
	goalsReg "github.com/horizon/core/services/experiences/goals/register"
	planReg "github.com/horizon/core/services/experiences/planning/register"
	insReg "github.com/horizon/core/services/experiences/insights/register"
	tlReg "github.com/horizon/core/services/experiences/timeline/register"
	advReg "github.com/horizon/core/services/experiences/advisor/register"
	notifReg "github.com/horizon/core/services/experiences/notifications/register"
	userReg "github.com/horizon/core/services/experiences/user/register"
	txReg "github.com/horizon/core/services/domains/financial-event/register"
	aiReg "github.com/horizon/core/services/ai/register"
)

func main() {
	cfg := bootstrap.LoadConfig()

	pool, err := bootstrap.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", bootstrap.HealthHandler("live"))
	mux.HandleFunc("GET /health/ready", bootstrap.HealthHandler("ready"))

	// Auth routes (no DB needed — JWT-based MVP)
	authReg.RegisterRoutes(mux)

	// Register all experience routes
	dashReg.RegisterRoutes(mux, pool)
	acctReg.RegisterRoutes(mux, pool)
	pfReg.RegisterRoutes(mux, pool)
	goalsReg.RegisterRoutes(mux, pool)
	planReg.RegisterRoutes(mux, pool)
	insReg.RegisterRoutes(mux, pool)
	tlReg.RegisterRoutes(mux, pool)
	advReg.RegisterRoutes(mux, pool)
	notifReg.RegisterRoutes(mux, pool)

	// User profile, preferences, privacy
	userReg.RegisterRoutes(mux, pool)

	// Transactions (Financial Event domain)
	txReg.RegisterRoutes(mux, pool)

	// AI service (self-contained, no DB needed)
	aiReg.RegisterRoutes(mux)

	// Wrap with CORS, then recovery
	handler := authReg.CORS(bootstrap.RecoveryMiddleware(mux))
	srv := bootstrap.NewServer(cfg.Addr(), handler)
	srv.StartAndWait()
}
