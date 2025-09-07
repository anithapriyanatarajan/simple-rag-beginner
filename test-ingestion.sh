#!/bin/bash

# Ingestion Pipeline Test - Crawl + Parse + Embed + Store with Custom Collection
set -e

echo "🔄 Testing Ingestion Pipeline: Crawl → Parse → Embed → Store"
echo "This will test ONLY the ingestion services (no RAG API)"

# Check prerequisites
if [ -z "$OPENAI_API_KEY" ]; then
    echo "❌ Error: OPENAI_API_KEY environment variable is required"
    echo "   Run: export OPENAI_API_KEY='your-api-key-here'"
    exit 1
fi

echo "✅ OpenAI API key is set"

# Function to test ingestion with specific URL and collection
test_ingestion() {
    local url="$1"
    local collection="$2"
    
    echo ""
    echo "📡 Testing ingestion for:"
    echo "   URL: $url"
    echo "   Collection: $collection"
    
    # Start crawling with specific collection
    echo "🚀 Starting crawl..."
    CRAWL_RESPONSE=$(curl -s -X POST http://localhost:8081/crawl \
        -H "Content-Type: application/json" \
        -d "{\"url\": \"$url\", \"collection\": \"$collection\"}")
    
    echo "Crawl response: $CRAWL_RESPONSE" | jq .
    
    # Wait for processing pipeline
    echo "⏳ Waiting for processing pipeline (60 seconds)..."
    sleep 60
    
    # Check Qdrant collection
    echo "🔍 Checking Qdrant collection '$collection'..."
    COLLECTION_INFO=$(curl -s "http://localhost:6333/collections/$collection")
    echo "Collection info: $COLLECTION_INFO" | jq .
    
    # Get collection stats
    COLLECTION_STATS=$(curl -s "http://localhost:6333/collections/$collection" | jq -r '.result.points_count // "N/A"')
    echo "📊 Points in collection '$collection': $COLLECTION_STATS"
    
    return 0
}

# Start ingestion services only
echo "🐳 Starting ingestion services (crawler, parser, embedder, qdrant)..."
docker-compose -f docker-compose-ingestion.yml up -d

echo "⏳ Waiting for services to be ready (30 seconds)..."
sleep 30

# Check service health
echo "🏥 Checking service health..."
for service in "crawler:8081" "parser:8082" "embedder:8083"; do
    name=$(echo $service | cut -d: -f1)
    port=$(echo $service | cut -d: -f2)
    if curl -f http://localhost:${port}/health > /dev/null 2>&1; then
        echo "✅ $name is healthy"
    else
        echo "❌ $name is not responding"
    fi
done

# Check Qdrant
if curl -f http://localhost:6333/collections > /dev/null 2>&1; then
    echo "✅ Qdrant is healthy"
else
    echo "❌ Qdrant is not responding"
fi

# Show initial Qdrant collections
echo ""
echo "📋 Initial Qdrant collections:"
curl -s http://localhost:6333/collections | jq .

# Test multiple URLs with different collections
echo ""
echo "🧪 Running ingestion tests..."

# Test 1: Simple HTML page
test_ingestion "https://httpbin.org/html" "test-html"

# Test 2: JSON response page  
test_ingestion "https://httpbin.org/json" "test-json"

# Test 3: Custom collection name
test_ingestion "https://httpbin.org/robots.txt" "my-custom-collection"

# Show final Qdrant state
echo ""
echo "📊 Final Qdrant collections:"
FINAL_COLLECTIONS=$(curl -s http://localhost:6333/collections)
echo "$FINAL_COLLECTIONS" | jq .

# Show point counts for each collection
echo ""
echo "📈 Point counts by collection:"
echo "$FINAL_COLLECTIONS" | jq -r '.result.collections[]?.name' | while read collection; do
    if [ -n "$collection" ]; then
        count=$(curl -s "http://localhost:6333/collections/$collection" | jq -r '.result.points_count // 0')
        echo "   $collection: $count points"
    fi
done

# Manual testing instructions
echo ""
echo "🎯 Manual Testing:"
echo "Services are still running. You can test manually:"
echo ""
echo "1. Crawl with custom collection:"
echo "   curl -X POST http://localhost:8081/crawl \\"
echo "     -H \"Content-Type: application/json\" \\"
echo "     -d '{\"url\": \"https://example.com\", \"collection\": \"my-docs\"}'"
echo ""
echo "2. Check Qdrant collections:"
echo "   curl http://localhost:6333/collections"
echo ""
echo "3. Check specific collection:"
echo "   curl http://localhost:6333/collections/my-docs"
echo ""
echo "4. Browse Qdrant dashboard:"
echo "   http://localhost:6333/dashboard"
echo ""

read -p "Press Enter to stop services or Ctrl+C to keep them running..."

# Cleanup
echo "🧹 Stopping services..."
docker-compose -f docker-compose-ingestion.yml down

echo "✅ Ingestion pipeline test complete!"
