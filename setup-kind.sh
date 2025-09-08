#!/bin/bash

# Setup Kind cluster for RAG microservices with Ko

set -e

echo "🚀 Setting up Kind cluster for RAG microservices..."

# Check dependencies
command -v kind >/dev/null 2>&1 || { echo "❌ Kind not found. Install from: https://kind.sigs.k8s.io/docs/user/quick-start/"; exit 1; }
command -v ko >/dev/null 2>&1 || { echo "❌ Ko not found. Install from: https://ko.build/install/"; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo "❌ kubectl not found. Please install kubectl"; exit 1; }

# Create Kind cluster
CLUSTER_NAME="rag-cluster"

echo "📦 Creating Kind cluster: $CLUSTER_NAME"
cat <<EOF | kind create cluster --name $CLUSTER_NAME --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30080
    hostPort: 8080
    protocol: TCP
  - containerPort: 30333
    hostPort: 6333
    protocol: TCP
EOF

# Wait for cluster to be ready
echo "⏳ Waiting for cluster to be ready..."
kubectl wait --for=condition=Ready nodes --all --timeout=300s

# Set Ko environment for Kind
export KO_DOCKER_REPO="kind.local"
echo "🔧 Setting Ko repository to: $KO_DOCKER_REPO"

# Create namespace and secrets
echo "🏗️  Creating namespace and secrets..."
kubectl create namespace rag --dry-run=client -o yaml | kubectl apply -f -

# Create OpenAI secret (requires OPENAI_API_KEY env var)
if [ -z "$OPENAI_API_KEY" ]; then
    echo "⚠️  Warning: OPENAI_API_KEY not set. Creating dummy secret."
    echo "   Set your API key with: export OPENAI_API_KEY='your-key-here'"
    kubectl create secret generic openai-secret \
        --from-literal=OPENAI_API_KEY="dummy-key" \
        -n rag --dry-run=client -o yaml | kubectl apply -f -
else
    echo "🔑 Creating OpenAI secret..."
    kubectl create secret generic openai-secret \
        --from-literal=OPENAI_API_KEY="$OPENAI_API_KEY" \
        -n rag --dry-run=client -o yaml | kubectl apply -f -
fi

# Deploy Qdrant first
echo "📊 Deploying Qdrant vector database..."
kubectl apply -f deploy/k8s/qdrant.yaml

# Wait for Qdrant to be ready
echo "⏳ Waiting for Qdrant to be ready..."
kubectl wait --for=condition=ready pod -l app=qdrant -n rag --timeout=300s

# Build and deploy services with Ko
echo "🏗️  Building and deploying microservices with Ko..."

# Deploy all services
KO_DOCKER_REPO=kind.local ko apply -f deploy/k8s/services.yaml

# Wait for deployments
echo "⏳ Waiting for all deployments to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/crawler -n rag
kubectl wait --for=condition=available --timeout=300s deployment/parser -n rag  
kubectl wait --for=condition=available --timeout=300s deployment/embedder -n rag
kubectl wait --for=condition=available --timeout=300s deployment/rag-api -n rag

# Show status
echo "✅ Deployment complete!"
echo ""
echo "📊 Cluster Status:"
kubectl get all -n rag

echo ""
echo "🌐 Access your services:"
echo "   RAG API: http://localhost:8080"
echo "   Qdrant Dashboard: http://localhost:6333/dashboard"
echo ""
echo "🧪 Test commands:"
echo "   make test-kind"
echo "   kubectl logs -f deployment/rag-api -n rag"
echo ""
echo "🗑️  Cleanup: kind delete cluster --name $CLUSTER_NAME"
