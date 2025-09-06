package model

import (
	"context"
	"fmt"
	"log"
	"simple-rag-beginner/internal/config"
	"strings"
	"time"
)

// AIProvider interface for different AI providers
type AIProvider interface {
	Generate(ctx context.Context, prompt string, maxTokens int, temperature float64) (string, error)
	IsAvailable(ctx context.Context) error
}

// ModelService handles AI model operations
type ModelService struct {
	provider AIProvider
	config   *config.ModelConfig
}

var globalModelService *ModelService

// InitModelService initializes the global model service
func InitModelService(cfg *config.ModelConfig) error {
	var provider AIProvider

	switch cfg.GenerativeProvider {
	case "ollama":
		provider = NewOllamaClient(cfg)
	default:
		return fmt.Errorf("unsupported generative provider: %s (only 'ollama' is supported)", cfg.GenerativeProvider)
	}

	// Test if provider is available - this is required, no fallbacks
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := provider.IsAvailable(ctx); err != nil {
		return fmt.Errorf("ollama is required but not available: %w\n\nTo fix this:\n1. Install Ollama: curl -fsSL https://ollama.ai/install.sh | sh\n2. Start Ollama: ollama serve\n3. Pull model: ollama pull %s\n4. Restart the application", err, cfg.GenerativeModel)
	}

	globalModelService = &ModelService{
		provider: provider,
		config:   cfg,
	}

	log.Printf("✅ Ollama connected successfully with model: %s", cfg.GenerativeModel)
	return nil
}

// GenerateText generates text using the configured AI model
func GenerateText(query, contextStr string) string {
	if globalModelService == nil {
		log.Fatal("model service not initialized - this should not happen")
	}

	return globalModelService.GenerateWithContext(query, contextStr)
}

// GenerateWithContext generates a response using the provided query and context
func (ms *ModelService) GenerateWithContext(query, contextStr string) string {
	// Create prompt with system instruction, context, and query
	prompt := ms.buildPrompt(query, contextStr)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// Generate with AI model - no fallbacks
	response, err := ms.provider.Generate(ctx, prompt, ms.config.MaxTokens, ms.config.Temperature)
	if err != nil {
		log.Printf("AI generation failed: %v", err)
		return fmt.Sprintf("Error generating response: %v", err)
	}

	// Clean up the response
	response = strings.TrimSpace(response)
	if response == "" {
		return "Generated response was empty. Please try rephrasing your question."
	}

	return response
}

// buildPrompt constructs the full prompt for the AI model
func (ms *ModelService) buildPrompt(query, contextStr string) string {
	var prompt strings.Builder

	// Add system prompt
	prompt.WriteString(ms.config.SystemPrompt)
	prompt.WriteString("\n\n")

	// Add context if available
	if contextStr != "" && contextStr != "No context available" {
		prompt.WriteString("Context Information:\n")
		prompt.WriteString(contextStr)
		prompt.WriteString("\n\n")
	}

	// Add user query
	prompt.WriteString("User Question: ")
	prompt.WriteString(query)
	prompt.WriteString("\n\n")

	// Add instruction for response
	prompt.WriteString("Please provide a helpful and accurate response based on the context provided above:")

	return prompt.String()
}

// GetModelInfo returns information about the current model configuration
func GetModelInfo() map[string]interface{} {
	if globalModelService == nil {
		return map[string]interface{}{
			"error":  "model service not initialized",
			"status": "error",
		}
	}

	return map[string]interface{}{
		"generative_provider": globalModelService.config.GenerativeProvider,
		"generative_model":    globalModelService.config.GenerativeModel,
		"generative_endpoint": globalModelService.config.GenerativeEndpoint,
		"embedding_provider":  globalModelService.config.EmbeddingProvider,
		"embedding_model":     globalModelService.config.EmbeddingModel,
		"max_tokens":          globalModelService.config.MaxTokens,
		"temperature":         globalModelService.config.Temperature,
		"status":              "connected",
	}
}
