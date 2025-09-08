package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"simple-rag-beginner/cmd/rag-api/config"
)

// OllamaLLM implements LLMProvider for Ollama (local models)
type OllamaLLM struct {
	baseURL string
	config  *config.Config
}

// NewOllamaLLM creates a new Ollama LLM provider
func NewOllamaLLM(cfg *config.Config) (*OllamaLLM, error) {
	baseURL := cfg.LLMBaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434" // Default Ollama port
	}

	return &OllamaLLM{
		baseURL: baseURL,
		config:  cfg,
	}, nil
}

// OllamaRequest represents the request structure for Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents the response structure from Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// GenerateResponse generates a response using Ollama
func (o *OllamaLLM) GenerateResponse(ctx context.Context, prompt string, config *config.Config) (string, error) {
	reqBody := OllamaRequest{
		Model:  o.config.LLMModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(o.baseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Ollama API error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API returned status: %d", resp.StatusCode)
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Response, nil
}

// GetName returns the provider name
func (o *OllamaLLM) GetName() string {
	return "ollama"
}

// OllamaEmbedding implements EmbeddingProvider for Ollama
type OllamaEmbedding struct {
	baseURL string
	config  *config.Config
}

// NewOllamaEmbedding creates a new Ollama embedding provider
func NewOllamaEmbedding(cfg *config.Config) (*OllamaEmbedding, error) {
	baseURL := cfg.EmbeddingBaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434" // Default Ollama port
	}

	return &OllamaEmbedding{
		baseURL: baseURL,
		config:  cfg,
	}, nil
}

// OllamaEmbeddingRequest represents the embedding request structure
type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaEmbeddingResponse represents the embedding response structure
type OllamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// GenerateEmbedding generates embeddings using Ollama
func (o *OllamaEmbedding) GenerateEmbedding(ctx context.Context, text string, config *config.Config) ([]float32, error) {
	reqBody := OllamaEmbeddingRequest{
		Model:  o.config.EmbeddingModel,
		Prompt: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(o.baseURL+"/api/embeddings", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Ollama embedding API error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama embedding API returned status: %d", resp.StatusCode)
	}

	var ollamaResp OllamaEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to decode embedding response: %w", err)
	}

	return ollamaResp.Embedding, nil
}

// GetName returns the provider name
func (o *OllamaEmbedding) GetName() string {
	return "ollama"
}

// GetDimensions returns the embedding dimensions
func (o *OllamaEmbedding) GetDimensions() int {
	return o.config.EmbeddingDimensions
}
