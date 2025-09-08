package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "parser_requests_total"},
		[]string{"status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "parser_request_duration_seconds"},
		[]string{"endpoint"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter, requestDuration)
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	log.Println(`{"level":"info","msg":"parser service started","port":8080}`)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/parse", parseHandler)
	http.HandleFunc("/health", healthHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		requestDuration.WithLabelValues("parse").Observe(time.Since(start).Seconds())
	}()

	var req struct {
		URL        string `json:"url"`
		HTML       string `json:"html"`
		Title      string `json:"title"`
		Collection string `json:"collection"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		requestCounter.WithLabelValues("error").Inc()
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Received parse request. URL: %s, Collection: %s", req.URL, req.Collection)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(req.HTML))
	if err != nil {
		requestCounter.WithLabelValues("error").Inc()
		http.Error(w, "Failed to parse HTML", http.StatusInternalServerError)
		return
	}

	// Extract main content, skip navigation/footer
	doc.Find("nav, footer, script, style").Remove()
	text := strings.TrimSpace(doc.Text())

	chunkSizeStr := os.Getenv("CHUNK_SIZE")
	chunkSize := 2000
	if chunkSizeStr != "" {
		if size, err := strconv.Atoi(chunkSizeStr); err == nil {
			chunkSize = size
		}
	}

	chunks := chunkText(text, chunkSize)

	// Create chunks with UUID-based IDs for Qdrant compatibility
	var enrichedChunks []map[string]interface{}
	for i, chunk := range chunks {
		chunkID := uuid.New().String()
		enrichedChunks = append(enrichedChunks, map[string]interface{}{
			"id":      chunkID,
			"url":     req.URL,
			"title":   req.Title,
			"content": chunk,
			"index":   i,
		})
	}

	// Forward to embedder
	embedderURL := os.Getenv("EMBEDDER_URL")
	if embedderURL == "" {
		embedderURL = "http://embedder:8080"
	}

	collection := req.Collection
	if collection == "" {
		collection = "documents"
	}

	payload := map[string]interface{}{
		"chunks":     enrichedChunks,
		"collection": collection,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Forward to embedder with proper error handling
	go func() {
		log.Printf("Forwarding to embedder. Collection: %s, Chunks: %d", collection, len(enrichedChunks))
		resp, err := http.Post(embedderURL+"/embed", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Failed to forward to embedder: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Embedder returned error status: %d", resp.StatusCode)
		} else {
			log.Printf("Successfully forwarded to embedder. Collection: %s", collection)
		}
	}()

	requestCounter.WithLabelValues("success").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"chunks": enrichedChunks,
		"count":  len(enrichedChunks),
	})
}

func chunkText(text string, size int) []string {
	// Split by sentences first, then by size if needed
	sentences := strings.Split(text, ". ")
	var chunks []string
	var currentChunk strings.Builder

	for _, sentence := range sentences {
		if currentChunk.Len()+len(sentence) > size && currentChunk.Len() > 0 {
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()
		}
		if currentChunk.Len() > 0 {
			currentChunk.WriteString(". ")
		}
		currentChunk.WriteString(sentence)
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
