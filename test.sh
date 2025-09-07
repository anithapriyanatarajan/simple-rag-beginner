#!/bin/bash

# RAG Pipeline Test Script
set -e

echo "🧪 Testing RAG Pipeline..."

# Check if OpenAI API key is set
if [ -z "$OPENAI_API_KEY" ]; then
    echo "❌ Error: OPENAI_API_KEY environment variable is required"
    echo "   Run: export OPENAI_API_KEY='your-api-key-here'"
    exit 1
fi

echo "✅ OpenAI API key is set"

# Option 1: Test with Docker Compose
test_docker_compose() {
    echo "🐳 Testing with Docker Compose..."
    
    # Start services
    echo "Starting services..."
    docker-compose up -d
    
    # Wait for services to be ready
    echo "Waiting for services to start..."
    sleep 30
    
    # Test health endpoints
    echo "Testing health endpoints..."
    curl -f http://localhost:8081/health || echo "Crawler health check failed"
    curl -f http://localhost:8082/health || echo "Parser health check failed"
    curl -f http://localhost:8083/health || echo "Embedder health check failed"
    curl -f http://localhost:8080/health || echo "RAG API health check failed"
    
    # Test crawler
    echo "Testing crawler..."
    curl -X POST http://localhost:8081/crawl \
        -H "Content-Type: application/json" \
        -d '{"url": "https://httpbin.org/html"}' \
        -o crawler_response.json
    
    echo "Crawler response saved to crawler_response.json"
    
    # Wait a bit for processing
    sleep 10
    
    # Test RAG query
    echo "Testing RAG query..."
    curl -X POST http://localhost:8080/query \
        -H "Content-Type: application/json" \
        -d '{"query": "What is this about?"}' \
        -o rag_response.json
    
    echo "RAG response saved to rag_response.json"
    
    # Show responses
    echo "=== Crawler Response ==="
    cat crawler_response.json | jq .
    echo ""
    echo "=== RAG Response ==="
    cat rag_response.json | jq .
    
    # Cleanup
    echo "Stopping services..."
    docker-compose down
}

# Option 2: Test individual binaries
test_individual() {
    echo "🔧 Testing individual services..."
    
    # Build binaries
    echo "Building binaries..."
    make build
    
    # Test crawler binary
    echo "Testing crawler binary..."
    ./bin/crawler &
    CRAWLER_PID=$!
    sleep 2
    curl -f http://localhost:8080/health && echo "Crawler binary works!" || echo "Crawler binary failed"
    kill $CRAWLER_PID 2>/dev/null || true
    
    echo "Individual testing complete."
}

# Option 3: Test with Kubernetes (if available)
test_kubernetes() {
    echo "☸️  Testing with Kubernetes..."
    
    # Check if kubectl is available
    if ! command -v kubectl &> /dev/null; then
        echo "kubectl not found, skipping Kubernetes test"
        return
    fi
    
    # Check if cluster is available
    if ! kubectl cluster-info &> /dev/null; then
        echo "No Kubernetes cluster available, skipping"
        return
    fi
    
    echo "Setting up OpenAI secret..."
    make setup-openai
    
    echo "Building and deploying..."
    make docker-build
    make k8s-deploy
    
    echo "Testing deployed services..."
    make test-crawl
    make test-rag
    
    echo "Checking status..."
    make k8s-status
}

# Main menu
echo "Choose testing method:"
echo "1) Docker Compose (Recommended)"
echo "2) Individual binaries"
echo "3) Kubernetes"
echo "4) All methods"

read -p "Enter choice (1-4): " choice

case $choice in
    1)
        test_docker_compose
        ;;
    2)
        test_individual
        ;;
    3)
        test_kubernetes
        ;;
    4)
        test_individual
        test_docker_compose
        test_kubernetes
        ;;
    *)
        echo "Invalid choice"
        exit 1
        ;;
esac

echo "🎉 Testing complete!"
