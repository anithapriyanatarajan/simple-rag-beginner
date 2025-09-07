package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sashabaranov/go-openai"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "embedder_requests_total"},
		[]string{"status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "embedder_request_duration_seconds"},
		[]string{"endpoint"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter, requestDuration)
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	log.Println(`{"level":"info","msg":"embedder service started","port":8080}`)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/embed", embedHandler)
	http.HandleFunc("/health", healthHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func embedHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		requestDuration.WithLabelValues("embed").Observe(time.Since(start).Seconds())
	}()

	var req struct {
		Chunks     []map[string]interface{} `json:"chunks"`
		Collection string                   `json:"collection"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		requestCounter.WithLabelValues("error").Inc()
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		requestCounter.WithLabelValues("error").Inc()
		http.Error(w, "OpenAI API key not configured", http.StatusInternalServerError)
		return
	}

	client := openai.NewClient(apiKey)

	// Process chunks in batches
	batchSize := 100
	for i := 0; i < len(req.Chunks); i += batchSize {
		end := i + batchSize
		if end > len(req.Chunks) {
			end = len(req.Chunks)
		}

		batch := req.Chunks[i:end]
		collection := req.Collection
		if collection == "" {
			collection = "documents"
		}
		if err := processBatch(client, batch, collection); err != nil {
			log.Printf(`{"level":"error","msg":"Failed to process batch","error":"%v"}`, err)
			requestCounter.WithLabelValues("error").Inc()
			http.Error(w, fmt.Sprintf("Failed to process embeddings: %v", err), http.StatusInternalServerError)
			return
		}
	}

	requestCounter.WithLabelValues("success").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"count":  len(req.Chunks),
	})
}

func processBatch(client *openai.Client, chunks []map[string]interface{}, collection string) error {
	var texts []string
	for _, chunk := range chunks {
		if content, ok := chunk["content"].(string); ok {
			texts = append(texts, content)
		}
	}

	if len(texts) == 0 {
		return nil
	}

	// Get embeddings from OpenAI
	resp, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: texts,
		Model: openai.AdaEmbeddingV2,
	})
	if err != nil {
		return fmt.Errorf("OpenAI API error: %w", err)
	}

	// Upsert to Qdrant
	var points []map[string]interface{}
	for i, chunk := range chunks {
		if i < len(resp.Data) {
			points = append(points, map[string]interface{}{
				"id":      chunk["id"],
				"vector":  resp.Data[i].Embedding,
				"payload": chunk,
			})
		}
	}

	return upsertToQdrant(points, collection)
}

func upsertToQdrant(points []map[string]interface{}, collection string) error {
	qdrantURL := os.Getenv("QDRANT_URL")
	if qdrantURL == "" {
		qdrantURL = "http://qdrant:6333"
	}

	// Create collection if not exists
	collectionPayload := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     1536, // Ada v2 embedding size
			"distance": "Cosine",
		},
	}

	collectionData, err := json.Marshal(collectionPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal collection payload: %w", err)
	}

	req, err := http.NewRequest("PUT", qdrantURL+"/collections/"+collection, bytes.NewBuffer(collectionData))
	if err != nil {
		return fmt.Errorf("failed to create collection request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	defer resp.Body.Close()

	// Log collection creation result
	if resp.StatusCode != 200 && resp.StatusCode != 409 { // 409 = already exists
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("collection creation failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	log.Printf(`{"level":"info","msg":"Collection ensured","collection":"%s","status":%d}`, collection, resp.StatusCode)

	// Upsert points
	upsertPayload := map[string]interface{}{
		"points": points,
	}

	upsertData, err := json.Marshal(upsertPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal upsert payload: %w", err)
	}

	upsertReq, err := http.NewRequest("PUT", qdrantURL+"/collections/"+collection+"/points", bytes.NewBuffer(upsertData))
	if err != nil {
		return fmt.Errorf("failed to create upsert request: %w", err)
	}
	upsertReq.Header.Set("Content-Type", "application/json")

	upsertResp, err := client.Do(upsertReq)
	if err != nil {
		return fmt.Errorf("failed to upsert points: %w", err)
	}
	defer upsertResp.Body.Close()

	if upsertResp.StatusCode != 200 {
		body, _ := io.ReadAll(upsertResp.Body)
		return fmt.Errorf("upsert failed: status %d, body: %s", upsertResp.StatusCode, string(body))
	}

	log.Printf(`{"level":"info","msg":"Points upserted","collection":"%s","count":%d}`, collection, len(points))
	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
