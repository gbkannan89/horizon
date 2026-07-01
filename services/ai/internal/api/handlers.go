package api

import (
	stdctx "context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/horizon/core/services/internal/auth"
	"time"

	"github.com/horizon/core/services/ai/internal/config"
	ctxpkg "github.com/horizon/core/services/ai/internal/context"
	"github.com/horizon/core/services/ai/internal/prompts"
	"github.com/horizon/core/services/ai/internal/provider"
	"github.com/horizon/core/services/ai/internal/provider/ollama"
	"github.com/horizon/core/services/ai/internal/registry"
	"github.com/horizon/core/services/ai/internal/runtime"
	"github.com/horizon/core/services/ai/internal/session"
)

type Handlers struct {
	registry    *registry.Registry
	prompts     *prompts.Manager
	sessions    *session.Manager
	ctxBuilder  *ctxpkg.Builder
	cfg         *config.Config
	runtimeSt   *runtime.Status
	ollamaProv  *ollama.OllamaProvider
}

func New(reg *registry.Registry, pm *prompts.Manager, sm *session.Manager, cb *ctxpkg.Builder, cfg *config.Config, rs *runtime.Status, op *ollama.OllamaProvider) *Handlers {
	return &Handlers{registry: reg, prompts: pm, sessions: sm, ctxBuilder: cb, cfg: cfg, runtimeSt: rs, ollamaProv: op}
}

func (h *Handlers) Register(mux *http.ServeMux) {
	// AI core
	mux.HandleFunc("GET /api/v1/ai/health", h.GetHealth)
	mux.HandleFunc("POST /api/v1/ai/chat", h.PostChat)
	mux.HandleFunc("POST /api/v1/ai/explain", h.PostExplain)
	mux.HandleFunc("POST /api/v1/ai/summarize", h.PostSummarize)
	mux.HandleFunc("GET /api/v1/ai/prompts", h.GetPrompts)
	mux.HandleFunc("GET /api/v1/ai/context", h.GetContext)
	// Provider management
	mux.HandleFunc("GET /api/v1/ai/providers", h.GetProviders)
	mux.HandleFunc("GET /api/v1/ai/providers/active", h.GetActiveProvider)
	mux.HandleFunc("GET /api/v1/ai/providers/capabilities", h.GetCapabilities)
	mux.HandleFunc("POST /api/v1/ai/provider/switch", h.PostSwitchProvider)
	// Runtime
	mux.HandleFunc("GET /api/v1/ai/runtime/health", h.GetRuntimeHealth)
	mux.HandleFunc("GET /api/v1/ai/runtime/config", h.GetRuntimeConfig)
	mux.HandleFunc("GET /api/v1/ai/runtime/status", h.GetRuntimeStatus)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]interface{}{
		"success": false, "error": map[string]string{"code": code, "message": message},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func getDefaultUserID(r *http.Request) string {
	return auth.UserIDFromRequest(r)
}

// AI core endpoints

func (h *Handlers) GetHealth(w http.ResponseWriter, r *http.Request) {
	ai := h.registry.Active()
	resp, err := ai.Health(r.Context())
	if err != nil {
		writeError(w, http.StatusServiceUnavailable, "AI_UNAVAILABLE", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]string{"status": resp.Status, "provider": resp.Provider, "message": resp.Message},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) PostChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SessionID string                 `json:"session_id"`
		Message   string                 `json:"message"`
		Context   map[string]interface{} `json:"context,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "invalid chat request"); return
	}
	if req.Message == "" { writeError(w, http.StatusBadRequest, "MISSING_MESSAGE", "message is required"); return }
	if req.SessionID == "" { req.SessionID = "session-" + time.Now().Format("20060102-150405") }

	userID := getDefaultUserID(r)
	conv := h.sessions.GetOrCreate(req.SessionID, userID)
	if req.Context != nil { conv.Context = req.Context }

	h.sessions.AddMessage(req.SessionID, session.Message{Role: "user", Content: req.Message, Timestamp: time.Now().UTC().Format(time.RFC3339)})

	ai := h.registry.Active()
	chatReq := provider.ChatRequest{SessionID: req.SessionID, Message: req.Message, Context: conv.Context}
	for _, m := range h.sessions.GetHistory(req.SessionID) {
		chatReq.History = append(chatReq.History, provider.ChatMessage{Role: m.Role, Content: m.Content})
	}

	resp, err := ai.Chat(r.Context(), chatReq)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true, "data": map[string]string{
				"reply": "[System] The AI provider is currently unavailable. Please try again later or contact support.",
				"confidence": "Low", "provider": "fallback", "session_id": req.SessionID,
			},
			"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
		})
		return
	}

	h.sessions.AddMessage(req.SessionID, session.Message{Role: "assistant", Content: resp.Reply, Timestamp: time.Now().UTC().Format(time.RFC3339), Confidence: resp.Confidence})
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": resp,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) PostExplain(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PromptType string                 `json:"prompt_type"`
		Data       map[string]interface{} `json:"data"`
		Version    int                    `json:"version,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "invalid explain request"); return
	}
	if req.PromptType == "" { writeError(w, http.StatusBadRequest, "MISSING_PROMPT_TYPE", "prompt_type is required"); return }
	if req.Version <= 0 { req.Version = h.prompts.LatestVersion(req.PromptType) }

	tmpl, ok := h.prompts.Get(req.PromptType, req.Version)
	if !ok { writeError(w, http.StatusNotFound, "PROMPT_NOT_FOUND", "no template for prompt_type/version"); return }

	rendered := prompts.Render(tmpl, req.Data)
	log.Printf("Rendered prompt [%s v%d]: %s", req.PromptType, req.Version, rendered)

	ai := h.registry.Active()
	resp, err := ai.Explain(r.Context(), provider.ExplainRequest{PromptType: req.PromptType, Data: req.Data, Version: req.Version})
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"success": true, "data": map[string]string{
				"explanation": "[Fallback] Explanation unavailable. The AI provider could not process this request.",
				"confidence": "Low", "provider": "fallback", "version": fmt.Sprintf("%d", req.Version),
			},
			"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": resp,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) PostSummarize(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Topic string                 `json:"topic"`
		Data  map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "invalid summarize request"); return
	}
	if req.Topic == "" { writeError(w, http.StatusBadRequest, "MISSING_TOPIC", "topic is required"); return }

	ai := h.registry.Active()
	resp, err := ai.Summarize(r.Context(), provider.SummarizeRequest{Topic: req.Topic, Data: req.Data})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "AI_ERROR", err.Error()); return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": resp, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetPrompts(w http.ResponseWriter, r *http.Request) {
	summaries := h.prompts.List()
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": summaries, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetContext(w http.ResponseWriter, r *http.Request) {
	inputs := h.buildTestInputs(getDefaultUserID(r))
	ctx := h.ctxBuilder.Build(inputs)
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": ctx, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)}})
}

