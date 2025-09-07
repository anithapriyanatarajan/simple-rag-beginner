#!/bin/bash

# Manual Ingestion Testing - Simple commands for manual testing
echo "🎯 Manual Ingestion Pipeline Testing"
echo "===================================="

# Check if services are running
check_services() {
    echo "🏥 Checking service health..."
    
    for service in "crawler:8081" "parser:8082" "embedder:8083" "qdrant:6333"; do
        name=$(echo $service | cut -d: -f1)
        port=$(echo $service | cut -d: -f2)
        if curl -f http://localhost:${port}/health > /dev/null 2>&1 || curl -f http://localhost:${port}/collections > /dev/null 2>&1; then
            echo "✅ $name (localhost:$port) - HEALTHY"
        else
            echo "❌ $name (localhost:$port) - NOT RESPONDING"
        fi
    done
}

# Show usage
show_usage() {
    echo ""
    echo "📖 Usage Commands:"
    echo "=================="
    echo ""
    echo "1. Start ingestion services:"
    echo "   docker-compose -f docker-compose-ingestion.yml up -d"
    echo ""
    echo "2. Test crawling with custom collection:"
    echo "   curl -X POST http://localhost:8081/crawl \\"
    echo "     -H \"Content-Type: application/json\" \\"
    echo "     -d '{\"url\": \"https://httpbin.org/html\", \"collection\": \"test-collection\"}'"
    echo ""
    echo "3. Check Qdrant collections:"
    echo "   curl http://localhost:6333/collections | jq ."
    echo ""
    echo "4. Check specific collection points:"
    echo "   curl http://localhost:6333/collections/test-collection | jq ."
    echo ""
    echo "5. Search in collection:"
    echo "   curl -X POST http://localhost:6333/collections/test-collection/points/search \\"
    echo "     -H \"Content-Type: application/json\" \\"
    echo "     -d '{\"vector\": [0.1, 0.2, 0.3], \"limit\": 5}'"
    echo ""
    echo "6. Browse Qdrant dashboard:"
    echo "   http://localhost:6333/dashboard"
    echo ""
    echo "7. Stop services:"
    echo "   docker-compose -f docker-compose-ingestion.yml down"
    echo ""
}

# Quick test function
quick_test() {
    local url="${1:-https://httpbin.org/html}"
    local collection="${2:-quick-test}"
    
    echo "🚀 Quick test with URL: $url, Collection: $collection"
    
    curl -X POST http://localhost:8081/crawl \
        -H "Content-Type: application/json" \
        -d "{\"url\": \"$url\", \"collection\": \"$collection\"}" \
        | jq .
    
    echo "⏳ Waiting 30 seconds for processing..."
    sleep 30
    
    echo "📊 Collection stats:"
    curl -s http://localhost:6333/collections/$collection | jq .
}

# Main script logic
case "${1:-help}" in
    "check")
        check_services
        ;;
    "test")
        quick_test "$2" "$3"
        ;;
    "help"|*)
        show_usage
        echo ""
        echo "🔧 Script Commands:"
        echo "==================="
        echo "./manual-test.sh check          - Check service health"
        echo "./manual-test.sh test [url] [collection] - Quick test (optional params)"
        echo "./manual-test.sh help           - Show this help"
        echo ""
        check_services
        ;;
esac
