

# RAG Microservices Pipeline

A production-ready Retrieval-Augmented Generation (RAG) system built with Go microservices and Docker Compose.

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- OpenAI API key

### 1. Setup Environment
```bash
export OPENAI_API_KEY="your-api-key-here"
```

### 2. Start Services
```bash
make up
```

### 3. Access the System
- **Web UI**: http://localhost:8080
- **RAG API**: http://localhost:8080/query
- **Direct LLM**: http://localhost:8080/query-direct
- **Qdrant Dashboard**: http://localhost:6333/dashboard

## 🏗️ Architecture

```
Web Scraping → Text Processing → Vector Embeddings → Query & Response
    ↓              ↓                 ↓                    ↓
  Crawler      →  Parser        →  Embedder         →  RAG API
    ↓              ↓                 ↓                    ↓
 Colly            goquery        OpenAI API         Gin + OpenAI
                                     ↓
                                  Qdrant Vector DB
```

### Services
- **Crawler** (`:8081`): Web scraping with Colly
- **Parser** (`:8082`): HTML parsing and text chunking  
- **Embedder** (`:8083`): OpenAI embeddings + Qdrant storage
- **RAG API** (`:8080`): Query interface with web UI
- **Qdrant** (`:6333`): Vector database

## 🧪 Testing

### Web Interface
Visit http://localhost:8080 and toggle between:
- **RAG Mode**: Uses vector database retrieval
- **Direct Mode**: Pure LLM responses

### API Testing
```bash
# Test RAG with vector retrieval
curl -X POST http://localhost:8080/query 
  -H "Content-Type: application/json" 
  -d '{"query": "What is Tekton?"}'

# Test direct LLM (no vector DB)
curl -X POST http://localhost:8080/query-direct 
  -H "Content-Type: application/json" 
  -d '{"query": "What is Tekton?"}'

# Test crawl + query pipeline
curl -X POST http://localhost:8080/crawl-and-query 
  -H "Content-Type: application/json" 
  -d '{"url": "https://tekton.dev", "query": "What is Tekton?", "collection": "my-docs"}'
```

### Health Checks
```bash
make health
```

## 📋 Available Commands

```bash
make help           # Show all commands
make up             # Start all services  
make down           # Stop all services
make logs           # View all logs
make test           # Test both RAG modes
make restart-api    # Restart just the API
make clean          # Clean up resources
```

## 🔧 Configuration

### Environment Variables
- `OPENAI_API_KEY`: Required for embeddings and LLM
- `CHUNK_SIZE`: Text chunk size (default: 2000)
- `VECTOR_DIM`: Embedding dimensions (default: 1536)

### Ports
- `8080`: RAG API (main interface)
- `8081`: Crawler service
- `8082`: Parser service  
- `8083`: Embedder service
- `6333`: Qdrant HTTP API
- `6334`: Qdrant gRPC API

## 🎯 Use Cases

1. **RAG Comparison**: Compare responses with/without vector retrieval
2. **Document Q&A**: Crawl websites and ask questions about content
3. **Knowledge Base**: Build searchable knowledge from web content
4. **Research Tool**: Extract insights from multiple web sources

## 🚀 Kubernetes Deployment

For production Kubernetes deployment:
```bash
make k8s-deploy
```

## 📁 Project Structure

```
cmd/                    # Microservices
├── crawler/           # Web scraping service
├── parser/            # Text processing service  
├── embedder/          # Vector embedding service
└── rag-api/           # Main API + web interface (inline HTML)
deploy/k8s/            # Kubernetes manifests
docker-compose.yml     # Local development
Makefile              # Build & run commands
```

### Local Development

1. **Docker Compose:**
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   make dev-up
   ```

2. **Test locally:**
   ```bash
   # Crawl a website
   curl -X POST http://localhost:8081/crawl -H "Content-Type: application/json" -d '{"url": "https://example.com"}'
   
   # Query the RAG system
   curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{"query": "What is this about?"}'
   ```

---

## Pipeline Flow

1. **Crawl** → POST `/crawl` with `{"url": "..."}`
2. **Parse** → Automatically triggered, chunks text  
3. **Embed** → Automatically triggered, generates embeddings and stores in Qdrant
4. **Query** → POST `/query` with `{"query": "..."}` returns AI response with sources

---

## Management

- **View status:** `make k8s-status`
- **View logs:** `make k8s-logs`  
- **Cleanup:** `make k8s-undeploy`
- **Build only:** `make build`

---

## Configuration

- Qdrant: 50Gi PVC, cosine similarity
- Embeddings: OpenAI Ada v2 (1536 dimensions)
- Chunking: ~2000 chars, sentence-aware
- LLM: GPT-3.5-turbo for responses

---

## Files

- `cmd/*/`: Microservice source code
- `deploy/k8s/`: Kubernetes manifests
- `Makefile`: Build and deployment automation
- `docker-compose.yml`: Local development stack

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