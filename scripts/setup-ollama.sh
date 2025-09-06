#!/bin/bash

# Simple RAG Beginner - Ollama Setup Script
# Installs Ollama and downloads the default model

echo "🤖 Setting up Ollama for Simple RAG Beginner..."
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default model (can be overridden)
MODEL_NAME=${1:-"llama3.2"}
OLLAMA_HOST=${OLLAMA_HOST:-"http://localhost:11434"}

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if Ollama is running
ollama_running() {
    curl -s "$OLLAMA_HOST/api/tags" >/dev/null 2>&1
}

# Function to check if model is installed
model_installed() {
    curl -s "$OLLAMA_HOST/api/tags" | grep -q "\"name\":\"$1\""
}

echo "📋 Setup Configuration:"
echo "   Model: $MODEL_NAME"
echo "   Ollama Host: $OLLAMA_HOST"
echo ""

# Step 1: Install Ollama if not present
echo "🔍 Checking for Ollama installation..."
if command_exists ollama; then
    echo -e "${GREEN}✅ Ollama is already installed${NC}"
    ollama --version
else
    echo -e "${YELLOW}📦 Installing Ollama...${NC}"
    
    # Detect OS and install accordingly
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "Installing on Linux..."
        curl -fsSL https://ollama.ai/install.sh | sh
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "Installing on macOS..."
        if command_exists brew; then
            brew install ollama
        else
            echo "Please install Homebrew first, then run: brew install ollama"
            echo "Or download from: https://ollama.ai/download"
            exit 1
        fi
    else
        echo -e "${RED}❌ Unsupported OS. Please install Ollama manually from: https://ollama.ai/download${NC}"
        exit 1
    fi
    
    if command_exists ollama; then
        echo -e "${GREEN}✅ Ollama installed successfully${NC}"
    else
        echo -e "${RED}❌ Ollama installation failed${NC}"
        exit 1
    fi
fi

# Step 2: Start Ollama service
echo ""
echo "🚀 Starting Ollama service..."

# Try to start ollama serve in background
if ! ollama_running; then
    echo "Starting ollama serve..."
    nohup ollama serve > /tmp/ollama.log 2>&1 &
    OLLAMA_PID=$!
    echo "Ollama PID: $OLLAMA_PID"
    
    # Wait for Ollama to start
    echo "Waiting for Ollama to start..."
    for i in {1..30}; do
        if ollama_running; then
            echo -e "${GREEN}✅ Ollama service is running${NC}"
            break
        fi
        echo -n "."
        sleep 1
    done
    
    if ! ollama_running; then
        echo -e "${RED}❌ Failed to start Ollama service${NC}"
        echo "Check logs: tail /tmp/ollama.log"
        exit 1
    fi
else
    echo -e "${GREEN}✅ Ollama service is already running${NC}"
fi

# Step 3: Download the model
echo ""
echo "📥 Checking for model: $MODEL_NAME"

if model_installed "$MODEL_NAME"; then
    echo -e "${GREEN}✅ Model '$MODEL_NAME' is already installed${NC}"
else
    echo -e "${YELLOW}📥 Downloading model '$MODEL_NAME'...${NC}"
    echo "This may take several minutes depending on model size and internet speed."
    
    if ollama pull "$MODEL_NAME"; then
        echo -e "${GREEN}✅ Model '$MODEL_NAME' downloaded successfully${NC}"
    else
        echo -e "${RED}❌ Failed to download model '$MODEL_NAME'${NC}"
        echo ""
        echo "Available models to try:"
        echo "  • llama3.2 (default, ~2GB)"
        echo "  • llama3.2:1b (smaller, ~1.3GB)"
        echo "  • qwen2.5:0.5b (very small, ~0.5GB)"
        echo "  • phi3:mini (small, ~2.3GB)"
        echo ""
        echo "Usage: $0 [model-name]"
        echo "Example: $0 llama3.2:1b"
        exit 1
    fi
fi

# Step 4: Test the installation
echo ""
echo "🧪 Testing Ollama installation..."

# Test basic functionality
TEST_RESPONSE=$(curl -s -X POST "$OLLAMA_HOST/api/generate" \
    -H "Content-Type: application/json" \
    -d "{\"model\":\"$MODEL_NAME\",\"prompt\":\"Hello\",\"stream\":false}" \
    --max-time 30)

if echo "$TEST_RESPONSE" | grep -q '"response"'; then
    echo -e "${GREEN}✅ Ollama is working correctly with model '$MODEL_NAME'${NC}"
    RESPONSE=$(echo "$TEST_RESPONSE" | grep -o '"response":"[^"]*"' | cut -d'"' -f4)
    echo "Test response: $RESPONSE"
else
    echo -e "${RED}❌ Ollama test failed${NC}"
    echo "Response: $TEST_RESPONSE"
    exit 1
fi

# Step 5: Show configuration
echo ""
echo "🎉 Ollama setup complete!"
echo ""
echo "📋 Configuration Summary:"
echo "   • Ollama Host: $OLLAMA_HOST"
echo "   • Model: $MODEL_NAME" 
echo "   • Status: ✅ Ready"
echo ""
echo "🔧 Environment Variables (optional):"
echo "   export GENERATIVE_PROVIDER=ollama"
echo "   export GENERATIVE_MODEL=$MODEL_NAME"
echo "   export GENERATIVE_ENDPOINT=$OLLAMA_HOST"
echo ""
echo "🚀 Next Steps:"
echo "   1. Start your RAG application: go run cmd/main.go"
echo "   2. Test the /model/info endpoint: curl http://localhost:8080/model/info"
echo "   3. Try a query: curl -X POST http://localhost:8080/query -H 'Content-Type: application/json' -d '{\"query\":\"Hello\"}'"
echo ""
echo "📊 Monitor Ollama:"
echo "   • Logs: tail /tmp/ollama.log"
echo "   • Models: ollama list"
echo "   • Stop: killall ollama"
echo ""

# Optional: Show available models
echo "📦 Available models on this system:"
ollama list
