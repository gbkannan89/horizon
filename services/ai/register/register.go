package register

import (
	stdctx "context"
	"log"
	"net/http"

	"github.com/horizon/core/services/ai/internal/api"
	"github.com/horizon/core/services/ai/internal/config"
	ctxpkg "github.com/horizon/core/services/ai/internal/context"
	"github.com/horizon/core/services/ai/internal/prompts"
	"github.com/horizon/core/services/ai/internal/provider/ollama"
	"github.com/horizon/core/services/ai/internal/registry"
	"github.com/horizon/core/services/ai/internal/runtime"
	"github.com/horizon/core/services/ai/internal/session"
)

// RegisterRoutes wires the AI service into the shared monolith mux.
// The AI service is self-contained (no PostgreSQL dependency) — it uses
// in-memory state and optional HTTP calls to an Ollama instance.
func RegisterRoutes(mux *http.ServeMux) {
	cfg := config.Load()
	reg := registry.New()

	ollamaProvider := ollama.New(ollama.Config{
		Endpoint: cfg.OllamaURL,
		Timeout:  cfg.Timeout,
		Retries:  cfg.Retries,
		Model:    cfg.OllamaModel,
		Enabled:  cfg.AIEnabled && cfg.Provider == "ollama",
	})
	reg.Replace("ollama", ollamaProvider)

	rtStatus := runtime.NewStatus(cfg)

	healthChecker := runtime.NewHealthChecker(rtStatus,
		func(ctx stdctx.Context) bool {
			if !ollamaProvider.Enabled() { return false }
			if _, err := ollamaProvider.Health(ctx); err != nil { return false }
			return true
		},
		30,
	)
	healthChecker.Start(stdctx.Background())

	if cfg.Provider != "mock" {
		if err := reg.Switch(cfg.Provider); err != nil {
			log.Printf("AI: configured provider %q not available, using mock: %v", cfg.Provider, err)
		}
	}
	rtStatus.SetProvider(reg.ActiveName())

	promptManager := prompts.NewManager()
	sessionManager := session.NewManager(cfg.MaxHistory)
	ctxBuilder := ctxpkg.NewBuilder()
	handlers := api.New(reg, promptManager, sessionManager, ctxBuilder, cfg, rtStatus, ollamaProvider)

	// Register all 13 AI endpoints on the shared mux
	handlers.Register(mux)

	log.Printf("AI service registered (provider: %s, AI enabled: %v)", reg.ActiveName(), cfg.AIEnabled)
}
