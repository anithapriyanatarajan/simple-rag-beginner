#!/bin/bash

# Docker deployment script for RAG Chatbot
# This script sets up and runs the full RAG stack with Docker Compose

set -e

echo "🚀 RAG Chatbot Docker Deployment"
echo "================================="

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Function to wait for service to be healthy
wait_for_service() {
    local service_name=$1
    local max_attempts=30
    local attempt=1
    
    echo "⏳ Waiting for $service_name to be healthy..."
    
    while [ $attempt -le $max_attempts ]; do
        if docker-compose ps $service_name | grep -q "healthy"; then
            echo "✅ $service_name is healthy"
            return 0
        fi
        echo "   Attempt $attempt/$max_attempts: $service_name not ready yet..."
        sleep 10
        attempt=$((attempt + 1))
    done
    
    echo "❌ $service_name failed to become healthy after $max_attempts attempts"
    return 1
}

# Clean up any existing containers
echo "🧹 Cleaning up existing containers..."
docker-compose down --volumes --remove-orphans 2>/dev/null || true

# Start the infrastructure services first
echo "🏗️  Starting infrastructure services..."
docker-compose up -d qdrant ollama

# Wait for infrastructure to be ready
wait_for_service qdrant
wait_for_service ollama

# Initialize Ollama models
echo "📥 Pulling Ollama models..."
docker-compose --profile init run --rm ollama-init

echo "✅ Models pulled successfully!"

# Start the application
echo "🚀 Starting RAG application..."
docker-compose up -d rag-app

# Wait for application to be ready
wait_for_service rag-app

echo ""
echo "🎉 RAG Chatbot deployed successfully!"
echo ""
echo "📋 Service URLs:"
echo "   • RAG Chatbot:      http://localhost:8080"
echo "   • Qdrant Dashboard: http://localhost:6333/dashboard"
echo "   • Ollama API:       http://localhost:11434"
echo ""
echo "🧪 Test the deployment:"
echo "   curl -X POST http://localhost:8080/query \\"
echo "        -H 'Content-Type: application/json' \\"
echo "        -d '{\"query\": \"hello there\"}'"
echo ""
echo "📊 Monitor services:"
echo "   docker-compose logs -f rag-app"
echo "   docker-compose ps"
echo ""
echo "🛑 Stop services:"
echo "   docker-compose down"
