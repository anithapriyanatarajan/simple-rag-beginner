# Kubernetes RAG Pipeline Makefile

REGISTRY ?= localhost:5000
TAG ?= latest

# Service names
SERVICES = crawler parser embedder rag-api

.PHONY: help build docker-build docker-push k8s-deploy k8s-undeploy clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build all Go binaries
	@echo "Building Go services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd cmd/$$service && go build -o ../../bin/$$service . && cd ../..; \
	done

docker-build: ## Build all Docker images
	@echo "Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building $(REGISTRY)/rag/$$service:$(TAG)..."; \
		docker build -f cmd/$$service/Dockerfile -t $(REGISTRY)/rag/$$service:$(TAG) .; \
	done

docker-push: docker-build ## Push Docker images to registry
	@echo "Pushing Docker images..."
	@for service in $(SERVICES); do \
		echo "Pushing $(REGISTRY)/rag/$$service:$(TAG)..."; \
		docker push $(REGISTRY)/rag/$$service:$(TAG); \
	done

k8s-deploy: ## Deploy to Kubernetes
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

k8s-logs: ## Show logs from all services
	@echo "=== Crawler Logs ==="
	kubectl logs -l app=crawler -n rag --tail=50
	@echo "=== Parser Logs ==="
	kubectl logs -l app=parser -n rag --tail=50
	@echo "=== Embedder Logs ==="
	kubectl logs -l app=embedder -n rag --tail=50
	@echo "=== RAG API Logs ==="
	kubectl logs -l app=rag-api -n rag --tail=50
	@echo "=== Qdrant Logs ==="
	kubectl logs -l app=qdrant -n rag --tail=50

k8s-status: ## Show Kubernetes status
	kubectl get all -n rag
	kubectl get pvc -n rag
	kubectl get configmaps -n rag
	kubectl get secrets -n rag

test-crawl: ## Test crawler service
	@echo "Testing crawler..."
	kubectl port-forward svc/crawler 8081:8080 -n rag &
	sleep 2
	curl -X POST http://localhost:8081/crawl -H "Content-Type: application/json" -d '{"url": "https://example.com"}'
	pkill -f "kubectl port-forward svc/crawler"

test-rag: ## Test RAG API
	@echo "Testing RAG API..."
	kubectl port-forward svc/rag-api 8082:8080 -n rag &
	sleep 2
	curl -X POST http://localhost:8082/query -H "Content-Type: application/json" -d '{"query": "What is this about?"}'
	pkill -f "kubectl port-forward svc/rag-api"

clean: ## Clean up build artifacts
	rm -rf bin/
	docker image prune -f

setup-openai: ## Setup OpenAI API key (requires OPENAI_API_KEY env var)
	@if [ -z "$(OPENAI_API_KEY)" ]; then \
		echo "Error: OPENAI_API_KEY environment variable is required"; \
		exit 1; \
	fi
	@echo "Setting up OpenAI API key..."
	kubectl create secret generic openai-secret \
		--from-literal=OPENAI_API_KEY=$(OPENAI_API_KEY) \
		-n rag --dry-run=client -o yaml | kubectl apply -f -

# Development targets
dev-build: ## Build for development
	@mkdir -p bin
	@for service in $(SERVICES); do \
		echo "Building $$service for development..."; \
		cd cmd/$$service && go build -race -o ../../bin/$$service . && cd ../..; \
	done

dev-test: ## Run tests
	go test ./...

dev-lint: ## Run linter
	golangci-lint run

# Local development with Docker Compose (alternative to K8s)
dev-up: ## Start services with Docker Compose
	docker-compose up -d

dev-down: ## Stop services with Docker Compose
	docker-compose down -v

# Default target
all: build docker-build k8s-deploy
