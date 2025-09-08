package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"simple-rag-beginner/cmd/rag-api/config"
	"simple-rag-beginner/cmd/rag-api/providers"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "rag_api_requests_total"},
		[]string{"status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "rag_api_request_duration_seconds"},
		[]string{"endpoint"},
	)
	appConfig *config.Config
)

func init() {
	prometheus.MustRegister(requestCounter, requestDuration)
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	// Load configuration
	appConfig = config.LoadConfig()
	if err := appConfig.Validate(); err != nil {
		log.Printf(`{"level":"warn","msg":"Configuration validation failed","error":"%s"}`, err.Error())
		log.Println(`{"level":"info","msg":"Continuing with current configuration - some features may not work"}`)
	}

	log.Printf(`{"level":"info","msg":"rag-api service started","port":8080,"llm_provider":"%s","embedding_provider":"%s"}`,
		appConfig.LLMProvider, appConfig.EmbeddingProvider)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Serve web interface
	r.GET("/", webHandler)
	r.GET("/config", configHandler)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/health", healthHandler)
	r.POST("/query", queryHandler)              // RAG query with vector DB
	r.POST("/query-direct", directQueryHandler) // Direct LLM query without vector DB
	r.POST("/crawl-and-query", crawlAndQueryHandler)

	r.Run(":8080")
}

func queryHandler(c *gin.Context) {
	start := time.Now()
	defer func() {
		requestDuration.WithLabelValues("query").Observe(time.Since(start).Seconds())
	}()

	var req struct {
		Query string `json:"query"`
	}

	if err := c.BindJSON(&req); err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Initialize providers
	embeddingProvider, err := providers.NewEmbeddingProvider(appConfig)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to initialize embedding provider: %v", err)})
		return
	}

	llmProvider, err := providers.NewLLMProvider(appConfig)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to initialize LLM provider: %v", err)})
		return
	}

	// 1. Get query embedding
	embedding, err := embeddingProvider.GenerateEmbedding(context.Background(), req.Query, appConfig)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to embed query: %v", err)})
		return
	}

	// 2. Search Qdrant for top-k similar chunks
	results, err := searchQdrant(embedding, appConfig.TopK)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to search: %v", err)})
		return
	}

	// 3. Generate response using configured LLM
	response, sources, err := generateResponseWithProvider(llmProvider, req.Query, results)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate response: %v", err)})
		return
	}

	requestCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, gin.H{
		"response":                response,
		"sources":                 sources,
		"query":                   req.Query,
		"llm_provider":            string(appConfig.LLMProvider),
		"embedding_provider":      string(appConfig.EmbeddingProvider),
		"processing_time_seconds": time.Since(start).Seconds(),
	})
}

func directQueryHandler(c *gin.Context) {
	start := time.Now()
	defer func() {
		requestDuration.WithLabelValues("query-direct").Observe(time.Since(start).Seconds())
	}()

	var req struct {
		Query string `json:"query"`
	}

	if err := c.BindJSON(&req); err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Initialize LLM provider
	llmProvider, err := providers.NewLLMProvider(appConfig)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to initialize LLM provider: %v", err)})
		return
	}

	// Direct LLM call without any vector DB retrieval
	response, err := generateDirectResponseWithProvider(llmProvider, req.Query)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate response: %v", err)})
		return
	}

	requestCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, gin.H{
		"response":                response,
		"sources":                 []string{}, // No sources for direct queries
		"query":                   req.Query,
		"mode":                    "direct",
		"llm_provider":            string(appConfig.LLMProvider),
		"processing_time_seconds": time.Since(start).Seconds(),
	})
}

func configHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"llm_provider":         string(appConfig.LLMProvider),
		"embedding_provider":   string(appConfig.EmbeddingProvider),
		"llm_model":            appConfig.LLMModel,
		"embedding_model":      appConfig.EmbeddingModel,
		"top_k":                appConfig.TopK,
		"llm_max_tokens":       appConfig.LLMMaxTokens,
		"llm_temperature":      appConfig.LLMTemperature,
		"qdrant_url":           appConfig.QdrantURL,
		"chunk_size":           appConfig.ChunkSize,
		"embedding_dimensions": appConfig.EmbeddingDimensions,
	})
}

