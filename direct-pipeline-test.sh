#!/bin/bash

# Direct Pipeline Test - Bypassing async issues
echo "🔧 Testing Pipeline Directly (Sync Mode)"
echo "========================================"

# Step 1: Crawl and get HTML
echo "1️⃣ Crawling URL..."
CRAWL_DATA=$(curl -s -X POST http://localhost:8081/crawl \
    -H "Content-Type: application/json" \
    -d '{"url": "https://httpbin.org/html", "collection": "direct-test"}')

echo "Crawl response received"

# Step 2: Extract data and send to parser
echo ""
echo "2️⃣ Parsing HTML..."
PARSE_DATA=$(echo "$CRAWL_DATA" | curl -s -X POST http://localhost:8082/parse \
    -H "Content-Type: application/json" \
    -d @-)

echo "Parse response:"
echo "$PARSE_DATA" | jq .

# Step 3: Extract chunks and send to embedder
echo ""
echo "3️⃣ Embedding chunks..."
CHUNKS=$(echo "$PARSE_DATA" | jq -c '.chunks')
EMBED_PAYLOAD=$(jq -n --argjson chunks "$CHUNKS" '{"chunks": $chunks, "collection": "direct-test"}')

EMBED_RESPONSE=$(echo "$EMBED_PAYLOAD" | curl -s -X POST http://localhost:8083/embed \
    -H "Content-Type: application/json" \
    -d @-)

echo "Embed response:"
echo "$EMBED_RESPONSE" | jq .

# Step 4: Check Qdrant
echo ""
echo "4️⃣ Checking Qdrant collections..."
sleep 5
curl -s http://localhost:6333/collections | jq .

echo ""
echo "5️⃣ Checking direct-test collection..."
curl -s http://localhost:6333/collections/direct-test 2>/dev/null | jq . || echo "Collection 'direct-test' not found"

# Show logs
echo ""
echo "📋 Service Logs After Direct Test:"
echo "--- Parser Logs ---"
docker logs parser-ingestion --tail 3

echo ""
echo "--- Embedder Logs ---"  
docker logs embedder-ingestion --tail 3
