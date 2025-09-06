# simple-rag-beginner

## Project Overview

This is a Go prototype for a Retrieval-Augmented Generation (RAG) chatbot. It includes:

- **REST API**: POST `/query` endpoint accepts `{ "query": "your question" }` and returns `{ "response": "generated text" }`.
- **In-memory vector DB stub**: Simple map-based retrieval for context.
- **AI stub**: Returns canned responses combined with retrieved context.
- **CLI**: Command-line interface to interact with the backend.
- **Web Frontend**: Simple chat UI served at `/` for browser-based interaction.

## Project Structure

- `cmd/main.go` — Main entry point (API server & CLI)
- `internal/api/` — REST API and static file serving
- `internal/rag/` — RAG logic
- `internal/vectordb/` — In-memory vector DB stub
- `internal/model/` — AI stub
- `web/index.html` — Chat UI

## How to Run & Test

### Prerequisites
- Go 1.20 or newer

### 1. Run the REST API Server

```
go run ./cmd/main.go
```
Server starts at `http://localhost:8080`.

### 2. Test via Web Frontend

Open your browser and go to:

```
http://localhost:8080
```
Type your message in the chat screen and see responses.

### 3. Test via CLI

```
go run ./cmd/main.go cli
```
Type your queries in the terminal. Type `exit` to quit.

### 4. Test via API (e.g., curl)

```
curl -X POST -H "Content-Type: application/json" \
	-d '{"query": "hello"}' \
	http://localhost:8080/query
```

## Notes
- No Docker required.
- Logging and error handling are included.
- The vector DB and AI are stubs for demonstration.