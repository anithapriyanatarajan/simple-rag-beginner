# RAG Microservices Pipeline Makefile

REGISTRY ?= localhost:5000
TAG ?= latest

# Service names
SERVICES = crawler parser embedder rag-api

.PHONY: help build docker-build docker-push k8s-deploy k8s-undeploy clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Docker Compose targets (Primary)
up: ## Start all services with Docker Compose
	docker-compose up -d

down: ## Stop all services with Docker Compose
	docker-compose down -v

build: ## Build all Docker images
	docker-compose build

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

# Kubernetes targets (Optional)
k8s-deploy: docker-build ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	kubectl apply -f deploy/k8s/qdrant.yaml
	kubectl apply -f deploy/k8s/services.yaml
	kubectl apply -f deploy/k8s/api.yaml
	@echo "Waiting for deployments to be ready..."
	kubectl wait --for=condition=available --timeout=300s deployment/crawler -n rag
	kubectl wait --for=condition=available --timeout=300s deployment/parser -n rag
	kubectl wait --for=condition=available --timeout=300s deployment/embedder -n rag
	kubectl wait --for=condition=available --timeout=300s deployment/rag-api -n rag
	kubectl wait --for=condition=ready --timeout=300s statefulset/qdrant -n rag

k8s-undeploy: ## Remove from Kubernetes
	@echo "Removing from Kubernetes..."
	kubectl delete -f deploy/k8s/api.yaml --ignore-not-found
	kubectl delete -f deploy/k8s/services.yaml --ignore-not-found
	kubectl delete -f deploy/k8s/qdrant.yaml --ignore-not-found
	kubectl delete namespace rag --ignore-not-found

k8s-status: ## Show Kubernetes status
	kubectl get all -n rag
	kubectl get pvc -n rag

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
