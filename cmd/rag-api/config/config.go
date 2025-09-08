package config

import (
	"fmt"
	"os"
	"strconv"
)

// LLMProvider represents different LLM providers
type LLMProvider string

const (
	OpenAI    LLMProvider = "openai"
	Anthropic LLMProvider = "anthropic"
	Cohere    LLMProvider = "cohere"
	Ollama    LLMProvider = "ollama"
	Azure     LLMProvider = "azure"
)

// EmbeddingProvider represents different embedding providers
type EmbeddingProvider string

const (
	OpenAIEmbedding EmbeddingProvider = "openai"
	CohereEmbedding EmbeddingProvider = "cohere"
	HuggingFace     EmbeddingProvider = "huggingface"
	OllamaEmbedding EmbeddingProvider = "ollama"
	AzureEmbedding  EmbeddingProvider = "azure"
)

// Config holds all configuration for the RAG system
type Config struct {
	// LLM Configuration
	LLMProvider    LLMProvider `json:"llm_provider"`
	LLMModel       string      `json:"llm_model"`
	LLMAPIKey      string      `json:"llm_api_key"`
	LLMBaseURL     string      `json:"llm_base_url"`
	LLMMaxTokens   int         `json:"llm_max_tokens"`
	LLMTemperature float32     `json:"llm_temperature"`

	// Embedding Configuration
	EmbeddingProvider   EmbeddingProvider `json:"embedding_provider"`
	EmbeddingModel      string            `json:"embedding_model"`
	EmbeddingAPIKey     string            `json:"embedding_api_key"`
	EmbeddingBaseURL    string            `json:"embedding_base_url"`
	EmbeddingDimensions int               `json:"embedding_dimensions"`

	// Vector DB Configuration
	QdrantURL string `json:"qdrant_url"`

	// Processing Configuration
	ChunkSize int `json:"chunk_size"`
	TopK      int `json:"top_k"`
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		// LLM Configuration
		LLMProvider:    LLMProvider(getEnvOrDefault("LLM_PROVIDER", "openai")),
		LLMModel:       getEnvOrDefault("LLM_MODEL", getDefaultLLMModel(LLMProvider(getEnvOrDefault("LLM_PROVIDER", "openai")))),
		LLMAPIKey:      getEnvOrDefault("LLM_API_KEY", getEnvOrDefault("OPENAI_API_KEY", "")),
		LLMBaseURL:     getEnvOrDefault("LLM_BASE_URL", ""),
		LLMMaxTokens:   getEnvIntOrDefault("LLM_MAX_TOKENS", 500),
		LLMTemperature: getEnvFloatOrDefault("LLM_TEMPERATURE", 0.7),

		// Embedding Configuration
		EmbeddingProvider:   EmbeddingProvider(getEnvOrDefault("EMBEDDING_PROVIDER", "openai")),
		EmbeddingModel:      getEnvOrDefault("EMBEDDING_MODEL", getDefaultEmbeddingModel(EmbeddingProvider(getEnvOrDefault("EMBEDDING_PROVIDER", "openai")))),
		EmbeddingAPIKey:     getEnvOrDefault("EMBEDDING_API_KEY", getEnvOrDefault("OPENAI_API_KEY", "")),
		EmbeddingBaseURL:    getEnvOrDefault("EMBEDDING_BASE_URL", ""),
		EmbeddingDimensions: getEnvIntOrDefault("EMBEDDING_DIMENSIONS", getDefaultEmbeddingDimensions(EmbeddingProvider(getEnvOrDefault("EMBEDDING_PROVIDER", "openai")))),

		// Vector DB Configuration
		QdrantURL: getEnvOrDefault("QDRANT_URL", "http://qdrant:6333"),

		// Processing Configuration
		ChunkSize: getEnvIntOrDefault("CHUNK_SIZE", 2000),
		TopK:      getEnvIntOrDefault("TOP_K", 5),
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.LLMAPIKey == "" && c.LLMProvider != Ollama {
		return fmt.Errorf("LLM API key is required for provider: %s", c.LLMProvider)
	}

	if c.EmbeddingAPIKey == "" && c.EmbeddingProvider != OllamaEmbedding && c.EmbeddingProvider != HuggingFace {
		return fmt.Errorf("embedding API key is required for provider: %s", c.EmbeddingProvider)
	}

	if c.QdrantURL == "" {
		return fmt.Errorf("Qdrant URL is required")
	}

	return nil
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloatOrDefault(key string, defaultValue float32) float32 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 32); err == nil {
			return float32(floatValue)
		}
	}
	return defaultValue
}

func getDefaultLLMModel(provider LLMProvider) string {
	switch provider {
	case OpenAI:
		return "gpt-3.5-turbo"
	case Anthropic:
		return "claude-3-haiku-20240307"
	case Cohere:
		return "command-light"
	case Ollama:
		return "llama2"
	case Azure:
		return "gpt-35-turbo"
	default:
		return "gpt-3.5-turbo"
	}
}

func getDefaultEmbeddingModel(provider EmbeddingProvider) string {
	switch provider {
	case OpenAIEmbedding:
		return "text-embedding-ada-002"
	case CohereEmbedding:
		return "embed-english-v3.0"
	case HuggingFace:
		return "sentence-transformers/all-MiniLM-L6-v2"
	case OllamaEmbedding:
		return "nomic-embed-text"
	case AzureEmbedding:
		return "text-embedding-ada-002"
	default:
		return "text-embedding-ada-002"
	}
}

func getDefaultEmbeddingDimensions(provider EmbeddingProvider) int {
	switch provider {
	case OpenAIEmbedding:
		return 1536
	case CohereEmbedding:
		return 1024
	case HuggingFace:
		return 384
	case OllamaEmbedding:
		return 768
	case AzureEmbedding:
		return 1536
	default:
		return 1536
	}
}

// GetProvidersList returns available providers for display
func GetLLMProviders() []string {
	return []string{string(OpenAI), string(Anthropic), string(Cohere), string(Ollama), string(Azure)}
}

func GetEmbeddingProviders() []string {
	return []string{string(OpenAIEmbedding), string(CohereEmbedding), string(HuggingFace), string(OllamaEmbedding), string(AzureEmbedding)}
}
