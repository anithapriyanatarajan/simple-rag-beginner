package config

import (
	"os"
	"strconv"
)

// ModelConfig holds configuration for AI models
type ModelConfig struct {
	// Generative Model Settings (Ollama only)
	GenerativeProvider string // "ollama" (only supported provider)
	GenerativeModel    string // e.g., "llama3.2:latest", "mistral:latest"
	GenerativeEndpoint string // Ollama API endpoint URL
	GenerativeAPIKey   string // Not used for Ollama (for future compatibility)

	// Embedding Model Settings (Ollama only)
	EmbeddingProvider string // "ollama" (only supported provider)
	EmbeddingModel    string // e.g., "nomic-embed-text:latest", "all-minilm:l6-v2"
	EmbeddingEndpoint string // Ollama API endpoint URL
	EmbeddingAPIKey   string // Not used for Ollama (for future compatibility)

	// Generation Settings
	MaxTokens    int     // Maximum tokens for generation
	Temperature  float64 // Temperature for generation (0.0-1.0)
	SystemPrompt string  // System prompt for the AI

	// Qdrant Vector Database Settings
	QdrantHost string // Qdrant host
	QdrantPort int    // Qdrant gRPC port
}

// DefaultConfig returns the default configuration
func DefaultConfig() *ModelConfig {
	return &ModelConfig{
		// Ollama is the only supported provider
		GenerativeProvider: getEnv("GENERATIVE_PROVIDER", "ollama"),
		GenerativeModel:    getEnv("GENERATIVE_MODEL", "llama3.2:latest"),
		GenerativeEndpoint: getEnv("GENERATIVE_ENDPOINT", "http://localhost:11434"),
		GenerativeAPIKey:   getEnv("GENERATIVE_API_KEY", ""),

		// Ollama embeddings
		EmbeddingProvider: getEnv("EMBEDDING_PROVIDER", "ollama"),
		EmbeddingModel:    getEnv("EMBEDDING_MODEL", "nomic-embed-text:latest"),
		EmbeddingEndpoint: getEnv("EMBEDDING_ENDPOINT", "http://localhost:11434"),
		EmbeddingAPIKey:   getEnv("EMBEDDING_API_KEY", ""),

		// Generation parameters
		MaxTokens:   getEnvInt("MAX_TOKENS", 512),
		Temperature: getEnvFloat("TEMPERATURE", 0.7),
		SystemPrompt: getEnv("SYSTEM_PROMPT",
			"You are a helpful AI assistant. Use the provided context to answer questions accurately and concisely. If the context doesn't contain relevant information, say so clearly."),

		// Qdrant configuration
		QdrantHost: getEnv("QDRANT_HOST", "localhost"),
		QdrantPort: getEnvInt("QDRANT_PORT", 6334),
	}
}

// getEnv gets environment variable or returns default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets environment variable as int or returns default
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvFloat gets environment variable as float64 or returns default
func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}
