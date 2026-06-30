package api

import (
	"github.com/horizon/core/services/ai/internal/context"
	"github.com/horizon/core/services/ai/internal/prompts"
)

type HealthResponse struct {
	Success  bool           `json:"success"`
	Data     *AIHealth      `json:"data,omitempty"`
	Metadata *Metadata      `json:"metadata,omitempty"`
}

type AIHealth struct {
	Status   string `json:"status"`
	Provider string `json:"provider"`
	Message  string `json:"message,omitempty"`
}

type ChatRequest struct {
	SessionID string                 `json:"session_id"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

type ChatResponse struct {
	Success  bool        `json:"success"`
	Data     *ChatResult `json:"data,omitempty"`
	Metadata *Metadata   `json:"metadata,omitempty"`
}

type ChatResult struct {
	Reply      string `json:"reply"`
	Confidence string `json:"confidence"`
	Provider   string `json:"provider"`
	SessionID  string `json:"session_id"`
}

type ExplainRequest struct {
	PromptType string                 `json:"prompt_type"`
	Data       map[string]interface{} `json:"data"`
	Version    int                    `json:"version,omitempty"`
}

type ExplainResponse struct {
	Success  bool           `json:"success"`
	Data     *ExplainResult `json:"data,omitempty"`
	Metadata *Metadata      `json:"metadata,omitempty"`
}

type ExplainResult struct {
	Explanation string `json:"explanation"`
	Confidence  string `json:"confidence"`
	Provider    string `json:"provider"`
	Version     int    `json:"version"`
}

type PromptListResponse struct {
	Success  bool                    `json:"success"`
	Data     []prompts.TemplateSummary `json:"data,omitempty"`
	Metadata *Metadata               `json:"metadata,omitempty"`
}

type ContextResponse struct {
	Success  bool                  `json:"success"`
	Data     *context.AdvisorContext `json:"data,omitempty"`
	Metadata *Metadata             `json:"metadata,omitempty"`
}

type Metadata struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}
