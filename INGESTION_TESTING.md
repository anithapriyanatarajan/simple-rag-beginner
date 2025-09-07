# Ingestion Pipeline Testing Guide

This guide helps you test the ingestion pipeline (`crawl → parse → embed → store`) with manual URL and collection inputs.

## Prerequisites

1. **OpenAI API Key**: Required for embeddings
   ```bash
   export OPENAI_API_KEY='your-api-key-here'
   ```

2. **Docker & Docker Compose**: Ensure both are installed and running

## Quick Start

### Option 1: Automated Test Script
Run the comprehensive test with multiple URLs:

```bash
./test-ingestion.sh
```

This script will:
- Start all ingestion services (crawler, parser, embedder, qdrant)
- Test multiple URLs with different collections
- Show Qdrant collection stats
- Provide manual testing instructions

### Option 2: Manual Testing

1. **Start Services**:
   ```bash
   docker-compose -f docker-compose-ingestion.yml up -d
   ```

2. **Check Service Health**:
   ```bash
   ./manual-test.sh check
   ```

3. **Test Single URL**:
   ```bash
   # Basic test
   ./manual-test.sh test
   
   # Custom URL and collection
   ./manual-test.sh test "https://example.com" "my-docs"
   ```

4. **Manual API Calls**:
   ```bash
   # Crawl with custom collection
   curl -X POST http://localhost:8081/crawl \
     -H "Content-Type: application/json" \
     -d '{"url": "https://httpbin.org/html", "collection": "test-docs"}'
   
   # Check collections in Qdrant
   curl http://localhost:6333/collections | jq .
   
   # Check specific collection
   curl http://localhost:6333/collections/test-docs | jq .
   ```

5. **Browse Qdrant Dashboard**:
   Open http://localhost:6333/dashboard in your browser

## Service Endpoints

| Service | Port | Health Check | Purpose |
|---------|------|--------------|---------|
| Crawler | 8081 | `/health` | Web scraping with Colly |
| Parser | 8082 | `/health` | HTML parsing and chunking |
| Embedder | 8083 | `/health` | OpenAI embeddings + Qdrant storage |
| Qdrant | 6333 | `/health` | Vector database |

## API Examples

### Crawl with Collection
```bash
curl -X POST http://localhost:8081/crawl \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://docs.example.com/guide",
    "collection": "product-docs"
  }'
```

### Check Qdrant Collections
```bash
# List all collections
curl http://localhost:6333/collections

# Get collection info
curl http://localhost:6333/collections/product-docs

# Search in collection (requires actual vector)
curl -X POST http://localhost:6333/collections/product-docs/points/search \
  -H "Content-Type: application/json" \
  -d '{
    "vector": [...], 
    "limit": 10,
    "with_payload": true
  }'
```

## Pipeline Flow

```
URL Input → Crawler (8081) → Parser (8082) → Embedder (8083) → Qdrant (6333)
```

1. **Crawler**: Fetches HTML content from URL
2. **Parser**: Extracts text and creates chunks
3. **Embedder**: Generates OpenAI embeddings for chunks
4. **Qdrant**: Stores vectors in specified collection

## Troubleshooting

### Services Not Starting
```bash
# Check Docker logs
docker-compose -f docker-compose-ingestion.yml logs

# Check specific service
docker logs crawler-ingestion
docker logs parser-ingestion
docker logs embedder-ingestion
docker logs qdrant-ingestion
```

### OpenAI API Issues
```bash
# Verify API key is set
echo $OPENAI_API_KEY

# Check embedder logs for API errors
docker logs embedder-ingestion
```

### Qdrant Connection Issues
```bash
# Test Qdrant directly
curl http://localhost:6333/health
curl http://localhost:6333/collections
```

## Cleanup

```bash
# Stop services
docker-compose -f docker-compose-ingestion.yml down

# Remove volumes (if needed)
docker-compose -f docker-compose-ingestion.yml down -v
```

## Collection Management

Each crawl can specify a custom collection name:
- Collections are created automatically if they don't exist
- Each collection is isolated (separate vector space)
- Use descriptive names like `product-docs`, `support-articles`, `api-guides`

## Next Steps

After verifying ingestion works:
1. Test the full RAG pipeline with `docker-compose.yml`
2. Deploy to Kubernetes with `kubectl apply -f deploy/k8s/`
3. Scale individual services based on load requirements
