#!/bin/bash

# Simple RAG Beginner - Cleanup Script
# Stops all containers and applications related to this project

echo "🧹 Cleaning up Simple RAG Beginner project..."
echo ""

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if a process is running on a port
port_in_use() {
    lsof -i :"$1" >/dev/null 2>&1
}

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 1. Stop Go applications running on port 8080
echo "🔍 Checking for Go applications on port 8080..."
if port_in_use 8080; then
    echo -e "${YELLOW}Found applications using port 8080${NC}"
    
    # Get PIDs of processes using port 8080
    PIDS=$(lsof -t -i :8080 2>/dev/null)
    
    if [ -n "$PIDS" ]; then
        echo "Stopping processes: $PIDS"
        kill $PIDS 2>/dev/null
        sleep 2
        
        # Force kill if still running
        REMAINING_PIDS=$(lsof -t -i :8080 2>/dev/null)
        if [ -n "$REMAINING_PIDS" ]; then
            echo "Force killing remaining processes: $REMAINING_PIDS"
            kill -9 $REMAINING_PIDS 2>/dev/null
        fi
        echo -e "${GREEN}✅ Go applications stopped${NC}"
    fi
else
    echo -e "${GREEN}✅ No applications running on port 8080${NC}"
fi

# 2. Stop any Go processes running main.go
echo ""
echo "🔍 Checking for Go processes running main.go..."
GO_PIDS=$(ps aux | grep "go run.*main.go" | grep -v grep | awk '{print $2}' 2>/dev/null)

if [ -n "$GO_PIDS" ]; then
    echo -e "${YELLOW}Found Go processes: $GO_PIDS${NC}"
    kill $GO_PIDS 2>/dev/null
    sleep 2
    
    # Force kill if still running
    REMAINING_GO_PIDS=$(ps aux | grep "go run.*main.go" | grep -v grep | awk '{print $2}' 2>/dev/null)
    if [ -n "$REMAINING_GO_PIDS" ]; then
        echo "Force killing remaining Go processes: $REMAINING_GO_PIDS"
        kill -9 $REMAINING_GO_PIDS 2>/dev/null
    fi
    echo -e "${GREEN}✅ Go processes stopped${NC}"
else
    echo -e "${GREEN}✅ No Go processes running main.go${NC}"
fi

# 3. Stop and remove Qdrant container
echo ""
echo "🔍 Checking for Qdrant containers..."

if command_exists docker; then
    # Check if qdrant-rag container exists and is running
    if docker ps | grep -q "qdrant-rag"; then
        echo -e "${YELLOW}Stopping Qdrant container...${NC}"
        docker stop qdrant-rag >/dev/null 2>&1
        echo -e "${GREEN}✅ Qdrant container stopped${NC}"
    else
        echo -e "${GREEN}✅ Qdrant container not running${NC}"
    fi
    
    # Check if qdrant-rag container exists (stopped)
    if docker ps -a | grep -q "qdrant-rag"; then
        echo -e "${YELLOW}Removing Qdrant container...${NC}"
        docker rm qdrant-rag >/dev/null 2>&1
        echo -e "${GREEN}✅ Qdrant container removed${NC}"
    else
        echo -e "${GREEN}✅ No Qdrant container to remove${NC}"
    fi
    
    # Check for any other qdrant containers
    OTHER_QDRANT=$(docker ps -a | grep qdrant | grep -v "qdrant-rag" | awk '{print $1}' 2>/dev/null)
    if [ -n "$OTHER_QDRANT" ]; then
        echo -e "${YELLOW}Found other Qdrant containers: $OTHER_QDRANT${NC}"
        echo "Stopping and removing them..."
        docker stop $OTHER_QDRANT >/dev/null 2>&1
        docker rm $OTHER_QDRANT >/dev/null 2>&1
        echo -e "${GREEN}✅ Other Qdrant containers cleaned up${NC}"
    fi
else
    echo -e "${RED}⚠️  Docker not found, skipping container cleanup${NC}"
fi

# 4. Verify ports are free
echo ""
echo "🔍 Verifying ports are free..."

# Check port 8080 (Go API)
if port_in_use 8080; then
    echo -e "${RED}⚠️  Port 8080 still in use${NC}"
    lsof -i :8080
else
    echo -e "${GREEN}✅ Port 8080 is free${NC}"
fi

# Check port 6333 (Qdrant HTTP)
if port_in_use 6333; then
    echo -e "${RED}⚠️  Port 6333 still in use${NC}"
    lsof -i :6333
else
    echo -e "${GREEN}✅ Port 6333 is free${NC}"
fi

# Check port 6334 (Qdrant gRPC)
if port_in_use 6334; then
    echo -e "${RED}⚠️  Port 6334 still in use${NC}"
    lsof -i :6334
else
    echo -e "${GREEN}✅ Port 6334 is free${NC}"
fi

echo ""
echo "🎉 Cleanup complete!"
echo ""
echo "📝 Summary:"
echo "   • Go applications stopped"
echo "   • Qdrant containers removed"
echo "   • Ports 8080, 6333, 6334 freed"
echo ""
echo "🚀 To restart the project:"
echo "   1. ./scripts/setup-qdrant.sh"
echo "   2. go run cmd/main.go"
echo ""
