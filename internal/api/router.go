package api

import (
	"encoding/json"
	"log"
	"net/http"
	"simple-rag-beginner/internal/model"
	"simple-rag-beginner/internal/rag"
)

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Response string   `json:"response"`
	Context  []string `json:"context"`
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/query", handleQuery)
	mux.HandleFunc("/model/info", handleModelInfo)
	// Serve static files at root
	mux.Handle("/", StaticHandler())
	return mux
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response, context := rag.GenerateResponseWithContext(req.Query)
	json.NewEncoder(w).Encode(QueryResponse{Response: response, Context: context})
}

func handleModelInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := model.GetModelInfo()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
