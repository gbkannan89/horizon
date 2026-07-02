package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/horizon/core/services/ai/internal/provider"
)

// OllamaProvider communicates with a local Ollama server.
type OllamaProvider struct {
	name    string
	endpoint string
	client  *http.Client
	timeout time.Duration
	retries int
	model   string
	enabled bool
}

type Config struct {
	Endpoint string
	Timeout  time.Duration
	Retries  int
	Model    string
	Enabled  bool
}

func New(cfg Config) *OllamaProvider {
	if cfg.Endpoint == "" { cfg.Endpoint = "http://localhost:11434" }
	if cfg.Timeout <= 0 { cfg.Timeout = 30 * time.Second }
	if cfg.Retries <= 0 { cfg.Retries = 2 }
	if cfg.Model == "" { cfg.Model = "llama3" }

	return &OllamaProvider{
		name:     "ollama",
		endpoint: strings.TrimRight(cfg.Endpoint, "/"),
		client: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: &http.Transport{MaxIdleConns: 5, IdleConnTimeout: 30 * time.Second},
		},
		timeout: cfg.Timeout,
		retries: cfg.Retries,
		model:   cfg.Model,
		enabled: cfg.Enabled,
	}
}

func (o *OllamaProvider) Chat(ctx context.Context, req provider.ChatRequest) (*provider.ChatResponse, error) {
	if !o.enabled { return nil, provider.ErrNotImplemented }

	messages := []map[string]string{}
	for _, h := range req.History {
		messages = append(messages, map[string]string{"role": h.Role, "content": h.Content})
	}
	messages = append(messages, map[string]string{"role": "user", "content": req.Message})

	body := map[string]interface{}{
		"model":    o.model,
		"messages": messages,
		"stream":   false,
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		Done bool `json:"done"`
	}

	if err := o.doRequest(ctx, "/api/chat", body, &result); err != nil {
		return nil, err
	}

	return &provider.ChatResponse{
		Reply:      result.Message.Content,
		Confidence: "Medium",
		Provider:   o.name,
		SessionID:  req.SessionID,
	}, nil
}

func (o *OllamaProvider) Explain(ctx context.Context, req provider.ExplainRequest) (*provider.ExplainResponse, error) {
	if !o.enabled { return nil, provider.ErrNotImplemented }

	prompt := fmt.Sprintf("Provide a clear explanation about %s based on the following data: %v. Keep the explanation concise and helpful.", req.PromptType, req.Data)

	body := map[string]interface{}{
		"model":  o.model,
		"prompt": prompt,
		"stream": false,
	}

	var result struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	if err := o.doRequest(ctx, "/api/generate", body, &result); err != nil {
		return nil, err
	}

	return &provider.ExplainResponse{
		Explanation: result.Response,
		Confidence:  "Medium",
		Provider:    o.name,
		Version:     req.Version,
	}, nil
}

func (o *OllamaProvider) Summarize(ctx context.Context, req provider.SummarizeRequest) (*provider.SummarizeResponse, error) {
	if !o.enabled { return nil, provider.ErrNotImplemented }

	prompt := fmt.Sprintf("Summarize the following financial information about '%s': %v", req.Topic, req.Data)

	body := map[string]interface{}{
		"model":  o.model,
		"prompt": prompt,
		"stream": false,
	}

	var result struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	if err := o.doRequest(ctx, "/api/generate", body, &result); err != nil {
		return nil, err
	}

	return &provider.SummarizeResponse{
		Summary:    result.Response,
		Confidence: "Medium",
		Provider:   o.name,
	}, nil
}

func (o *OllamaProvider) Health(ctx context.Context) (*provider.HealthResponse, error) {
	if !o.enabled {
		return &provider.HealthResponse{Status: "disabled", Provider: o.name, Message: "Ollama provider is disabled via configuration."}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", o.endpoint+"/api/tags", nil)
	if err != nil {
		return &provider.HealthResponse{Status: "error", Provider: o.name, Message: err.Error()}, fmt.Errorf("health check failed: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return &provider.HealthResponse{Status: "unavailable", Provider: o.name, Message: fmt.Sprintf("Cannot reach Ollama at %s: %v", o.endpoint, err)},
			fmt.Errorf("ollama unavailable: %w", err)
	}
	defer resp.Body.Close()

	io.Copy(io.Discard, resp.Body)

	var models []string
	models = o.discoverModels(ctx)
	if len(models) == 0 { models = []string{o.model} }

	return &provider.HealthResponse{
		Status:   "operational",
		Provider: o.name,
		Message:  fmt.Sprintf("Ollama running at %s with model %s", o.endpoint, strings.Join(models, ", ")),
	}, nil
}

func (o *OllamaProvider) Capabilities() provider.Capabilities {
	return provider.Capabilities{
		Provider: o.name, Chat: o.enabled, Explanation: o.enabled, Summarize: o.enabled,
		Streaming: true, ToolCalling: false, Embeddings: true, Vision: false,
		MaxContext: 8192, Models: []string{o.model},
	}
}

func (o *OllamaProvider) Enabled() bool { return o.enabled }

func (o *OllamaProvider) doRequest(ctx context.Context, path string, body interface{}, result interface{}) error {
	var lastErr error
	for attempt := 0; attempt <= o.retries; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(time.Duration(attempt) * 500 * time.Millisecond):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		jsonBody, err := json.Marshal(body)
		if err != nil { return fmt.Errorf("marshal request: %w", err) }

		req, err := http.NewRequestWithContext(ctx, "POST", o.endpoint+path, bytes.NewReader(jsonBody))
		if err != nil { return fmt.Errorf("create request: %w", err) }
		req.Header.Set("Content-Type", "application/json")

		resp, err := o.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: %w", attempt+1, err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: read response: %w", attempt+1, err)
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("attempt %d: server error %d: %s", attempt+1, resp.StatusCode, string(respBody))
			continue
		}

		if resp.StatusCode >= 400 {
			return &provider.ProviderError{
				Code: "HTTP_ERROR", Message: fmt.Sprintf("Ollama error %d: %s", resp.StatusCode, string(respBody)),
				Provider: o.name, Retryable: false,
			}
		}

		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}
		return nil
	}
	return lastErr
}

func (o *OllamaProvider) discoverModels(ctx context.Context) []string {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", o.endpoint+"/api/tags", nil)
	if err != nil { return nil }

	resp, err := o.client.Do(req)
	if err != nil { return nil }
	defer resp.Body.Close()

	var tags struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil { return nil }

	var models []string
	for _, m := range tags.Models { models = append(models, m.Name) }
	return models
}
