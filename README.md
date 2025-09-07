

# Simple RAG Beginner

A minimal Go-based RAG chatbot using Qdrant for vector search.

---

## Install & Run

### VM Mode (No Docker)
1. Install Go 1.20+
2. Download and run Qdrant manually ([Qdrant docs](https://qdrant.tech/documentation/quick-start/))
3. Start API server:
   ```bash
   go run cmd/main.go
   ```
4. Open [http://localhost:8080](http://localhost:8080) in your browser

### Docker Mode (Recommended)
1. Install Docker
2. Run Qdrant setup script:
   ```bash
   ./scripts/setup-qdrant.sh
   ```
3. Start API server:
   ```bash
   go run cmd/main.go
   ```
4. Open [http://localhost:8080](http://localhost:8080)

---

## Test

- **Web UI:** Use browser at [http://localhost:8080](http://localhost:8080)
- **API:**
  ```bash
  curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{"query": "hello"}'
  ```
- **CLI:**
  ```bash
  go run cmd/main.go cli
  ```

---

## Troubleshoot & Cleanup

- **Stop all & cleanup:**
  ```bash
  ./scripts/cleanup.sh
  ```
- **Restart Qdrant only:**
  ```bash
  docker restart qdrant-rag
  ```
- **View Qdrant logs:**
  ```bash
  docker logs qdrant-rag
  ```
- **Check port usage:**
  ```bash
  lsof -i :8080
  ```

---

## Notes

- Default ports: 8080 (API), 6333 (Qdrant HTTP)
- Embedding/model logic is a stub; replace for production
- See `PROJECT_STRUCTURE.md` for file details

### 2. Start the RAG Application

**API Mode (Default):**
```bash
go run cmd/main.go
```

**Expected Output:**
```
2025/09/06 14:08:29 Starting RAG Chatbot REST API on :8080...
```

The application automatically:
- Connects to Qdrant (localhost:6334)
- Creates the `rag_collection` with proper vector configuration
- Inserts sample documents with embeddings
- Starts HTTP server on port 8080

### 3. Test the System

**Option A: Web Interface**
Open your browser: http://localhost:8080

**Option B: REST API**
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "hello world"}'
```

**Expected Response:**
```json
{
  "response": "[AI Stub] You asked: 'hello world'. Context: Hi there! This is context for 'hello'.",
  "context": ["Hi there! This is context for 'hello'."]
}
```

**Option C: CLI Mode (Interactive Terminal)**
```bash
# First: Start the API server (in one terminal)
go run cmd/main.go

# Then: Start CLI mode (in another terminal)  
go run cmd/main.go cli
# Follow interactive prompts
```

**CLI Mode Example Session:**
```
$ go run cmd/main.go cli
Enter your query (or 'exit'): who am i?
Response: [AI Stub] You asked: 'who am i?'. Context: You are a helpful AI assistant...
Retrieved context: [You are a helpful AI assistant designed to answer questions...]

Enter your query (or 'exit'): weather today
Response: [AI Stub] You asked: 'weather today'. Context: Today's weather is sunny...
Retrieved context: [Today's weather is sunny and bright. Perfect day for outdoor activities...]

Enter your query (or 'exit'): exit
```

> **Note**: CLI mode requires the API server to be running simultaneously as it makes HTTP requests to `localhost:8080` internally.

## 🔧 System Management

### Application Modes

The RAG chatbot supports two distinct operation modes:

**1. API Mode (Default)**
- Starts REST API server on port 8080
- Serves web interface at `http://localhost:8080`
- Accepts HTTP POST requests to `/query` endpoint
- Usage: `go run cmd/main.go`

**2. CLI Mode (Interactive)**
- Provides terminal-based chat interface
- Requires API server running simultaneously
- Interactive prompts for real-time queries
- Usage: `go run cmd/main.go cli`

### Quick Cleanup

**🧹 Stop Everything (One Command):**
```bash
./scripts/cleanup.sh
```

This automated cleanup script:
- ✅ Stops all Go applications (port 8080)
- ✅ Stops and removes Qdrant containers
- ✅ Frees up all project ports (8080, 6333, 6334)
- ✅ Provides clear status feedback
- ✅ Shows restart instructions

### Qdrant Database

**Container Management:**
```bash
# Check status
docker ps | grep qdrant-rag

# View real-time logs  
docker logs qdrant-rag -f

# Stop/Start
docker stop qdrant-rag
docker start qdrant-rag

# Complete removal and recreation
docker rm -f qdrant-rag
./scripts/setup-qdrant.sh
```

**Qdrant Web Dashboard:**
- URL: http://localhost:6333/dashboard
- Features: Browse collections, execute queries, monitor cluster status

**Data Persistence:**
- **Preferred**: Local directory `./qdrant_storage/` (when filesystem permits)
- **Fallback**: Docker internal volume (automatic fallback)
- **Backup**: `cp -r qdrant_storage qdrant_storage.backup` (if using local storage)

### Application Management

**Start/Stop Application:**
```bash
# API mode (default - no arguments needed)
go run cmd/main.go

# CLI interactive mode (requires API server running)
go run cmd/main.go cli

# Stop: Ctrl+C or kill the process
```

**Port Usage:**
- **8080**: Go application (REST API + Web UI)
- **6333**: Qdrant HTTP API + Web Dashboard
- **6334**: Qdrant gRPC API (internal use)

## 📊 API Reference

### POST /query
**Request:**
```json
{
  "query": "your question here"
}
```

**Response:**
```json
{
  "response": "AI-generated response with context",
  "context": ["retrieved context 1", "retrieved context 2"]
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/query 
  -H "Content-Type: application/json" 
  -d '{"query": "what is an agent?"}'
```

### GET /
Serves the interactive web chat interface.

## 🖥️ CLI Mode Usage

The RAG chatbot provides an interactive command-line interface for direct terminal usage.

### How to Use CLI Mode

**Step 1: Start the API Server**
```bash
go run cmd/main.go
```
Keep this running - it provides the backend services.

**Step 2: Start CLI Mode (in a new terminal)**
```bash
go run cmd/main.go cli
```

### CLI Features

- **Interactive Prompts**: Ask questions directly in the terminal
- **Real-time Responses**: Get AI responses with retrieved context
- **Context Display**: See exactly what documents were retrieved
- **Easy Exit**: Type `exit` to quit the CLI session
- **Multiple Sessions**: Run multiple CLI instances simultaneously

### Example CLI Session

```
$ go run cmd/main.go cli
2025/09/06 14:33:18 Collection might already exist: CreateCollection() failed: rag_collection: rpc error: code = AlreadyExists desc = Wrong input: Collection `rag_collection` already exists!

Enter your query (or 'exit'): hello there
Response: [AI Stub] You asked: 'hello there'. Context: Hi there! This is context for 'hello'. | You are a helpful AI assistant designed to answer questions and provide information. | Today's weather is sunny and bright. Perfect day for outdoor activities.
Retrieved context: [Hi there! This is context for 'hello'. You are a helpful AI assistant designed to answer questions and provide information. Today's weather is sunny and bright. Perfect day for outdoor activities.]

Enter your query (or 'exit'): what is an agent?
Response: [AI Stub] You asked: 'what is an agent?'. Context: Agents are autonomous entities that can perceive their environment and act upon it. | Agents are autonomous entities that can perceive their environment and act upon it. | Agents are autonomous entities that can perceive their environment and act upon it.
Retrieved context: [Agents are autonomous entities that can perceive their environment and act upon it. Agents are autonomous entities that can perceive their environment and act upon it. Agents are autonomous entities that can perceive their environment and act upon it.]

Enter your query (or 'exit'): exit
$
```

### CLI vs Web Interface

| Feature | CLI Mode | Web Interface |
|---------|----------|---------------|
| **Access** | Terminal required | Browser required |
| **Setup** | Two terminals needed | Single server |
| **Context Display** | Full raw context | Formatted display |
| **Session** | Single-user | Multi-user capable |
| **History** | Terminal scroll-back | Chat history in UI |
| **Best For** | Development/testing | End-user interaction |

## 🧠 Technical Details

### Vector Database Configuration
- **Collection**: `rag_collection`
- **Vector Dimensions**: 8 (simple embedding stub)
- **Distance Metric**: Cosine similarity
- **Top-K Results**: 3 (configurable)

### Sample Data
The system pre-loads these documents:
- "Hi there! This is context for 'hello'."
- "Today's weather is sunny."  
- "Agents are autonomous entities."

### Embedding Strategy
Current implementation uses a simple 8-dimensional embedding stub. For production:
- Replace with real embedding models (OpenAI, HuggingFace, etc.)
- Increase vector dimensions (typically 384, 768, or 1536)
- Add proper text preprocessing

### AI Response Generation
Current implementation returns formatted responses. For production:
- Integrate with LLM APIs (OpenAI GPT, Anthropic Claude, etc.)
- Add proper prompt engineering
- Implement response streaming

## 🚧 Development Status

**✅ Completed:**
- Basic RAG architecture
- Qdrant vector database integration
- REST API with JSON responses
- Web chat interface
- CLI interaction mode
- Automated deployment script
- Error handling and fallbacks

**🔄 Next Steps:**
- Real embedding model integration
- LLM API integration  
- Authentication and rate limiting
- Production deployment configuration
- Advanced vector search options
- Document ingestion pipeline

## 🛠️ Troubleshooting

**🧹 Quick Fix - Clean Restart:**
```bash
# Stop everything and start fresh
./scripts/cleanup.sh
./scripts/setup-qdrant.sh
go run cmd/main.go
```

**Port Conflicts:**
```bash
# Check what's using port 8080
lsof -i :8080

# Kill conflicting process  
kill <PID>

# Or use the cleanup script
./scripts/cleanup.sh
```

**Qdrant Connection Issues:**
```bash
# Check Qdrant status
curl http://localhost:6333/collections

# View container logs
docker logs qdrant-rag

# Restart Qdrant
docker restart qdrant-rag
```

**Permission Issues:**
The setup script automatically handles filesystem permission issues by falling back to Docker internal volumes.

## 📜 License

This is a prototype for learning purposes. Adapt as needed for your projects.

## 📁 Project Documentation

- **[README.md](README.md)** - Main project documentation (this file)
- **[PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)** - Detailed architecture and file organization

---

**Ready to enhance your RAG system?** 🚀
- Add real embeddings (OpenAI, HuggingFace)
- Integrate LLM APIs (GPT, Claude, Llama)  
- Scale with production vector databases
- Deploy with Docker Compose or Kubernetes```
Server starts at `http://localhost:8080` and automatically:
- Connects to Qdrant
- Creates the "rag_collection" 
- Inserts sample data with embeddings

### 3. Test via Web Frontend

Open your browser and go to:

```
http://localhost:8080
```
Type your message in the chat screen and see responses with retrieved context.

### 4. Test via API (e.g., curl)

```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"query": "hello"}' \
     http://localhost:8080/query
```

Expected response:
```json
{
  "response": "[AI Stub] You asked: 'hello'. Context: Hi there! This is context for 'hello'.",
  "context": ["Hi there! This is context for 'hello'."]
}
```

### 5. Stop Services

**Stop Qdrant:**
```bash
docker stop qdrant-rag
```

**Restart Qdrant:**
```bash
docker start qdrant-rag
```

**View Qdrant logs:**
```bash
docker logs qdrant-rag
```

**Remove Qdrant completely:**
```bash
docker stop qdrant-rag
docker rm qdrant-rag
```

## Qdrant Management

### Docker CLI Commands

The setup script creates a container named `qdrant-rag`. You can manage it using standard Docker commands:

```bash
# Check status
docker ps | grep qdrant

# View logs
docker logs qdrant-rag -f

# Stop/Start
docker stop qdrant-rag
docker start qdrant-rag

# Remove and recreate
docker rm -f qdrant-rag
./scripts/setup-qdrant.sh
```

### Qdrant Web Dashboard

Access the Qdrant web interface at: http://localhost:6333/dashboard

This allows you to:
- Browse collections and points
- Execute queries
- Monitor cluster status
- View metrics and logs

## Features

- **Real Vector Search**: Uses Qdrant for semantic similarity search
- **Automatic Collection Setup**: Creates collection with proper vector configuration
- **Embedding Integration**: Text-to-vector conversion for queries and documents
- **Top-K Retrieval**: Configurable number of similar results
- **Logging**: Request/response logging for debugging
- **Error Handling**: Graceful error handling with fallbacks
- **Web UI**: Interactive chat interface
- **CLI**: Command-line interface for testing

## Development Notes

- **Vector Dimensions**: Currently using 8-dimensional embeddings (simple stub)
- **Distance Metric**: Cosine similarity for vector search
- **Collection Name**: "rag_collection" (configurable)
- **Qdrant Ports**: HTTP 6333, gRPC 6334
- **Fallback**: Graceful fallback if Qdrant is unavailable

## Configuration

### Qdrant Connection
The application connects to Qdrant at `localhost:6334` (gRPC). To use a different host:

```go
// In vectordb/qdrant.go InitQdrant function
client, err = qdrant.NewClient(&qdrant.Config{
    Host: "your-qdrant-host",
    Port: 6334,
})
```

### Vector Dimensions
To change embedding dimensions, update:
1. `internal/embedding/embedding.go` - TextToVector function
2. `internal/vectordb/qdrant.go` - VectorParams Size

## Future Enhancements

- [ ] Real embedding models (OpenAI, Hugging Face)
- [ ] LLM integration (OpenAI, Anthropic, etc.)
- [ ] Collection management UI
- [ ] Authentication and rate limiting
- [ ] Batch processing for large datasets
- [ ] Multiple collection support