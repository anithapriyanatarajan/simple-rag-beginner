package embedding

import (
	"context"
	"fmt"
	"log"
	"simple-rag-beginner/internal/config"
	"time"
)

var embeddingClient *OllamaEmbeddingClient

// InitEmbeddingService initializes the embedding service with Ollama
func InitEmbeddingService(cfg *config.ModelConfig) error {
	if cfg.EmbeddingProvider != "ollama" {
		return fmt.Errorf("only ollama embedding provider is supported, got: %s", cfg.EmbeddingProvider)
	}

	embeddingClient = NewOllamaEmbeddingClient(cfg.EmbeddingEndpoint, cfg.EmbeddingModel)

	// Test if Ollama embedding is available - this is required
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := embeddingClient.IsAvailable(ctx); err != nil {
		return fmt.Errorf("ollama embedding is required but not available: %w\n\nTo fix this:\n1. Ensure Ollama is running: ollama serve\n2. Pull embedding model: ollama pull %s\n3. Restart the application", err, cfg.EmbeddingModel)
	}

	log.Printf("✅ Ollama embedding connected successfully with model: %s", cfg.EmbeddingModel)
	return nil
}

// TextToVector converts text to vector using Ollama embeddings
func TextToVector(text string) []float32 {
	if embeddingClient == nil {
		log.Fatal("embedding service not initialized - this should not happen")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	embedding, err := embeddingClient.GetEmbedding(ctx, text)
	if err != nil {
		log.Printf("Error generating embedding: %v", err)
		// Return zero vector as fallback
		return make([]float32, 768) // nomic-embed-text dimension
	}

	return embedding
}
