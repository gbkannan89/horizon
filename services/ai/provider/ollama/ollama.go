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

	"github.com/horizon/core/services/ai/provider"
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
	// Simple token/context limit: keep at most last 10 messages
	hist := req.History
	if len(hist) > 10 {
		hist = hist[len(hist)-10:]
	}
	for _, h := range hist {
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

func (o *OllamaProvider) StreamChat(ctx context.Context, req provider.ChatRequest) (<-chan string, error) {
	if !o.enabled {
		return nil, provider.ErrNotImplemented
	}

	messages := []map[string]string{}
	hist := req.History
	if len(hist) > 10 {
		hist = hist[len(hist)-10:]
	}
	for _, h := range hist {
		messages = append(messages, map[string]string{"role": h.Role, "content": h.Content})
	}
	messages = append(messages, map[string]string{"role": "user", "content": req.Message})

	body := map[string]interface{}{
		"model":    o.model,
		"messages": messages,
		"stream":   true,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", o.endpoint+"/api/chat", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, &provider.ProviderError{
			Code: "HTTP_ERROR", Message: fmt.Sprintf("Ollama error %d: %s", resp.StatusCode, string(respBody)),
			Provider: o.name, Retryable: false,
		}
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var chunk struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
				Done bool `json:"done"`
			}
			if err := decoder.Decode(&chunk); err != nil {
				break
			}
			if chunk.Message.Content != "" {
				select {
				case ch <- chunk.Message.Content:
				case <-ctx.Done():
					return
				}
			}
			if chunk.Done {
				break
			}
		}
	}()

	return ch, nil
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
	models := o.discoverModels(context.Background())
	return provider.Capabilities{
		Provider: o.name, Chat: true, Explanation: true, Summarize: true, Insights: true,
		Streaming: true, ToolCalling: false, Embeddings: false, Vision: false,
		MaxContext: 8192, Models: models,
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

func (o *OllamaProvider) GenerateInsights(ctx context.Context, req provider.GenerateInsightsRequest) (*provider.GenerateInsightsResponse, error) {
	// A simple prompt asking Ollama to generate insights in JSON format
	prompt := "Analyze the following financial data and provide proactive insights (like spending patterns, savings opportunities). Format as JSON array with type, title, summary, explanation, and confidence."
	dataStr, _ := json.Marshal(req.Data)
	prompt += "\nData: " + string(dataStr)

	apiReq := map[string]interface{}{
		"model":  "llama3",
		"prompt": prompt,
		"stream": false,
		"format": "json",
	}

	var apiResp struct {
		Response string `json:"response"`
	}
	if err := o.doRequest(ctx, "/api/generate", apiReq, &apiResp); err != nil {
		return nil, err
	}

	var insights []provider.GeneratedInsight
	if err := json.Unmarshal([]byte(apiResp.Response), &insights); err != nil {
		// Fallback to basic if parsing fails
		insights = []provider.GeneratedInsight{
			{
				Type:        "general",
				Title:       "AI Insight Generated",
				Summary:     "Generated via Ollama.",
				Explanation: apiResp.Response,
				Confidence:  "Medium",
			},
		}
	}

	return &provider.GenerateInsightsResponse{
		Insights: insights,
		Provider: o.name,
	}, nil
}
