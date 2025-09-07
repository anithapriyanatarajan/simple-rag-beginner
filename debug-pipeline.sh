#!/bin/bash

# Step-by-step pipeline testing
set -e

echo "🔍 Testing Pipeline Step by Step"
echo "================================"

echo ""
echo "1️⃣ Testing Crawler Service..."
CRAWLER_RESPONSE=$(curl -s -X POST http://localhost:8081/crawl \
    -H "Content-Type: application/json" \
    -d '{"url": "https://httpbin.org/html", "collection": "step-test"}')

echo "Crawler Response:"
echo "$CRAWLER_RESPONSE" | jq .

echo ""
echo "2️⃣ Wait 10 seconds for pipeline processing..."
sleep 10

echo ""
echo "3️⃣ Check all service logs for activity..."
echo "--- Crawler Logs ---"
docker logs crawler-ingestion --tail 5

echo ""
echo "--- Parser Logs ---"
docker logs parser-ingestion --tail 5

echo ""
echo "--- Embedder Logs ---"
docker logs embedder-ingestion --tail 5

echo ""
echo "4️⃣ Check Qdrant collections..."
curl -s http://localhost:6333/collections | jq .

echo ""
echo "5️⃣ Manual step testing..."
echo ""
echo "Test parser directly:"
echo "curl -X POST http://localhost:8082/parse \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"html\": \"<p>test content</p>\", \"url\": \"https://test.com\", \"collection\": \"manual-test\"}'"

echo ""
echo "Test embedder directly:"
echo "curl -X POST http://localhost:8083/embed \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{\"chunks\": [{\"id\": \"test1\", \"text\": \"test content\", \"url\": \"https://test.com\"}], \"collection\": \"manual-test\"}'"
