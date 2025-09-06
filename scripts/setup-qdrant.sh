#!/bin/bash

# Setup Qdrant using Docker CLI
# This script installs and runs Qdranwhile [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -s http://localhost:6333/collections >/dev/null 2>&1; then
        echo "✅ Qdrant is running and responding!"
        echo "🌐 Web UI available at: http://localhost:6333/dashboard"
        echo "🔗 HTTP API endpoint: http://localhost:6333"
        echo "🔗 gRPC endpoint: localhost:6334"
        
        if [ "$USE_PERSISTENT_STORAGE" = true ]; then
            echo "📁 Persistent data directory: $QDRANT_DATA_DIR"
            echo "📋 Management commands:"
            echo "   Stop Qdrant: docker stop qdrant-rag"
            echo "   Start Qdrant: docker start qdrant-rag"
            echo "   View logs: docker logs qdrant-rag"
            echo "   Backup data: cp -r $QDRANT_DATA_DIR $QDRANT_DATA_DIR.backup"
        else
            echo "📁 Data storage: Docker internal volume"
            echo "📋 Management commands:"
            echo "   Stop Qdrant: docker stop qdrant-rag"
            echo "   Start Qdrant: docker start qdrant-rag"
            echo "   View logs: docker logs qdrant-rag"
            echo "   Note: Data will be lost when container is removed"
        fi
        exit 0
    fi
    
    echo "⏳ Still waiting for Qdrant... (attempt $((RETRY_COUNT + 1))/$MAX_RETRIES)"
    sleep 2
    RETRY_COUNT=$((RETRY_COUNT + 1))
doner the RAG chatbot

set -e

echo "🚀 Setting up Qdrant Vector Database..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    echo "Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Stop and remove existing Qdrant container if it exists
echo "🧹 Cleaning up existing Qdrant containers..."
docker stop qdrant-rag 2>/dev/null || true
docker rm qdrant-rag 2>/dev/null || true

# Try persistent storage first, fallback to Docker volume if it fails
QDRANT_DATA_DIR="$(pwd)/qdrant_storage"
USE_PERSISTENT_STORAGE=true

echo "📁 Attempting to create persistent data directory: $QDRANT_DATA_DIR"
if mkdir -p "$QDRANT_DATA_DIR" 2>/dev/null; then
    echo "✅ Directory created successfully"
else
    echo "⚠️  Cannot create persistent directory, will use Docker volume instead"
    USE_PERSISTENT_STORAGE=false
fi

echo "📦 Pulling Qdrant Docker image..."
docker pull qdrant/qdrant:latest

echo "🏃 Starting Qdrant container..."
if [ "$USE_PERSISTENT_STORAGE" = true ]; then
    echo "   Trying with persistent storage: $QDRANT_DATA_DIR"
    # Try with persistent storage and check if it actually works
    docker run -d \
      --name qdrant-rag \
      --restart unless-stopped \
      -p 6333:6333 \
      -p 6334:6334 \
      -v "$QDRANT_DATA_DIR:/qdrant/storage" \
      -e QDRANT__SERVICE__HTTP_PORT=6333 \
      -e QDRANT__SERVICE__GRPC_PORT=6334 \
      qdrant/qdrant:latest
    
    # Give it a moment to start and check for permission errors
    sleep 5
    if docker logs qdrant-rag 2>&1 | grep -q "Permission denied"; then
        echo "⚠️  Permission denied error detected, retrying with Docker volume..."
        docker stop qdrant-rag 2>/dev/null || true
        docker rm qdrant-rag 2>/dev/null || true
        USE_PERSISTENT_STORAGE=false
    else
        echo "✅ Started with persistent storage"
    fi
fi

if [ "$USE_PERSISTENT_STORAGE" = false ]; then
    echo "   Using Docker internal volume"
    docker run -d \
      --name qdrant-rag \
      --restart unless-stopped \
      -p 6333:6333 \
      -p 6334:6334 \
      -e QDRANT__SERVICE__HTTP_PORT=6333 \
      -e QDRANT__SERVICE__GRPC_PORT=6334 \
      qdrant/qdrant:latest
fi

echo "⏳ Waiting for Qdrant to start..."
sleep 3

# Check if Qdrant is responding
MAX_RETRIES=10
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -s http://localhost:6333/collections >/dev/null 2>&1; then
        echo "✅ Qdrant is running and responding!"
        echo "🌐 Web UI available at: http://localhost:6333/dashboard"
        echo "🔗 HTTP API endpoint: http://localhost:6333"
        echo "🔗 gRPC endpoint: localhost:6334"
        echo "📁 Data directory: $QDRANT_DATA_DIR"
        echo ""
        echo "� Management commands:"
        echo "   Stop Qdrant: docker stop qdrant-rag"
        echo "   Start Qdrant: docker start qdrant-rag"
        echo "   View logs: docker logs qdrant-rag"
        echo "   Backup data: cp -r $QDRANT_DATA_DIR $QDRANT_DATA_DIR.backup"
        exit 0
    fi
    
    echo "⏳ Still waiting for Qdrant... (attempt $((RETRY_COUNT + 1))/$MAX_RETRIES)"
    sleep 2
    RETRY_COUNT=$((RETRY_COUNT + 1))
done

echo "❌ Qdrant failed to start properly. Check logs with: docker logs qdrant-rag"
exit 1
