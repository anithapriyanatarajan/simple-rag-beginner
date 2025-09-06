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
)

func main() {
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
		var result map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Printf("Decode error: %v", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		fmt.Printf("Response: %s\n", result["response"])
	}
}

func setupRouter() http.Handler {
	return api.NewRouter()
}
