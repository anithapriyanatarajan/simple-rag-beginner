#!/bin/bash

# Build and push images with Ko (without deploying)

set -e

echo "🏗️  Building container images with Ko..."

# Set Ko repository
export KO_DOCKER_REPO=${KO_DOCKER_REPO:-"ko.local"}
echo "📦 Using repository: $KO_DOCKER_REPO"

# Build all services
echo "Building microservices..."

services=("crawler" "parser" "embedder" "rag-api")

for service in "${services[@]}"; do
    echo "  📦 Building $service..."
    ko build ./cmd/$service --bare
done

echo "✅ All images built successfully!"
echo ""
echo "🏷️  Images tagged with repository: $KO_DOCKER_REPO"
echo "🚀 To deploy: make kind-deploy"
