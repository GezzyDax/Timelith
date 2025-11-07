#!/bin/bash

echo "🚀 Starting Timelith..."

# Check if .env exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found!"
    echo "📝 Creating .env from .env.example..."
    cp .env.example .env
    echo "✅ .env created. Please edit it with your credentials before continuing."
    exit 1
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Start services
echo "🐳 Starting Docker containers..."
docker compose up -d

# Wait for services to be healthy
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check service status
echo "📊 Service Status:"
docker compose ps

echo ""
echo "✅ Timelith is running!"
echo ""
echo "🌐 Web Dashboard: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"
echo ""
echo "📖 To view logs: docker compose logs -f"
echo "🛑 To stop: docker compose down"
