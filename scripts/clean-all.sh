#!/bin/bash
# Clean all build artifacts and caches
set -e

echo "🧹 Cleaning all build artifacts and caches..."
echo ""

# Go backend
echo "=== Cleaning Go backend ==="
cd go-backend
if [ -d "bin" ]; then
    rm -rf bin/
    echo "✓ Removed bin/"
fi
go clean -cache -testcache -modcache
echo "✓ Cleaned Go caches"
cd ..

# Web UI
echo ""
echo "=== Cleaning Web UI ==="
cd web-ui
if [ -d ".next" ]; then
    rm -rf .next/
    echo "✓ Removed .next/"
fi
if [ -d "node_modules" ]; then
    rm -rf node_modules/
    echo "✓ Removed node_modules/"
fi
if [ -f "package-lock.json" ]; then
    rm package-lock.json
    echo "✓ Removed package-lock.json"
fi
cd ..

# Docker
echo ""
echo "=== Cleaning Docker ==="
echo "Stopping containers..."
docker compose down -v 2>/dev/null || true
echo "✓ Stopped and removed containers"

# Remove dangling images
echo "Removing dangling images..."
docker image prune -f
echo "✓ Cleaned dangling images"

echo ""
echo "✅ Clean complete!"
echo ""
echo "To rebuild everything:"
echo "  1. cd web-ui && npm install"
echo "  2. cd go-backend && go mod download"
echo "  3. docker compose build"
