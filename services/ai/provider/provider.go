package provider

import "context"

// AIProvider defines the interface for AI providers.
// Implementations must be stateless and swappable via configuration.
type AIProvider interface {
	// Chat handles a conversational interaction.
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
	// StreamChat handles a conversational interaction with a streaming response.
	StreamChat(ctx context.Context, req ChatRequest) (<-chan string, error)
	// Explain generates an explanation for a given prompt type with deterministic data.
	Explain(ctx context.Context, req ExplainRequest) (*ExplainResponse, error)
	// Summarize generates a summary of financial data.
	Summarize(ctx context.Context, req SummarizeRequest) (*SummarizeResponse, error)
	// GenerateInsights generates proactive insights from financial data.
	GenerateInsights(ctx context.Context, req GenerateInsightsRequest) (*GenerateInsightsResponse, error)
	// Health returns provider health status.
	Health(ctx context.Context) (*HealthResponse, error)
	// Capabilities returns the provider's capabilities.
	Capabilities() Capabilities
}

type ChatRequest struct {
	SessionID string                 `json:"session_id"`
	Message   string                 `json:"message"`
	History   []ChatMessage          `json:"history,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Reply      string `json:"reply"`
	Confidence string `json:"confidence"`
	Provider   string `json:"provider"`
	SessionID  string `json:"session_id"`
}

type ExplainRequest struct {
	PromptType string                 `json:"prompt_type"`
	Data       map[string]interface{} `json:"data"`
	Version    int                    `json:"version"`
}

type ExplainResponse struct {
	Explanation string `json:"explanation"`
	Confidence  string `json:"confidence"`
	Provider    string `json:"provider"`
	Version     int    `json:"version"`
}

type SummarizeRequest struct {
	Topic string                 `json:"topic"`
	Data  map[string]interface{} `json:"data"`
}

type SummarizeResponse struct {
	Summary    string `json:"summary"`
	Confidence string `json:"confidence"`
	Provider   string `json:"provider"`
}

type GenerateInsightsRequest struct {
	Data map[string]interface{} `json:"data"`
}

type GeneratedInsight struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Explanation string `json:"explanation"`
	Confidence  string `json:"confidence"`
}

type GenerateInsightsResponse struct {
	Insights []GeneratedInsight `json:"insights"`
	Provider string             `json:"provider"`
}

type HealthResponse struct {
	Status   string `json:"status"`
	Provider string `json:"provider"`
	Message  string `json:"message,omitempty"`
}

// Capabilities defines what a provider supports.
type Capabilities struct {
	Provider     string   `json:"provider"`
	Chat         bool     `json:"chat"`
	Explanation  bool     `json:"explanation"`
	Summarize    bool     `json:"summarize"`
	Insights     bool     `json:"insights"`
	Streaming    bool     `json:"streaming"`
	ToolCalling  bool     `json:"tool_calling"`
	Embeddings   bool     `json:"embeddings"`
	Vision       bool     `json:"vision"`
	MaxContext   int      `json:"max_context"`
	Models       []string `json:"models,omitempty"`
}

// ProviderInfo holds metadata and status for a registered provider.
type ProviderInfo struct {
	Name         string       `json:"name"`
	Capabilities Capabilities `json:"capabilities"`
	Status       string       `json:"status"`
	Healthy      bool         `json:"healthy"`
	Error        string       `json:"error,omitempty"`
}

// ProviderError represents a typed error from a provider.
type ProviderError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Provider  string `json:"provider"`
	Retryable bool   `json:"retryable"`
}

func (e *ProviderError) Error() string { return e.Message }

var (
	ErrNotImplemented   = &ProviderError{Code: "NOT_IMPLEMENTED", Message: "Not implemented in this provider", Retryable: false}
	ErrProviderNotFound = &ProviderError{Code: "PROVIDER_NOT_FOUND", Message: "Provider not found in registry", Retryable: false}
)
