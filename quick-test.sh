#!/bin/bash

# Quick test to verify the complete pipeline workflow
echo "🔗 Testing Complete RAG Pipeline Workflow"

if [ -z "$OPENAI_API_KEY" ]; then
    echo "❌ Set OPENAI_API_KEY first: export OPENAI_API_KEY='your-key'"
    exit 1
fi

echo "🚀 Starting services..."
docker-compose up -d

echo "⏳ Waiting 30 seconds for services to start..."
sleep 30

echo "🌐 Testing single endpoint crawl-and-query..."
curl -X POST http://localhost:8080/crawl-and-query \
    -H "Content-Type: application/json" \
    -d '{
        "url": "https://httpbin.org/html",
        "query": "What is this webpage about?"
    }' | jq .

echo "🔍 Checking Qdrant for stored data..."
curl -s http://localhost:6333/collections/documents | jq .

echo "🧹 Cleaning up..."
docker-compose down

echo "✅ Quick test complete!"
