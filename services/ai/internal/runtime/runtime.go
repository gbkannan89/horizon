package runtime

import (
	"context"
	"sync"
	"time"

	"github.com/horizon/core/services/ai/internal/config"
)

// Status tracks the runtime health of the AI system.
type Status struct {
	mu            sync.RWMutex
	provider      string
	ollamaHealthy bool
	ollamaChecked time.Time
	checkInterval time.Duration
	history       []HealthEvent
}

type HealthEvent struct {
	Timestamp string `json:"timestamp"`
	Provider  string `json:"provider"`
	Status    string `json:"status"`
}

func NewStatus(cfg *config.Config) *Status {
	return &Status{
		provider:      cfg.Provider,
		checkInterval: 30 * time.Second,
		history:       make([]HealthEvent, 0),
	}
}

func (s *Status) SetProvider(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.provider = name
}

func (s *Status) SetOllamaHealthy(healthy bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	statusStr := "down"
	if healthy {
		statusStr = "up"
	}

	// Only record if state changed or if history is empty
	if len(s.history) == 0 || s.ollamaHealthy != healthy {
		s.history = append(s.history, HealthEvent{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Provider:  "ollama",
			Status:    statusStr,
		})
		// Keep last 10 events
		if len(s.history) > 10 {
			s.history = s.history[1:]
		}
	}

	s.ollamaHealthy = healthy
	s.ollamaChecked = time.Now()
}

func (s *Status) Get() StatusInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return StatusInfo{
		Provider:        s.provider,
		OllamaHealthy:   s.ollamaHealthy,
		OllamaChecked:   s.ollamaChecked.Format(time.RFC3339),
		CheckInterval:   s.checkInterval.String(),
		History:         s.history,
	}
}

type StatusInfo struct {
	Provider        string        `json:"provider"`
	OllamaHealthy   bool          `json:"ollama_healthy"`
	OllamaChecked   string        `json:"ollama_checked"`
	CheckInterval   string        `json:"check_interval"`
	History         []HealthEvent `json:"history"`
}

// RuntimeConfig exposes the active runtime configuration.
type RuntimeConfig struct {
	Provider    string        `json:"provider"`
	AIEnabled   bool          `json:"ai_enabled"`
	OllamaURL   string        `json:"ollama_url"`
	OllamaModel string        `json:"ollama_model"`
	Timeout     string        `json:"timeout"`
	Retries     int           `json:"retries"`
}

func GetRuntimeConfig(cfg *config.Config) RuntimeConfig {
	return RuntimeConfig{
		Provider: cfg.Provider, AIEnabled: cfg.AIEnabled,
		OllamaURL: cfg.OllamaURL, OllamaModel: cfg.OllamaModel,
		Timeout: cfg.Timeout.String(), Retries: cfg.Retries,
	}
}

// HealthChecker performs periodic health checks.
type HealthChecker struct {
	status   *Status
	checkFn  func(context.Context) bool
	interval time.Duration
	stopCh   chan struct{}
}

func NewHealthChecker(status *Status, checkFn func(context.Context) bool, interval time.Duration) *HealthChecker {
	return &HealthChecker{status: status, checkFn: checkFn, interval: interval, stopCh: make(chan struct{})}
}

func (h *HealthChecker) Start(ctx context.Context) {
	h.checkOnce(ctx)
	go func() {
		ticker := time.NewTicker(h.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				h.checkOnce(ctx)
			case <-h.stopCh:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (h *HealthChecker) Stop() { close(h.stopCh) }

func (h *HealthChecker) checkOnce(ctx context.Context) {
	healthy := h.checkFn(ctx)
	h.status.SetOllamaHealthy(healthy)
}
