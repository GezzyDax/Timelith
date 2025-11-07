#!/bin/bash
# Quick start script for local development
set -e

echo "🚀 Starting Timelith development environment..."
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "⚠️  Please update .env with your credentials before running the app!"
    echo ""
    read -p "Press enter to continue..."
fi

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check required tools
echo "🔍 Checking required tools..."
MISSING_TOOLS=()

if ! command_exists docker; then
    MISSING_TOOLS+=("docker")
fi

if ! command_exists go; then
    MISSING_TOOLS+=("go")
fi

if ! command_exists node; then
    MISSING_TOOLS+=("node")
fi

if [ ${#MISSING_TOOLS[@]} -ne 0 ]; then
    echo "❌ Missing required tools: ${MISSING_TOOLS[*]}"
    echo "Please install them before continuing."
    exit 1
fi

echo "✅ All required tools are installed"
echo ""

# Start infrastructure
echo "🐳 Starting infrastructure (PostgreSQL & Redis)..."
docker compose up -d postgres redis

# Wait for services to be healthy
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check if services are healthy
if docker compose ps postgres | grep -q "healthy"; then
    echo "✅ PostgreSQL is ready"
else
    echo "⚠️  PostgreSQL is not healthy yet, waiting..."
    sleep 5
fi

if docker compose ps redis | grep -q "healthy"; then
    echo "✅ Redis is ready"
else
    echo "⚠️  Redis is not healthy yet, waiting..."
    sleep 5
fi

echo ""
echo "📊 Services Status:"
docker compose ps

echo ""
echo "================================================"
echo "✅ Infrastructure is ready!"
echo ""
echo "Next steps:"
echo "  1. Backend:  cd go-backend && go run ./cmd/server"
echo "  2. Frontend: cd web-ui && npm run dev"
echo ""
echo "Or run everything with Docker:"
echo "  docker compose up"
echo ""
echo "To stop infrastructure:"
echo "  docker compose down"
echo "================================================"
