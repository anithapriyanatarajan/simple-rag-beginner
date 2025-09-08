package providers

import (
	"context"
	"fmt"
	"simple-rag-beginner/cmd/rag-api/config"

	"github.com/sashabaranov/go-openai"
)

// OpenAILLM implements LLMProvider for OpenAI
type OpenAILLM struct {
	client *openai.Client
	config *config.Config
}

// NewOpenAILLM creates a new OpenAI LLM provider
func NewOpenAILLM(cfg *config.Config) (*OpenAILLM, error) {
	if cfg.LLMAPIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	clientConfig := openai.DefaultConfig(cfg.LLMAPIKey)
	if cfg.LLMBaseURL != "" {
		clientConfig.BaseURL = cfg.LLMBaseURL
	}

	return &OpenAILLM{
		client: openai.NewClientWithConfig(clientConfig),
		config: cfg,
	}, nil
}

// GenerateResponse generates a response using OpenAI's chat completion
func (o *OpenAILLM) GenerateResponse(ctx context.Context, prompt string, config *config.Config) (string, error) {
	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: o.config.LLMModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		MaxTokens:   o.config.LLMMaxTokens,
		Temperature: o.config.LLMTemperature,
	})

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response generated from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// GetName returns the provider name
func (o *OpenAILLM) GetName() string {
	return "openai"
}

// OpenAIEmbedding implements EmbeddingProvider for OpenAI
type OpenAIEmbedding struct {
	client *openai.Client
	config *config.Config
}

// NewOpenAIEmbedding creates a new OpenAI embedding provider
func NewOpenAIEmbedding(cfg *config.Config) (*OpenAIEmbedding, error) {
	if cfg.EmbeddingAPIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required for embeddings")
	}

	clientConfig := openai.DefaultConfig(cfg.EmbeddingAPIKey)
	if cfg.EmbeddingBaseURL != "" {
		clientConfig.BaseURL = cfg.EmbeddingBaseURL
	}

	return &OpenAIEmbedding{
		client: openai.NewClientWithConfig(clientConfig),
		config: cfg,
	}, nil
}

// GenerateEmbedding generates embeddings using OpenAI's embedding API
func (o *OpenAIEmbedding) GenerateEmbedding(ctx context.Context, text string, config *config.Config) ([]float32, error) {
	resp, err := o.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.EmbeddingModel(o.config.EmbeddingModel),
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI embedding API error: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned from OpenAI")
	}

	return resp.Data[0].Embedding, nil
}

// GetName returns the provider name
func (o *OpenAIEmbedding) GetName() string {
	return "openai"
}

// GetDimensions returns the embedding dimensions
func (o *OpenAIEmbedding) GetDimensions() int {
	return o.config.EmbeddingDimensions
}