func generateResponseWithProvider(llmProvider providers.LLMProvider, query string, contextResults []map[string]interface{}) (string, []string, error) {
	// Build context from search results
	var contextTexts []string
	var sources []string

	for _, result := range contextResults {
		if payload, ok := result["payload"].(map[string]interface{}); ok {
			if text, ok := payload["text"].(string); ok {
				contextTexts = append(contextTexts, text)
			}
			if source, ok := payload["source"].(string); ok {
				sources = append(sources, source)
			}
		}
	}

	contextString := strings.Join(contextTexts, "\n\n")

	prompt := fmt.Sprintf(`Based on the following context, answer the question. If the answer cannot be found in the context, say so.

Context:
%s

Question: %s

Answer:`, contextString, query)

	response, err := llmProvider.GenerateResponse(context.Background(), prompt, appConfig)
	return response, sources, err
}

func generateDirectResponseWithProvider(llmProvider providers.LLMProvider, query string) (string, error) {
	prompt := fmt.Sprintf(`Please answer the following question directly:

Question: %s

Answer:`, query)

	return llmProvider.GenerateResponse(context.Background(), prompt, appConfig)
}

func searchQdrant(embedding []float32, topK int) ([]map[string]interface{}, error) {
	qdrantURL := os.Getenv("QDRANT_URL")
	if qdrantURL == "" {
		qdrantURL = "http://qdrant:6333"
	}

	searchPayload := map[string]interface{}{
		"vector":       embedding,
		"limit":        topK,
		"with_payload": true,
	}

	searchData, _ := json.Marshal(searchPayload)
	resp, err := http.Post(qdrantURL+"/collections/documents/points/search", "application/json", bytes.NewBuffer(searchData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResp struct {
		Result []struct {
			Payload map[string]interface{} `json:"payload"`
			Score   float64                `json:"score"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, hit := range searchResp.Result {
		hit.Payload["score"] = hit.Score
		results = append(results, hit.Payload)
	}

	return results, nil
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func crawlAndQueryHandler(c *gin.Context) {
	start := time.Now()
	defer func() {
		requestDuration.WithLabelValues("crawl-and-query").Observe(time.Since(start).Seconds())
	}()

	var req struct {
		URL   string `json:"url"`
		Query string `json:"query"`
	}

	if err := c.BindJSON(&req); err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Step 1: Trigger crawling
	crawlerURL := os.Getenv("CRAWLER_URL")
	if crawlerURL == "" {
		crawlerURL = "http://crawler:8080"
	}

	crawlPayload := map[string]string{"url": req.URL}
	crawlData, _ := json.Marshal(crawlPayload)

	resp, err := http.Post(crawlerURL+"/crawl", "application/json", bytes.NewBuffer(crawlData))
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to crawl: %v", err)})
		return
	}
	resp.Body.Close()

	// Step 2: Wait for processing pipeline to complete
	time.Sleep(30 * time.Second)

	// Step 3: Query the processed data
	// Initialize providers
	embeddingProvider, err := providers.NewEmbeddingProvider(appConfig)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to initialize embedding provider: %v", err)})
		return
	}

	llmProvider, err := providers.NewLLMProvider(appConfig)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to initialize LLM provider: %v", err)})
		return
	}

	// Get query embedding
	embedding, err := embeddingProvider.GenerateEmbedding(context.Background(), req.Query, appConfig)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to embed query: %v", err)})
		return
	}

	// Search Qdrant for relevant chunks
	results, err := searchQdrant(embedding, appConfig.TopK)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to search: %v", err)})
		return
	}

	// Generate response
	response, sources, err := generateResponseWithProvider(llmProvider, req.Query, results)
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate response: %v", err)})
		return
	}

	requestCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, gin.H{
		"response":                response,
		"sources":                 sources,
		"query":                   req.Query,
		"crawled_url":             req.URL,
		"processing_time_seconds": time.Since(start).Seconds(),
	})
}

func webHandler(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>RAG Comparison Tool</title>
  <style>
    body { font-family: Arial, sans-serif; background: #f4f4f4; margin: 0; }
    .container { max-width: 800px; margin: 40px auto; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px #ccc; padding: 24px; }
    .mode-selector { margin-bottom: 20px; }
    .mode-selector label { margin-right: 20px; font-weight: bold; }
    .mode-selector input[type=radio] { margin-right: 5px; }
    .chat-box { height: 400px; overflow-y: auto; border: 1px solid #ddd; border-radius: 6px; padding: 12px; background: #fafafa; margin-bottom: 16px; }
    .msg { margin: 8px 0; padding: 8px; border-radius: 4px; }
    .user { background: #e3f2fd; color: #1976d2; }
    .bot { background: #f3e5f5; color: #7b1fa2; }
    .bot.direct { background: #fff3e0; color: #f57c00; }
    .mode-label { font-size: 12px; font-weight: bold; margin-bottom: 4px; }
    form { display: flex; gap: 8px; }
    input[type=text] { flex: 1; padding: 12px; border-radius: 4px; border: 1px solid #ccc; }
    button { padding: 12px 20px; border-radius: 4px; border: none; background: #007bff; color: #fff; cursor: pointer; }
    button:disabled { background: #aaa; }
    .comparison-note { margin-bottom: 16px; padding: 12px; background: #e8f4fd; border-left: 4px solid #2196f3; }
  </style>
</head>
<body>
  <div class="container">
    <h2>RAG Comparison Tool</h2>
    <div class="comparison-note">
      <strong>Compare RAG vs Direct LLM responses:</strong><br>
      • <strong>RAG Mode:</strong> Uses vector DB retrieval + context-aware generation<br>
      • <strong>Direct Mode:</strong> Pure LLM responses without document retrieval
    </div>
    
    <div class="mode-selector">
      <label><input type="radio" name="mode" value="rag" checked> RAG Mode (with vector DB)</label>
      <label><input type="radio" name="mode" value="direct"> Direct Mode (no vector DB)</label>
    </div>
    
    <div class="chat-box" id="chat"></div>
    <form id="chatForm">
      <input type="text" id="query" placeholder="Ask about Tekton or any topic..." autocomplete="off" required />
      <button type="submit">Send</button>
    </form>
  </div>
  <script>
    const chat = document.getElementById('chat');
    const form = document.getElementById('chatForm');
    const queryInput = document.getElementById('query');

    function addMessage(sender, text, mode = '') {
      const msg = document.createElement('div');
      msg.className = 'msg ' + sender + (mode ? ' ' + mode : '');
      
      const modeLabel = document.createElement('div');
      modeLabel.className = 'mode-label';
      if (sender === 'user') {
        modeLabel.textContent = 'You:';
      } else {
        modeLabel.textContent = mode === 'direct' ? 'Bot (Direct LLM):' : 'Bot (RAG):';
      }
      
      const content = document.createElement('div');
      content.textContent = text;
      
      msg.appendChild(modeLabel);
      msg.appendChild(content);
      chat.appendChild(msg);
      chat.scrollTop = chat.scrollHeight;
    }

    form.addEventListener('submit', async (e) => {
      e.preventDefault();
      const query = queryInput.value.trim();
      if (!query) return;
      
      const mode = document.querySelector('input[name="mode"]:checked').value;
      const endpoint = mode === 'direct' ? '/query-direct' : '/query';
      
      addMessage('user', query);
      queryInput.value = '';
      form.querySelector('button').disabled = true;
      
      try {
        const res = await fetch(endpoint, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ query })
        });
        const data = await res.json();
        addMessage('bot', data.response, mode);
        
        if (mode === 'rag' && data.sources && data.sources.length > 0) {
          addMessage('bot', 'Sources: ' + data.sources.join(', '), mode);
        }
      } catch (err) {
        addMessage('bot', 'Error: Could not reach backend.', mode);
      }
      
      form.querySelector('button').disabled = false;
    });
  </script>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
