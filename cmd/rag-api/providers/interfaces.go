package providers

import (
	"context"
	"simple-rag-beginner/cmd/rag-api/config"
)

// LLMProvider interface for different language model providers
type LLMProvider interface {
	GenerateResponse(ctx context.Context, prompt string, config *config.Config) (string, error)
	GetName() string
}

// EmbeddingProvider interface for different embedding providers
type EmbeddingProvider interface {
	GenerateEmbedding(ctx context.Context, text string, config *config.Config) ([]float32, error)
	GetName() string
	GetDimensions() int
}

// LLMProviderFactory creates LLM providers based on configuration
func NewLLMProvider(cfg *config.Config) (LLMProvider, error) {
	switch cfg.LLMProvider {
	case config.OpenAI:
		return NewOpenAILLM(cfg)
	case config.Anthropic:
		return NewAnthropicLLM(cfg)
	case config.Cohere:
		return NewCohereLLM(cfg)
	case config.Ollama:
		return NewOllamaLLM(cfg)
	case config.Azure:
		return NewAzureLLM(cfg)
	default:
		return NewOpenAILLM(cfg) // Default to OpenAI
	}
}

// EmbeddingProviderFactory creates embedding providers based on configuration
func NewEmbeddingProvider(cfg *config.Config) (EmbeddingProvider, error) {
	switch cfg.EmbeddingProvider {
	case config.OpenAIEmbedding:
		return NewOpenAIEmbedding(cfg)
	case config.CohereEmbedding:
		return NewCohereEmbedding(cfg)
	case config.HuggingFace:
		return NewHuggingFaceEmbedding(cfg)
	case config.OllamaEmbedding:
		return NewOllamaEmbedding(cfg)
	case config.AzureEmbedding:
		return NewAzureEmbedding(cfg)
	default:
		return NewOpenAIEmbedding(cfg) // Default to OpenAI
	}
}
