package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"simple-rag-beginner/internal/api"
	"simple-rag-beginner/internal/vectordb"
)

func main() {
	// Initialize vector DB
	if err := vectordb.InitQdrant("localhost:6333"); err != nil {
		log.Printf("Warning: Failed to initialize Qdrant: %v", err)
	}
	defer vectordb.CloseQdrant()

	if len(os.Args) > 1 && os.Args[1] == "cli" {
		runCLI()
		return
	}
	log.Println("Starting RAG Chatbot REST API on :8080...")
	if err := http.ListenAndServe(":8080", setupRouter()); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func runCLI() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter your query (or 'exit'): ")
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]
		if input == "exit" {
			break
		}
		payload := map[string]string{"query": input}
		body, _ := json.Marshal(payload)
		resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("Request error: %v", err)
			continue
		}
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("Decode error: %v", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		fmt.Printf("Response: %s\n", result["response"])
		if ctx, ok := result["context"]; ok {
			fmt.Printf("Retrieved context: %v\n", ctx)
		}
	}
}

func setupRouter() http.Handler {
	return api.NewRouter()
}
