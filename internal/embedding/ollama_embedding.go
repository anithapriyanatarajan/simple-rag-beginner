package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaEmbeddingClient handles embedding requests to Ollama
type OllamaEmbeddingClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
}

// NewOllamaEmbeddingClient creates a new Ollama embedding client
func NewOllamaEmbeddingClient(endpoint, model string) *OllamaEmbeddingClient {
	return &OllamaEmbeddingClient{
		baseURL: endpoint,
		model:   model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// OllamaEmbeddingRequest represents the request structure for Ollama embeddings
type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaEmbeddingResponse represents the response structure from Ollama embeddings
type OllamaEmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// GetEmbedding gets an embedding vector from Ollama
func (c *OllamaEmbeddingClient) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	request := OllamaEmbeddingRequest{
		Model:  c.model,
		Prompt: text,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embedding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send embedding request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama embedding API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedding response: %w", err)
	}

	var response OllamaEmbeddingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal embedding response: %w", err)
	}

	// Convert float64 to float32
	embedding := make([]float32, len(response.Embedding))
	for i, val := range response.Embedding {
		embedding[i] = float32(val)
	}

	return embedding, nil
}

// IsAvailable checks if the embedding model is available
func (c *OllamaEmbeddingClient) IsAvailable(ctx context.Context) error {
	// Test with a simple embedding request
	_, err := c.GetEmbedding(ctx, "test")
	return err
}