// Provider management endpoints

func (h *Handlers) GetProviders(w http.ResponseWriter, r *http.Request) {
	infos := h.registry.List()
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": infos, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) GetActiveProvider(w http.ResponseWriter, r *http.Request) {
	name := h.registry.ActiveName()
	ai := h.registry.Active()
	caps := ai.Capabilities()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]interface{}{"name": name, "capabilities": caps},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetCapabilities(w http.ResponseWriter, r *http.Request) {
	ai := h.registry.Active()
	caps := ai.Capabilities()
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": caps, "metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)}})
}

func (h *Handlers) PostSwitchProvider(w http.ResponseWriter, r *http.Request) {
	var req struct{ Provider string `json:"provider"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "invalid switch request"); return
	}
	if req.Provider == "" { writeError(w, http.StatusBadRequest, "MISSING_PROVIDER", "provider is required"); return }
	if err := h.registry.Switch(req.Provider); err != nil {
		writeError(w, http.StatusNotFound, "PROVIDER_NOT_FOUND", err.Error()); return
	}
	h.runtimeSt.SetProvider(h.registry.ActiveName())
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]string{"active_provider": h.registry.ActiveName()},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

// Runtime endpoints

func (h *Handlers) GetRuntimeHealth(w http.ResponseWriter, r *http.Request) {
	var ollamaAvailable bool
	if h.ollamaProv != nil && h.ollamaProv.Enabled() {
		ctx, cancel := stdctx.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		if _, err := h.ollamaProv.Health(ctx); err == nil {
			ollamaAvailable = true
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]interface{}{
			"ai_enabled":       h.cfg.AIEnabled,
			"active_provider":  h.registry.ActiveName(),
			"ollama_available": ollamaAvailable,
			"status":           "operational",
		},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetRuntimeConfig(w http.ResponseWriter, r *http.Request) {
	rc := runtime.GetRuntimeConfig(h.cfg)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": rc,
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) GetRuntimeStatus(w http.ResponseWriter, r *http.Request) {
	st := h.runtimeSt.Get()
	var ollamaStatus string
	if h.ollamaProv != nil && h.ollamaProv.Enabled() {
		ctx, cancel := stdctx.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		if _, err := h.ollamaProv.Health(ctx); err == nil {
			ollamaStatus = "connected"
		} else {
			ollamaStatus = "unavailable"
		}
	} else {
		ollamaStatus = "disabled"
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true, "data": map[string]interface{}{
			"provider":        st.Provider,
			"ollama":          ollamaStatus,
			"ai_enabled":      h.cfg.AIEnabled,
			"last_checked":    st.OllamaChecked,
			"check_interval":  st.CheckInterval,
		},
		"metadata": map[string]string{"timestamp": time.Now().UTC().Format(time.RFC3339)},
	})
}

func (h *Handlers) buildTestInputs(userID string) ctxpkg.Inputs {
	return ctxpkg.Inputs{
		UserID: userID, NetWorth: 5000000, CashBalance: 500000,
		MonthlyIncome: 150000, MonthlyExp: 100000,
		TotalGoals: 5, OnTrack: 3, AtRisk: 1, FundingGap: 500000,
		TotalAccounts: 8, TotalBalance: 1500000,
		PortfolioValue: 3500000, TotalReturn: 75000, ReturnPct: 2.3,
		RiskScore: 30, HealthScore: 75, HealthGrade: "Good", HealthChg: 3, RiskLevel: "Low",
		NWProjected: 8500000, ProjOnTrack: true, ProjConf: "High",
		HasRecs: true, RecCount: 3, HasOpts: true, OptCount: 2,
		HasSims: true, SimCount: 1, EventCount: 15,
		NotifUnread: 3, InsightTotal: 8,
	}
}
