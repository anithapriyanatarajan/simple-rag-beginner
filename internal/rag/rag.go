package rag

import (
	"log"
	"simple-rag-beginner/internal/model"
	"simple-rag-beginner/internal/vectordb"
)

func GenerateResponse(query string) string {
	context := vectordb.Retrieve(query)
	response := model.GenerateText(query, context)
	log.Printf("Query: %s | Context: %s | Response: %s", query, context, response)
	return response
}
