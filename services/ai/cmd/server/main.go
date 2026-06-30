package main

import (
	stdctx "context"
	"log"

	pkghttp "github.com/horizon/core/packages/http"
	"github.com/horizon/core/services/ai/internal/api"
	"github.com/horizon/core/services/ai/internal/config"
	ctxpkg "github.com/horizon/core/services/ai/internal/context"
	"github.com/horizon/core/services/ai/internal/prompts"
	"github.com/horizon/core/services/ai/internal/provider/ollama"
	"github.com/horizon/core/services/ai/internal/registry"
	"github.com/horizon/core/services/ai/internal/runtime"
	"github.com/horizon/core/services/ai/internal/session"
)

func main() {
	cfg := config.Load()
	reg := registry.New()

	// Build Ollama provider
	ollamaProvider := ollama.New(ollama.Config{
		Endpoint: cfg.OllamaURL,
		Timeout:  cfg.Timeout,
		Retries:  cfg.Retries,
		Model:    cfg.OllamaModel,
		Enabled:  cfg.AIEnabled && cfg.Provider == "ollama",
	})

	// Replace the Ollama stub with the real provider
	reg.Replace("ollama", ollamaProvider)

	// Set up runtime status tracking
	rtStatus := runtime.NewStatus(cfg)

	// Health checker — periodically checks Ollama availability
	healthChecker := runtime.NewHealthChecker(rtStatus,
		func(ctx stdctx.Context) bool {
			if !ollamaProvider.Enabled() { return false }
			if _, err := ollamaProvider.Health(ctx); err != nil { return false }
			return true
		},
		30,
	)
	healthChecker.Start(stdctx.Background())
	defer healthChecker.Stop()

	// Switch to configured provider
	if cfg.Provider != "mock" {
		if err := reg.Switch(cfg.Provider); err != nil {
			log.Printf("Warning: configured provider %q not available, using mock: %v", cfg.Provider, err)
		}
	}
	rtStatus.SetProvider(reg.ActiveName())

	// Initialize components
	promptManager := prompts.NewManager()
	sessionManager := session.NewManager(cfg.MaxHistory)
	ctxBuilder := ctxpkg.NewBuilder()
	handlers := api.New(reg, promptManager, sessionManager, ctxBuilder, cfg, rtStatus, ollamaProvider)
	router := api.NewRouter(handlers)
	srv := pkghttp.New(cfg.Addr(), router)

	log.Printf("Starting AI Service on %s (provider: %s, AI enabled: %v)", cfg.Addr(), reg.ActiveName(), cfg.AIEnabled)
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
