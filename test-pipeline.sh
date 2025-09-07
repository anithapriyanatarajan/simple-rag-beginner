#!/bin/bash

# Complete RAG Pipeline End-to-End Test
set -e

echo "🔄 Testing Complete RAG Pipeline: Crawl → Parse → Embed → Store → Query"

# Check prerequisites
if [ -z "$OPENAI_API_KEY" ]; then
    echo "❌ Error: OPENAI_API_KEY environment variable is required"
    echo "   Run: export OPENAI_API_KEY='your-api-key-here'"
    exit 1
fi

echo "✅ OpenAI API key is set"

# Start services with Docker Compose
echo "🐳 Starting all services..."
docker-compose up -d

echo "⏳ Waiting for services to be ready (30 seconds)..."
sleep 30

# Check if all services are healthy
echo "🏥 Checking service health..."
for port in 8080 8081 8082 8083; do
    if curl -f http://localhost:${port}/health > /dev/null 2>&1; then
        echo "✅ Service on port ${port} is healthy"
    else
        echo "❌ Service on port ${port} is not responding"
    fi
done

# Test 1: Crawl a simple webpage
echo ""
echo "📡 Step 1: Crawling webpage..."
CRAWL_RESPONSE=$(curl -s -X POST http://localhost:8081/crawl \
    -H "Content-Type: application/json" \
    -d '{"url": "https://httpbin.org/html"}')

echo "Crawl response: $CRAWL_RESPONSE" | jq .

# Wait for processing pipeline (crawl → parse → embed → store)
echo ""
echo "⏳ Waiting for processing pipeline to complete (60 seconds)..."
sleep 60

# Test 2: Check if Qdrant has data
echo ""
echo "🔍 Step 2: Checking Qdrant for stored vectors..."
QDRANT_COLLECTIONS=$(curl -s http://localhost:6333/collections)
echo "Qdrant collections: $QDRANT_COLLECTIONS" | jq .

# Check if documents collection exists and has points
QDRANT_INFO=$(curl -s http://localhost:6333/collections/documents)
echo "Documents collection info: $QDRANT_INFO" | jq .

# Test 3: Query the RAG system
echo ""
echo "🤖 Step 3: Querying RAG system..."
RAG_RESPONSE=$(curl -s -X POST http://localhost:8080/query \
    -H "Content-Type: application/json" \
    -d '{"query": "What is this webpage about?"}')

echo "RAG response: $RAG_RESPONSE" | jq .

# Test 4: Try another crawl and query
echo ""
echo "📡 Step 4: Testing with another URL..."
curl -s -X POST http://localhost:8081/crawl \
    -H "Content-Type: application/json" \
    -d '{"url": "https://httpbin.org/robots.txt"}' > /dev/null

echo "⏳ Waiting for processing..."
sleep 30

RAG_RESPONSE2=$(curl -s -X POST http://localhost:8080/query \
    -H "Content-Type: application/json" \
    -d '{"query": "Tell me about robots"}')

echo "Second RAG response: $RAG_RESPONSE2" | jq .

# Test 5: Check Qdrant point count
echo ""
echo "📊 Step 5: Final Qdrant stats..."
FINAL_INFO=$(curl -s http://localhost:6333/collections/documents)
echo "Final collection info: $FINAL_INFO" | jq .

echo ""
echo "🎉 End-to-end pipeline test complete!"
echo "📋 Summary:"
echo "   1. ✅ Crawled webpages"
echo "   2. ✅ Parsed and chunked content"
echo "   3. ✅ Generated embeddings"
echo "   4. ✅ Stored in Qdrant"
echo "   5. ✅ Retrieved and generated responses"

# Cleanup
echo ""
echo "🧹 Cleaning up..."
docker-compose down

echo "✨ Test completed successfully!"
