# Docker Deployment Guide

This guide covers deploying the RAG Chatbot using Docker and Docker Compose.

## 🐳 Quick Start with Docker

### Prerequisites
- Docker (v20.0+)
- Docker Compose (v2.0+)
- NVIDIA Docker runtime (for GPU acceleration - optional)

### One-Command Deployment
```bash
# Deploy the complete stack (auto-detects CPU/GPU)
./scripts/docker-deploy.sh
```

**🎮 GPU vs 💻 CPU Mode:**
- **GPU Mode**: Automatically enabled if NVIDIA GPU + Docker GPU support detected
- **CPU Mode**: Default fallback - works on any system (slower but functional)
- **Performance**: GPU is ~5-10x faster for AI inference

This script will:
1. Detect GPU support and configure accordingly
2. Pull and start Qdrant vector database
3. Pull and start Ollama AI server (CPU or GPU mode)
3. Download required AI models (LLaMA 3.2 and nomic-embed-text)
4. Build and start the RAG application
5. Verify all services are healthy

## 📋 Service URLs
After deployment, access these services:

- **RAG Chatbot**: http://localhost:8080
- **Qdrant Dashboard**: http://localhost:6333/dashboard  
- **Ollama API**: http://localhost:11434
- **Health Check**: http://localhost:8080/health

## 🔧 Manual Docker Commands

### Build the Application
```bash
# Build the Go application image
docker build -t rag-chatbot .
```

### Start Infrastructure Services

**CPU-Only Systems (default):**
```bash
# Start Qdrant and Ollama in CPU mode
docker-compose up -d qdrant ollama
```

**GPU-Enabled Systems:**
```bash
# Start with GPU acceleration
docker-compose -f docker-compose.yml -f docker-compose.gpu.yml up -d qdrant ollama
```

### Initialize AI Models
```bash
# Pull required models (one-time setup)
# Note: This may take 10-15 minutes on first run
docker-compose --profile init run --rm ollama-init
```

### Start the Application
```bash
# Start the RAG application
docker-compose up -d rag-app
```

### View Logs
```bash
# Follow application logs
docker-compose logs -f rag-app

# View all service logs
docker-compose logs -f
```

## 🔧 Development Mode

For development with hot-reloading:

```bash
# Start only infrastructure services
docker-compose -f docker-compose.dev.yml up -d qdrant-dev ollama-dev

# Set environment variables
export QDRANT_HOST=localhost
export GENERATIVE_ENDPOINT=http://localhost:11434
export EMBEDDING_ENDPOINT=http://localhost:11434

# Run application on host (with hot-reload)
go run cmd/main.go
```

## 🌍 Environment Configuration

The application supports these environment variables:

### AI Model Configuration
```bash
GENERATIVE_ENDPOINT=http://ollama:11434    # Ollama API URL
EMBEDDING_ENDPOINT=http://ollama:11434     # Embedding API URL  
GENERATIVE_MODEL=llama3.2:latest          # LLM model name
EMBEDDING_MODEL=nomic-embed-text:latest    # Embedding model name
```

### Qdrant Configuration
```bash
QDRANT_HOST=qdrant                         # Qdrant hostname
QDRANT_PORT=6334                           # Qdrant gRPC port
```

### Application Settings
```bash
MAX_TOKENS=512                             # Max response tokens
TEMPERATURE=0.7                            # Generation temperature
```

## 🐛 Troubleshooting

### Check Service Health
```bash
# Check all services
docker-compose ps

# Check specific service health
docker-compose exec rag-app wget -qO- http://localhost:8080/health
```

### Common Issues

**Port Conflicts:**
```bash
# Check what's using ports
lsof -i :8080
lsof -i :6333
lsof -i :11434

# Stop conflicting services
docker-compose down
```

**GPU Support:**
```bash
# Verify NVIDIA Docker runtime
docker run --rm --gpus all nvidia/cuda:11.0-base nvidia-smi

# If GPU not available, Ollama will use CPU (slower)
```

**Model Download Issues:**
```bash
# Manually download models
docker-compose exec ollama ollama pull llama3.2:latest
docker-compose exec ollama ollama pull nomic-embed-text:latest
```

**Memory Issues:**
```bash
# Check container resource usage
docker stats

# Increase Docker memory limit (Docker Desktop)
# Settings > Resources > Memory > 8GB+
```

## 🧹 Cleanup

### Stop Services
```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down --volumes
```

### Complete Cleanup
```bash
# Remove all containers, volumes, and networks
./scripts/docker-cleanup.sh
```

### Remove Images
```bash
# Remove built images
docker rmi rag-chatbot
docker rmi $(docker images -f "dangling=true" -q)
```

## 🚀 Production Deployment

For production deployment:

1. **Use environment-specific compose files:**
   ```bash
   docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
   ```

2. **Set production environment variables:**
   ```bash
   # Use external Ollama instance
   GENERATIVE_ENDPOINT=https://your-ollama-server:11434
   EMBEDDING_ENDPOINT=https://your-ollama-server:11434
   
   # Use external Qdrant instance  
   QDRANT_HOST=your-qdrant-server
   QDRANT_PORT=6334
   ```

3. **Enable resource limits:**
   ```yaml
   deploy:
     resources:
       limits:
         memory: 2G
         cpus: '1.0'
   ```

4. **Add monitoring and logging:**
   ```bash
   # Add to docker-compose.yml
   logging:
     driver: "json-file"
     options:
       max-size: "10m"
       max-file: "3"
   ```

## 📊 Monitoring

### Health Checks
All services include health checks:
- **rag-app**: HTTP health endpoint
- **qdrant**: Qdrant health API
- **ollama**: Ollama API availability

### Resource Monitoring
```bash
# Real-time resource usage
docker stats

# Service status
docker-compose ps

# Detailed service info
docker-compose exec rag-app ps aux
```

### Log Analysis
```bash
# Application logs
docker-compose logs rag-app | grep ERROR

# Infrastructure logs
docker-compose logs qdrant
docker-compose logs ollama
```
