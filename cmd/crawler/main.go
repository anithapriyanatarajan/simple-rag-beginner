package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gocolly/colly"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "crawler_requests_total"},
		[]string{"status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "crawler_request_duration_seconds"},
		[]string{"endpoint"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter, requestDuration)
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	log.Println(`{"level":"info","msg":"crawler service started","port":8080}`)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/crawl", crawlHandler)
	http.HandleFunc("/health", healthHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func crawlHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		requestDuration.WithLabelValues("crawl").Observe(time.Since(start).Seconds())
	}()

	var req struct {
		URL        string `json:"url"`
		Collection string `json:"collection"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		requestCounter.WithLabelValues("error").Inc()
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	c := colly.NewCollector()
	var htmlContent string
	var title string

	c.OnHTML("html", func(e *colly.HTMLElement) {
		htmlContent, _ = e.DOM.Html()
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		title = e.Text
	})

	if err := c.Visit(req.URL); err != nil {
		requestCounter.WithLabelValues("error").Inc()
		http.Error(w, fmt.Sprintf("Failed to crawl: %v", err), http.StatusInternalServerError)
		return
	}

	// Forward to parser
	collection := req.Collection
	if collection == "" {
		collection = "documents"
	}

	response := map[string]string{
		"url":        req.URL,
		"html":       htmlContent,
		"title":      title,
		"collection": collection,
	}

	parserURL := os.Getenv("PARSER_URL")
	if parserURL == "" {
		parserURL = "http://parser:8080"
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Forward to parser with proper error handling
	go func() {
		resp, err := http.Post(parserURL+"/parse", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Failed to forward to parser: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Parser returned error status: %d", resp.StatusCode)
		} else {
			log.Printf("Successfully forwarded to parser. Collection: %s", collection)
		}
	}()

	requestCounter.WithLabelValues("success").Inc()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
