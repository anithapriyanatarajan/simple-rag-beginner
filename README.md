# Simple RAG Beginner

A Go-based Retrieval-Augmented Generation (RAG) chatbot with **real AI integration** using Ollama for both text generation and embeddings.

## 📚 Documentation

- **[🏗️ Architecture Guide](ARCHITECTURE.md)** - Complete system design, component details, and technical specifications
- **[🐳 Docker Guide](DOCKER.md)** - Comprehensive Docker deployment, development setup, and troubleshooting
- **[📖 Main Guide](#-quick-start)** - Getting started, features, and basic usage (this file)

## ✨ Features

- **🤖 Real AI Models**: Powered by Ollama with LLaMA 3.2 and nomic-embed-text
- **🔍 Semantic Search**: High-quality 768D embeddings with Qdrant vector database
- **🌐 REST API**: JSON endpoints for query processing with context retrieval
- **💬 Web Chat UI**: Interactive browser-based chat interface
- **⚡ Auto Setup**: One-command deployment for Qdrant and Ollama
- **🛡️ Production Ready**: No stubs - real AI models required
- **📊 Monitoring**: Built-in logging and Qdrant web dashboard

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Frontend  │────│   Go REST API    │────│   Qdrant DB     │────│   Ollama AI     │
│   (Chat UI)     │    │   (RAG Logic)    │    │   (768D Vectors)│    │ (LLM + Embed)   │
└─────────────────┘    └──────────────────┘    └─────────────────┘    └─────────────────┘
```

📖 **Detailed Architecture**: See [ARCHITECTURE.md](ARCHITECTURE.md) for complete system design, component details, and data flow diagrams.

**Core Components:**
- `cmd/main.go` — Application entry point with Ollama integration
- `internal/api/` — REST endpoints and static file serving
- `internal/rag/` — RAG orchestration with real context retrieval
- `internal/vectordb/` — Qdrant client with dynamic vector dimensions
- `internal/embedding/` — Ollama embedding integration (768D vectors)
- `internal/model/` — Ollama LLM integration (LLaMA 3.2)
- `web/index.html` — Interactive chat interface
- `scripts/setup-qdrant.sh` — Automated Qdrant deployment
- `scripts/setup-ollama.sh` — Automated Ollama setup with models
- `scripts/cleanup.sh` — Complete project cleanup

> 📋 **Detailed Structure**: See [`PROJECT_STRUCTURE.md`](PROJECT_STRUCTURE.md) for comprehensive architecture documentation.

## 🚀 Quick Start

Choose your deployment method:

### 🐳 Docker Deployment (Recommended)
```bash
# One-command full deployment
./scripts/docker-deploy.sh
```
✅ **Includes everything**: Qdrant, Ollama, AI models, and the application  
✅ **Production ready**: Health checks, proper networking, persistent storage  
✅ **Cross-platform**: Works on Linux, macOS, and Windows  

**Access after deployment:**
- Web UI: http://localhost:8080  
- API: http://localhost:8080/query
- Qdrant Dashboard: http://localhost:6333/dashboard

📖 **Full Docker guide**: See [DOCKER.md](DOCKER.md)

### ⚡ Local Development

For local development or if you prefer running services individually:

### Prerequisites
- **Go 1.20+**
- **Docker** (for Qdrant vector database)
- **Ollama** (for AI models - auto-installed by setup script)

### 1. Setup AI Models (Ollama)

```bash
./scripts/setup-ollama.sh
```

This automated script:
- ✅ Installs Ollama if not present
- ✅ Downloads LLaMA 3.2 (generative model)
- ✅ Downloads nomic-embed-text (embedding model)
- ✅ Starts Ollama service
- ✅ Tests model availability

### 2. Setup Vector Database

```bash
./scripts/setup-qdrant.sh
```

This automated script:
This automated script:
- ✅ Pulls latest Qdrant Docker image
- ✅ Starts Qdrant container with optimal configuration
- ✅ Configures for 768D vector storage (Ollama embeddings)
- ✅ Validates Qdrant is responding correctly

**Output Example:**
```
🚀 Setting up Qdrant Vector Database...
✅ Qdrant is running and responding!
🌐 Web UI available at: http://localhost:6333/dashboard
📁 Data storage: Docker internal volume
```

### 3. Start the RAG Application

```bash
go run cmd/main.go
```

**Expected Output:**
```
2025/09/06 15:38:32 ✅ Ollama embedding connected successfully with model: nomic-embed-text:latest
2025/09/06 15:38:32 ✅ Ollama connected successfully with model: llama3.2:latest
2025/09/06 15:38:42 Using vector dimension: 768
2025/09/06 15:38:43 Starting RAG Chatbot REST API on :8080...
```

The application automatically:
- Connects to Ollama for AI models (localhost:11434)
- Connects to Qdrant for vector storage (localhost:6334)
- Creates collection with 768D vectors for real embeddings
- Inserts sample documents with semantic embeddings
- Starts HTTP server on port 8080

### 4. Test the System

**Option A: Web Interface**
Open your browser: http://localhost:8080

**Option B: REST API**
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "hello world"}'
```

**Expected Response (Real AI):**
```json
{
  "response": "Hello! How can I assist you today? I'm here to help with any questions or tasks you might have.",
  "context": ["Hi there! This is context for 'hello'. Welcome to our system!"]
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
Response: Based on the context provided, you are a helpful AI assistant designed to answer questions and provide information. You're here to assist users with their queries and help them find the information they need.
Retrieved context: [You are a helpful AI assistant designed to answer questions...]

Enter your query (or 'exit'): weather today
Response: According to the available information, today's weather is sunny and bright. It's described as a perfect day for outdoor activities, so it would be great for spending time outside!
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
Response: Hello! It's great to meet you. Based on the context, I'm here to help you as a helpful AI assistant designed to answer questions and provide information. The weather today is sunny and bright - it's described as a perfect day for outdoor activities. How can I assist you today?
Retrieved context: [Hi there! This is context for 'hello'. You are a helpful AI assistant designed to answer questions and provide information. Today's weather is sunny and bright. Perfect day for outdoor activities.]

Enter your query (or 'exit'): what is an agent?
Response: An agent is an autonomous entity that can perceive its environment and act upon it. These entities are designed to operate independently, gathering information about their surroundings and making decisions or taking actions based on that information. Agents can be software-based (like AI agents) or physical entities (like robots) that interact with their environment to achieve specific goals.
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
- **Vector Dimensions**: 768 (real semantic embeddings via Ollama)
- **Distance Metric**: Cosine similarity  
- **Top-K Results**: 3 (configurable)
- **Dynamic Dimension Support**: Automatically detects vector dimensions

### Sample Data
The system pre-loads these documents:
- "Hi there! This is context for 'hello'."
- "Today's weather is sunny."  
- "Agents are autonomous entities."

### Embedding Strategy
Real semantic embeddings powered by Ollama:
- **Model**: `nomic-embed-text:latest` (768-dimensional vectors)
- **API Integration**: HTTP REST API with Ollama server
- **Quality**: Production-grade semantic similarity matching
- **Processing**: Automatic text normalization and embedding generation

### AI Response Generation
Full LLM integration via Ollama:
- **Model**: `llama3.2:latest` for natural language generation
- **Context Integration**: Retrieved documents injected into prompts
- **Response Quality**: Real AI reasoning and natural language output
- **API**: HTTP REST interface with Ollama server

## 🚧 Development Status

**✅ Completed:**
- Production RAG architecture with Ollama integration
- Real semantic embeddings (768D via nomic-embed-text)
- LLM text generation (via llama3.2)
- Qdrant vector database integration with dynamic dimensions
- REST API with JSON responses
- Web chat interface with real AI responses
- CLI interaction mode
- Automated deployment script with Ollama setup
- Error handling requiring AI services

**🔄 Next Steps:**
- Authentication and rate limiting
- Production deployment configuration  
- Advanced vector search options
- Document ingestion pipeline
- Multi-model support and configuration
- Response streaming and async processing

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
  "response": "Hello! It's nice to meet you. Based on the context, I can see there's a greeting here for you. How can I help you today?",
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

- **Vector Dimensions**: Using 768-dimensional real semantic embeddings via Ollama
- **Distance Metric**: Cosine similarity for vector search  
- **Collection Name**: "rag_collection" (configurable)
- **Qdrant Ports**: HTTP 6333, gRPC 6334
- **Requirements**: Ollama server must be running with required models

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