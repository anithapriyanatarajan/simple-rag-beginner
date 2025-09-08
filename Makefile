# RAG Microservices Pipeline Makefile

REGISTRY ?= localhost:5000
TAG ?= latest

# Service names
SERVICES = crawler parser embedder rag-api

.PHONY: help build docker-build docker-push k8s-deploy k8s-undeploy clean config-check

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Docker Compose targets (Primary)
up: ## Start all services with Docker Compose
	docker-compose up -d

down: ## Stop all services with Docker Compose
	docker-compose down -v

build: ## Build all Docker images
	docker-compose build

build-ko: ## Build container images with Ko
	./build-ko.sh

restart: ## Restart all services
	docker-compose restart

logs: ## Show logs from all services
	docker-compose logs -f

status: ## Show service status
	docker-compose ps

# Individual service operations
restart-api: ## Restart just the RAG API
	docker-compose restart rag-api

logs-api: ## Show RAG API logs
	docker-compose logs -f rag-api

# Testing targets
test: ## Test the RAG system (both modes)
	@echo "Testing RAG mode..."
	curl -X POST http://localhost:8080/query \
		-H "Content-Type: application/json" \
		-d '{"query": "What is Tekton?"}'
	@echo "\nTesting Direct mode..."
	curl -X POST http://localhost:8080/query-direct \
		-H "Content-Type: application/json" \
		-d '{"query": "What is Tekton?"}'

test-crawl: ## Test web crawling
	curl -X POST http://localhost:8080/crawl-and-query \
		-H "Content-Type: application/json" \
		-d '{"url": "https://tekton.dev", "query": "What is Tekton?", "collection": "test-collection"}'

health: ## Check service health
	@echo "Checking RAG API health..."
	curl -s http://localhost:8080/health | jq .
	@echo "Checking Qdrant health..."
	curl -s http://localhost:6333/health | jq .

config-check: ## Check current provider configuration
	@echo "Current RAG API configuration:"
	curl -s http://localhost:8080/config | jq .

config-demo: ## Demo different provider configurations
	@echo "Testing OpenAI configuration:"
	cd cmd/rag-api && LLM_PROVIDER=openai EMBEDDING_PROVIDER=openai timeout 2s go run . 2>&1 | grep "rag-api service started"
	@echo "Testing Ollama configuration:"
	cd cmd/rag-api && LLM_PROVIDER=ollama EMBEDDING_PROVIDER=ollama timeout 2s go run . 2>&1 | grep "rag-api service started" || true
	@echo "Testing Mixed configuration (OpenAI LLM + Cohere Embedding):"
	cd cmd/rag-api && LLM_PROVIDER=openai EMBEDDING_PROVIDER=cohere timeout 2s go run . 2>&1 | grep "rag-api service started" || true

# Kubernetes targets (Kind cluster)
kind-setup: ## Setup Kind cluster with Ko
	./setup-kind.sh

kind-deploy: ## Deploy to existing Kind cluster with Ko
	@echo "🚀 Deploying to Kind cluster with Ko..."
	export KO_DOCKER_REPO=kind.local && ko apply -f deploy/k8s/

kind-undeploy: ## Remove from Kind cluster
	@echo "🗑️  Removing from Kind cluster..."
	kubectl delete -f deploy/k8s/services.yaml --ignore-not-found
	kubectl delete -f deploy/k8s/qdrant.yaml --ignore-not-found
	kubectl delete namespace rag --ignore-not-found

kind-destroy: ## Destroy Kind cluster
	kind delete cluster --name rag-cluster

kind-status: ## Show Kind cluster status
	kubectl get all -n rag
	kubectl get pvc -n rag

kind-logs: ## Show logs from Kind cluster
	kubectl logs -f deployment/rag-api -n rag

test-kind: ## Test RAG system on Kind cluster
	@echo "Testing RAG on Kind cluster..."
	curl -X POST http://localhost:8080/query \
		-H "Content-Type: application/json" \
		-d '{"query": "What is Tekton?"}'

# Legacy Kubernetes targets (Generic)
k8s-deploy: kind-deploy ## Alias for kind-deploy
k8s-undeploy: kind-undeploy ## Alias for kind-undeploy
k8s-status: kind-status ## Alias for kind-status

setup-openai: ## Setup OpenAI API key (requires OPENAI_API_KEY env var)
	@if [ -z "$(OPENAI_API_KEY)" ]; then \
		echo "Error: OPENAI_API_KEY environment variable is required"; \
		exit 1; \
	fi
	@echo "Setting up OpenAI API key..."
	kubectl create secret generic openai-secret \
		--from-literal=OPENAI_API_KEY=$(OPENAI_API_KEY) \
		-n rag --dry-run=client -o yaml | kubectl apply -f -

# Cleanup
clean: ## Clean up Docker resources
	docker-compose down -v
	docker system prune -f

# Development
dev: ## Start services in development mode with auto-rebuild
	docker-compose up --build

# Default target
all: build up
