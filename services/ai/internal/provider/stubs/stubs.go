package stubs

import (
	"context"

	"github.com/horizon/core/services/ai/internal/provider"
)

// OllamaStub is a placeholder for the future Ollama provider.
// No network calls, no model downloads, no external dependencies.
type OllamaStub struct{ Name string }

func NewOllamaStub() *OllamaStub { return &OllamaStub{Name: "ollama"} }

func (s *OllamaStub) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	return nil, stubError("Chat")
}
func (s *OllamaStub) Explain(ctx context.Context, req provider.ExplainRequest) (*provider.ExplainResponse, error) {
	return nil, stubError("Explain")
}
func (s *OllamaStub) Summarize(ctx context.Context, req provider.SummarizeRequest) (*provider.SummarizeResponse, error) {
	return nil, stubError("Summarize")
}
func (s *OllamaStub) Health(ctx context.Context) (*provider.HealthResponse, error) {
	return &provider.HealthResponse{Status: "stub", Provider: s.Name, Message: "Ollama provider stub — not yet implemented."}, nil
}
func (s *OllamaStub) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		Provider: s.Name, Chat: true, Explanation: true, Summarize: true,
		Streaming: true, ToolCalling: false, Embeddings: true, Vision: false,
		MaxContext: 8192, Models: []string{"llama3", "mistral", "codellama"},
	}
}

// OpenAIStub is a placeholder for the future OpenAI provider.
type OpenAIStub struct{ Name string }

func NewOpenAIStub() *OpenAIStub { return &OpenAIStub{Name: "openai"} }

func (s *OpenAIStub) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	return nil, stubError("Chat")
}
func (s *OpenAIStub) Explain(ctx context.Context, req provider.ExplainRequest) (*provider.ExplainResponse, error) {
	return nil, stubError("Explain")
}
func (s *OpenAIStub) Summarize(ctx context.Context, req provider.SummarizeRequest) (*provider.SummarizeResponse, error) {
	return nil, stubError("Summarize")
}
func (s *OpenAIStub) Health(ctx context.Context) (*provider.HealthResponse, error) {
	return &provider.HealthResponse{Status: "stub", Provider: s.Name, Message: "OpenAI provider stub — not yet implemented."}, nil
}
func (s *OpenAIStub) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		Provider: s.Name, Chat: true, Explanation: true, Summarize: true,
		Streaming: true, ToolCalling: true, Embeddings: true, Vision: true,
		MaxContext: 128000, Models: []string{"gpt-4o", "gpt-4o-mini"},
	}
}

// AzureStub is a placeholder for the future Azure OpenAI provider.
type AzureStub struct{ Name string }

func NewAzureStub() *AzureStub { return &AzureStub{Name: "azure"} }

func (s *AzureStub) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	return nil, stubError("Chat")
}
func (s *AzureStub) Explain(ctx context.Context, req provider.ExplainRequest) (*provider.ExplainResponse, error) {
	return nil, stubError("Explain")
}
func (s *AzureStub) Summarize(ctx context.Context, req provider.SummarizeRequest) (*provider.SummarizeResponse, error) {
	return nil, stubError("Summarize")
}
func (s *AzureStub) Health(ctx context.Context) (*provider.HealthResponse, error) {
	return &provider.HealthResponse{Status: "stub", Provider: s.Name, Message: "Azure OpenAI provider stub — not yet implemented."}, nil
}
func (s *AzureStub) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		Provider: s.Name, Chat: true, Explanation: true, Summarize: true,
		Streaming: true, ToolCalling: true, Embeddings: true, Vision: true,
		MaxContext: 128000, Models: []string{"gpt-4o", "gpt-4o-mini"},
	}
}

// AnthropicStub is a placeholder for the future Anthropic (Claude) provider.
type AnthropicStub struct{ Name string }

func NewAnthropicStub() *AnthropicStub { return &AnthropicStub{Name: "anthropic"} }

func (s *AnthropicStub) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	return nil, stubError("Chat")
}
func (s *AnthropicStub) Explain(ctx context.Context, req provider.ExplainRequest) (*provider.ExplainResponse, error) {
	return nil, stubError("Explain")
}
func (s *AnthropicStub) Summarize(ctx context.Context, req provider.SummarizeRequest) (*provider.SummarizeResponse, error) {
	return nil, stubError("Summarize")
}
func (s *AnthropicStub) Health(ctx context.Context) (*provider.HealthResponse, error) {
	return &provider.HealthResponse{Status: "stub", Provider: s.Name, Message: "Anthropic provider stub — not yet implemented."}, nil
}
func (s *AnthropicStub) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		Provider: s.Name, Chat: true, Explanation: true, Summarize: true,
		Streaming: true, ToolCalling: true, Embeddings: false, Vision: true,
		MaxContext: 200000, Models: []string{"claude-3-opus", "claude-3-sonnet", "claude-3-haiku"},
	}
}

func stubError(method string) *provider.ProviderError {
	return &provider.ProviderError{
		Code: "NOT_IMPLEMENTED", Message: method + " not implemented in stub provider",
		Retryable: false,
	}
}
