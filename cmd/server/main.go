package main

import (
	"context"
	"log"
	"net/http"

	"github.com/horizon/core/internal/app/middleware"
	"github.com/horizon/core/internal/bootstrap"

	apiGatewayReg "github.com/horizon/core/services/infra/api-gateway/register"
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
	dataReg "github.com/horizon/core/services/experiences/data/register"
	txReg "github.com/horizon/core/services/domains/financial-event/register"
	aiReg "github.com/horizon/core/services/ai/register"
	hhReg "github.com/horizon/core/services/domains/household/register"
	budgetReg "github.com/horizon/core/services/domains/budget/register"
	acctDomainReg "github.com/horizon/core/services/domains/account/register"
	recurringReg "github.com/horizon/core/services/domains/recurring/register"
	rulesReg "github.com/horizon/core/services/domains/rules/register"
	streamReg "github.com/horizon/core/services/infra/events/register"
)

func main() {
	cfg := bootstrap.LoadConfig()

	pool, err := bootstrap.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	// Create shared JWT service
	secret, accessTTL, refreshTTL := authReg.DefaultJWTConfig()
	authSvc := authReg.NewJWTSvc(secret, accessTTL, refreshTTL)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", bootstrap.HealthHandler("live"))
	mux.HandleFunc("GET /health/ready", bootstrap.HealthHandler("ready"))

	// Auth routes with shared JWT service
	authReg.RegisterRoutesWithSvc(mux, authSvc)

	// Initialize Real-time Events Hub
	hub := streamReg.InitHub(context.Background())
	streamReg.RegisterRoutes(mux, hub)

	// Register all experience routes
	// AI service (self-contained, no DB needed)
	aiProvider := aiReg.RegisterRoutes(mux)

	dashReg.RegisterRoutes(mux, pool)
	acctReg.RegisterRoutes(mux, pool)
	pfReg.RegisterRoutes(mux, pool)
	goalsReg.RegisterRoutes(mux, pool)
	planReg.RegisterRoutes(mux, pool)
	insReg.RegisterRoutes(mux, pool, aiProvider)
	tlReg.RegisterRoutes(mux, pool)
	advReg.RegisterRoutes(mux, pool)
	notifReg.RegisterRoutes(mux, pool, hub)

	// User profile, preferences, privacy
	userReg.RegisterRoutes(mux, pool)

	// Data management (export, import, backup, restore)
	dataReg.RegisterRoutes(mux, pool)

	// Transactions (Financial Event domain)
	txReg.RegisterRoutes(mux, pool)

	// Household domain
	hhReg.RegisterRoutes(mux, pool)

	// Budget domain
	budgetReg.RegisterRoutes(mux, pool)

	// Recurring Transactions domain
	recurringReg.RegisterRoutes(mux, pool)

	// Rules domain — automation rules
	rulesReg.RegisterRoutes(mux, pool)

	// Account domain — household account listings
	acctDomainReg.RegisterRoutes(mux, pool)

	// API Gateway bridge — registers domain CRUD endpoints
	apiGatewayReg.RegisterGatewayRoutes(mux, pool)

	// Wrap with auth, then CORS, then recovery
	handler := middleware.Authenticate(authSvc)(mux)
	handler = authReg.CORS(handler)
	handler = bootstrap.RecoveryMiddleware(handler)

	srv := bootstrap.NewServer(cfg.Addr(), handler)
	srv.StartAndWait()
}
