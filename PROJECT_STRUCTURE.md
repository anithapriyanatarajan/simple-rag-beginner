# Project Structure

> **Last Updated**: September 6, 2025  
> **Total Lines of Code**: 489 lines  
> **Languages**: Go, HTML, Bash  

## 📁 Directory Layout

```
simple-rag-beginner/
├── cmd/
│   └── main.go                     # 64 lines - Application entry point
├── internal/
│   ├── api/
│   │   ├── router.go               # 40 lines - REST API endpoints
│   │   └── static.go               # 9 lines - Static file serving
│   ├── embedding/
│   │   └── embedding.go            # 10 lines - Text-to-vector conversion
│   ├── model/
│   │   └── model.go                # 5 lines - AI response generation stub
│   ├── rag/
│   │   └── rag.go                  # 37 lines - RAG orchestration logic
│   └── vectordb/
│       └── qdrant.go               # 122 lines - Qdrant vector DB client
├── scripts/
│   └── setup-qdrant.sh             # 138 lines - Automated Qdrant deployment
├── web/
│   └── index.html                  # 64 lines - Interactive chat UI
├── go.mod                          # Go module definition
├── go.sum                          # Go dependency checksums
├── README.md                       # Main project documentation
└── PROJECT_STRUCTURE.md            # This file
```

## 📊 Component Statistics

| Component | Files | Lines | Purpose |
|-----------|-------|-------|---------|
| **Entry Point** | 1 | 64 | Application startup & CLI |
| **API Layer** | 2 | 49 | REST endpoints & static serving |
| **RAG Core** | 1 | 37 | Retrieval-Augmented Generation logic |
| **Vector DB** | 1 | 122 | Qdrant client & vector operations |
| **AI Model** | 1 | 5 | Response generation stub |
| **Embedding** | 1 | 10 | Text vectorization stub |
| **Frontend** | 1 | 64 | Interactive web chat interface |
| **Infrastructure** | 1 | 138 | Automated deployment script |
| **Configuration** | 2 | - | Go modules & dependencies |
| **Documentation** | 2 | - | README & this structure doc |

## 🏗️ Architecture Overview

### Data Flow
```
User Query → API Router → RAG Logic → Vector Search → Context Retrieval → AI Response
     ↓             ↓           ↓            ↓              ↓              ↓
Web/CLI → router.go → rag.go → qdrant.go → embedding.go → model.go → JSON Response
```

### Component Dependencies
```
cmd/main.go
├── internal/api/router.go
│   ├── internal/rag/rag.go
│   │   ├── internal/vectordb/qdrant.go
│   │   │   └── internal/embedding/embedding.go
│   │   └── internal/model/model.go
│   └── internal/api/static.go
│       └── web/index.html
└── scripts/setup-qdrant.sh (external)
```

## 📝 File Descriptions

### Core Application (`cmd/`)
- **`main.go`**: Application entry point with CLI argument parsing, server startup, and mode selection (API vs CLI)

### API Layer (`internal/api/`)
- **`router.go`**: HTTP request handlers, JSON parsing/generation, and REST endpoint definitions
- **`static.go`**: Static file serving for the web frontend

### RAG Logic (`internal/rag/`)
- **`rag.go`**: Orchestrates the retrieval-augmented generation process, combining vector search with AI responses

### Vector Database (`internal/vectordb/`)
- **`qdrant.go`**: Qdrant client implementation, collection management, vector operations, and connection handling

### AI Components (`internal/model/`, `internal/embedding/`)
- **`model.go`**: AI response generation stub (placeholder for LLM integration)
- **`embedding.go`**: Text-to-vector conversion stub (placeholder for real embedding models)

### Frontend (`web/`)
- **`index.html`**: Interactive chat interface with JavaScript for API communication

### Infrastructure (`scripts/`)
- **`setup-qdrant.sh`**: Automated Qdrant deployment with Docker, including fallback mechanisms and error handling

### Configuration
- **`go.mod`**: Go module definition with dependencies
- **`go.sum`**: Dependency checksums for reproducible builds

## 🔧 Key Design Patterns

### 1. **Clean Architecture**
- Clear separation of concerns between layers
- Internal packages prevent external access to implementation details
- Dependency injection through function parameters

### 2. **Error Handling**
- Graceful fallbacks when external services are unavailable
- Comprehensive error logging and user feedback
- Automatic retry mechanisms in deployment scripts

### 3. **Configuration Management**
- Environment-based configuration
- Sensible defaults with override capabilities
- Centralized configuration constants

### 4. **Modular Design**
- Each component has a single responsibility
- Easy to swap implementations (e.g., vector DB, AI model)
- Clear interfaces between components

## 🚀 Deployment Architecture

### Development Setup
```
Developer Machine
├── Go Application (Port 8080)
│   ├── REST API endpoints
│   └── Static web serving
└── Docker Container (Qdrant)
    ├── HTTP API (Port 6333)
    ├── gRPC API (Port 6334)
    └── Web Dashboard (Port 6333/dashboard)
```

### Production Considerations
```
Load Balancer
├── Multiple Go App Instances
├── Qdrant Cluster
├── Real Embedding Service
└── LLM API Integration
```

## 📈 Growth Path

### Phase 1: Current State ✅
- Basic RAG architecture
- Stub implementations
- Docker deployment
- Web interface

### Phase 2: Production Ready 🔄
- Real embedding models
- LLM API integration
- Authentication
- Rate limiting

### Phase 3: Scale 📈
- Multiple vector databases
- Distributed deployment
- Advanced retrieval strategies
- Performance optimization

## 🔍 Code Quality Metrics

| Metric | Value | Status |
|--------|-------|--------|
| **Total Files** | 12 | ✅ Manageable |
| **Average File Size** | 40 lines | ✅ Focused |
| **Dependency Count** | 2 external | ✅ Minimal |
| **Test Coverage** | 0% | ⚠️ Needs tests |
| **Documentation** | 100% | ✅ Complete |

## 🛠️ Maintenance Notes

### Regular Updates Needed
- [ ] Dependency versions in `go.mod`
- [ ] Line counts in this document
- [ ] Architecture diagrams as system evolves
- [ ] Performance benchmarks

### Code Quality Checklist
- [ ] Add unit tests for each component
- [ ] Add integration tests for API endpoints
- [ ] Implement proper logging levels
- [ ] Add performance monitoring
- [ ] Create CI/CD pipeline

---

> 💡 **Tip**: Update this document whenever you add new files or significantly modify the architecture. This helps maintain project clarity as the codebase grows.
