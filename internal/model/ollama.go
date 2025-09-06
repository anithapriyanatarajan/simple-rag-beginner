package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"simple-rag-beginner/internal/config"
	"time"
)

// OllamaClient handles communication with Ollama API
type OllamaClient struct {
	baseURL    string
	httpClient *http.Client
	model      string
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(cfg *config.ModelConfig) *OllamaClient {
	return &OllamaClient{
		baseURL: cfg.GenerativeEndpoint,
		model:   cfg.GenerativeModel,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// OllamaGenerateRequest represents the request structure for Ollama
type OllamaGenerateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// OllamaGenerateResponse represents the response structure from Ollama
type OllamaGenerateResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	Context            []int     `json:"context"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

// Generate sends a prompt to Ollama and returns the generated text
func (c *OllamaClient) Generate(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error) {
	request := OllamaGenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"num_predict": maxTokens,
			"temperature": temperature,
			"top_p":       0.9,
			"stop":        []string{},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var response OllamaGenerateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Response, nil
}

// IsAvailable checks if Ollama is running and the model is available
func (c *OllamaClient) IsAvailable(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ollama not available: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	// Check if our specific model is available
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var tagsResponse struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.Unmarshal(body, &tagsResponse); err != nil {
		return fmt.Errorf("failed to parse tags response: %w", err)
	}

	// Check if our model is in the list
	for _, model := range tagsResponse.Models {
		if model.Name == c.model {
			return nil // Model found
		}
	}

	return fmt.Errorf("model '%s' not found in Ollama. Available models: %v", c.model, tagsResponse.Models)
}
