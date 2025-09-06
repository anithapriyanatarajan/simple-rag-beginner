#!/bin/bash

# Docker cleanup script for RAG Chatbot
# This script stops and removes all Docker containers, volumes, and networks

set -e

echo "🧹 RAG Chatbot Docker Cleanup"
echo "============================="

echo "🛑 Stopping all services..."
docker-compose down --volumes --remove-orphans

echo "🗑️  Removing unused Docker resources..."
docker system prune -f

echo "📦 Removing unused volumes..."
docker volume prune -f

echo "🌐 Removing unused networks..."
docker network prune -f

echo ""
echo "✅ Cleanup completed!"
echo ""
echo "🔄 To redeploy:"
echo "   ./scripts/docker-deploy.sh"
