package registry

import (
	"context"
	"fmt"
	"sync"

	"github.com/horizon/core/services/ai/internal/provider"
	"github.com/horizon/core/services/ai/internal/provider/stubs"
)

type Registry struct {
	mu        sync.RWMutex
	providers map[string]provider.AIProvider
	active    string
}

func New() *Registry {
	r := &Registry{
		providers: make(map[string]provider.AIProvider),
		active:    "mock",
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
func (m *mockNamed) Health(ctx context.Context) (*provider.HealthResponse, error) {
	return &provider.HealthResponse{Status: "operational", Provider: m.name, Message: "Mock provider — ready."}, nil
}
func (m *mockNamed) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		Provider: m.name, Chat: true, Explanation: true, Summarize: true,
		Streaming: false, ToolCalling: false, Embeddings: false, Vision: false,
		MaxContext: 4096, Models: []string{"mock-1"},
	}
}

func (r *Registry) Register(name string, p provider.AIProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = p
}

func (r *Registry) Replace(name string, p provider.AIProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.providers[name]; ok {
		r.providers[name] = p
	}
}

func (r *Registry) Get(name string) (provider.AIProvider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	return p, ok
}

func (r *Registry) Active() provider.AIProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[r.active]
	if !ok { return r.providers["mock"] }
	return p
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
	for name, p := range r.providers {
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
