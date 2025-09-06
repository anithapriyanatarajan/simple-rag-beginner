package rag

import (
	"log"
	"simple-rag-beginner/internal/model"
	"simple-rag-beginner/internal/vectordb"
	"strings"
)

// GenerateResponseWithContext returns both the AI response and the top-K context
func GenerateResponseWithContext(query string) (string, []string) {
	topK := 3
	results, err := vectordb.RetrieveTopK(query, topK)
	var context string
	if err != nil {
		context = "Error retrieving context: " + err.Error()
		results = []string{context}
	} else {
		context = strings.Join(results, " | ")
	}
	response := model.GenerateText(query, context)
	return response, results
}

func GenerateResponse(query string) string {
	topK := 3
	results, err := vectordb.RetrieveTopK(query, topK)
	var context string
	if err != nil {
		context = "Error retrieving context: " + err.Error()
	} else {
		context = strings.Join(results, " | ")
	}
	response := model.GenerateText(query, context)
	log.Printf("Query: %s | Context: %s | Response: %s", query, context, response)
	return response
}
