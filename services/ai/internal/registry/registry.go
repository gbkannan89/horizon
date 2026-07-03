package registry

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/horizon/core/services/ai/provider"
	"github.com/horizon/core/services/ai/provider/stubs"
)

type Registry struct {
	mu          sync.RWMutex
	providers   map[string]*registryEntry
	active      string
	cbThreshold int
	cbTimeout   time.Duration
}

type registryEntry struct {
	provider provider.AIProvider
	cb       *CircuitBreaker
}

type CircuitBreaker struct {
	mu           sync.Mutex
	failures     int
	threshold    int
	timeout      time.Duration
	openedAt     time.Time
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.openedAt = time.Time{}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.failures >= cb.threshold && cb.openedAt.IsZero() {
		cb.openedAt = time.Now()
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.failures < cb.threshold {
		return true
	}
	if !cb.openedAt.IsZero() && time.Since(cb.openedAt) > cb.timeout {
		// Half-open: allow one request to test the waters
		// We don't reset failures here; if it fails again, it immediately re-opens.
		// We reset openedAt so that if it fails, it sets a new openedAt time.
		cb.openedAt = time.Time{}
		return true
	}
	return false
}

func New(cbThreshold int, cbTimeout time.Duration) *Registry {
	if cbThreshold <= 0 { cbThreshold = 3 }
	if cbTimeout <= 0 { cbTimeout = 60 * time.Second }
	r := &Registry{
		providers:   make(map[string]*registryEntry),
		active:      "mock",
		cbThreshold: cbThreshold,
		cbTimeout:   cbTimeout,
	}
	r.Register("mock", NewMockWithName("mock"))
	r.Register("ollama", stubs.NewOllamaStub())
	r.Register("openai", stubs.NewOpenAIStub())
	r.Register("azure", stubs.NewAzureStub())
	r.Register("anthropic", stubs.NewAnthropicStub())
	return r
}

func NewMockWithName(name string) provider.AIProvider {
	return &mockNamed{name: name}
}

type mockNamed struct{ name string }

func (m *mockNamed) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	return &provider.ChatResponse{
		Reply: fmt.Sprintf("[%s] This is a mock response. When connected to a real provider, I would provide an intelligent response.", m.name),
		Confidence: "Medium", Provider: m.name, SessionID: req.SessionID,
	}, nil
}
func (m *mockNamed) StreamChat(ctx context.Context, req provider.ChatRequest) (<-chan string, error) {
	ch := make(chan string)
	go func() {
		defer close(ch)
		msg := fmt.Sprintf("[%s] This is a mock streaming response.", m.name)
		for _, word := range strings.Split(msg, " ") {
			select {
			case ch <- word + " ":
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch, nil
}
func (m *mockNamed) Explain(ctx context.Context, req provider.ExplainRequest) (*provider.ExplainResponse, error) {
	return &provider.ExplainResponse{
		Explanation: fmt.Sprintf("[%s] Mock explanation for %s. Replace with real provider for actual explanations.", m.name, req.PromptType),
		Confidence: "Medium", Provider: m.name, Version: req.Version,
	}, nil
}
func (m *mockNamed) Summarize(ctx context.Context, req provider.SummarizeRequest) (*provider.SummarizeResponse, error) {
	return &provider.SummarizeResponse{
		Summary: fmt.Sprintf("[%s] Mock summary for '%s'. Replace with real provider.", m.name, req.Topic),
		Confidence: "Medium", Provider: m.name,
	}, nil
}
func (m *mockNamed) GenerateInsights(ctx context.Context, req provider.GenerateInsightsRequest) (*provider.GenerateInsightsResponse, error) {
	return &provider.GenerateInsightsResponse{
		Insights: []provider.GeneratedInsight{
			{
				Type:        "mock_insight",
				Title:       fmt.Sprintf("[%s] Mock Insight", m.name),
				Summary:     "This is a mock insight summary.",
				Explanation: "Replace with real provider.",
				Confidence:  "Medium",
			},
		},
		Provider: m.name,
	}, nil
}
func (m *mockNamed) Health(ctx context.Context) (*provider.HealthResponse, error) {
	return &provider.HealthResponse{Status: "operational", Provider: m.name, Message: "Mock provider — ready."}, nil
}
func (m *mockNamed) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		Provider: m.name, Chat: true, Explanation: true, Summarize: true,
		Streaming: true, ToolCalling: false, Embeddings: false, Vision: false,
		MaxContext: 4096, Models: []string{"mock-1"},
	}
}

func (r *Registry) Register(name string, p provider.AIProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = &registryEntry{
		provider: p,
		cb:       &CircuitBreaker{threshold: r.cbThreshold, timeout: r.cbTimeout},
	}
}

func (r *Registry) Replace(name string, p provider.AIProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if entry, ok := r.providers[name]; ok {
		entry.provider = p
		// Keep existing circuit breaker state or reset? Resetting is safer.
		entry.cb = &CircuitBreaker{threshold: r.cbThreshold, timeout: r.cbTimeout}
	}
}

func (r *Registry) Get(name string) (provider.AIProvider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if entry, ok := r.providers[name]; ok {
		return entry.provider, true
	}
	return nil, false
}

func (r *Registry) Active() (provider.AIProvider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Try the active provider first
	if entry, ok := r.providers[r.active]; ok && entry.cb.Allow() {
		return entry.provider, true
	}
	
	// Fallback logic
	if r.active != "mock" {
		if mockEntry, ok := r.providers["mock"]; ok {
			return mockEntry.provider, true
		}
	}
	
	// Absolute fallback
	if mockEntry, ok := r.providers["mock"]; ok {
		return mockEntry.provider, true
	}
	return nil, false
}

func (r *Registry) RecordSuccess(name string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if entry, ok := r.providers[name]; ok {
		entry.cb.RecordSuccess()
	}
}

func (r *Registry) RecordFailure(name string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if entry, ok := r.providers[name]; ok {
		entry.cb.RecordFailure()
	}
}

func (r *Registry) ActiveName() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.active
}

func (r *Registry) Switch(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.providers[name]; !ok {
		return fmt.Errorf("provider %q not found", name)
	}
	r.active = name
	return nil
}

func (r *Registry) List() []provider.ProviderInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var infos []provider.ProviderInfo
	for name, entry := range r.providers {
		p := entry.provider
		health, err := p.Health(nil)
		info := provider.ProviderInfo{
			Name:         name,
			Capabilities: p.Capabilities(),
			Healthy:      err == nil,
		}
		if health != nil { info.Status = health.Status; info.Error = health.Message }
		if err != nil { info.Status = "error"; info.Error = err.Error() }
		infos = append(infos, info)
	}
	return infos
}
