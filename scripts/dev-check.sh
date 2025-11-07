#!/bin/bash
# Local development check script - runs all checks before committing
set -e

echo "🔍 Running development checks..."
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  .env file not found. Creating from .env.example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ Created .env file. Please update it with your credentials.${NC}"
fi

echo "=== Go Backend Checks ==="
cd go-backend

# Format check
echo "📝 Checking Go formatting..."
if ! gofmt -l . | grep -q .; then
    echo -e "${GREEN}✓ Go formatting is correct${NC}"
else
    echo -e "${RED}✗ Go formatting issues found:${NC}"
    gofmt -l .
    echo "Run: cd go-backend && gofmt -w ."
    exit 1
fi

# Build check
echo "🔨 Building Go backend..."
if go build -o bin/server ./cmd/server > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

# Run tests
echo "🧪 Running Go tests..."
if go test -race ./... > /dev/null 2>&1; then
    echo -e "${GREEN}✓ All tests passed${NC}"
else
    echo -e "${RED}✗ Tests failed${NC}"
    go test -v -race ./...
    exit 1
fi

# Go mod tidy
echo "📦 Checking Go modules..."
go mod tidy
if git diff --exit-code go.mod go.sum > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Go modules are tidy${NC}"
else
    echo -e "${YELLOW}⚠️  Go modules were updated${NC}"
fi

cd ..

echo ""
echo "=== Web UI Checks ==="
cd web-ui

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "📦 Installing npm dependencies..."
    npm ci
fi

# Lint check
echo "📝 Running ESLint..."
if npm run lint > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Linting passed${NC}"
else
    echo -e "${RED}✗ Linting failed${NC}"
    npm run lint
    exit 1
fi

# Type check
echo "🔍 Running TypeScript type check..."
if npx tsc --noEmit; then
    echo -e "${GREEN}✓ Type check passed${NC}"
else
    echo -e "${RED}✗ Type check failed${NC}"
    exit 1
fi

# Build check
echo "🔨 Building Next.js app..."
if NEXT_PUBLIC_API_URL=http://localhost:8080 npm run build > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

cd ..

echo ""
echo "=== Docker Checks ==="
# Check if docker-compose.yml is valid
if docker compose config > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Docker Compose config is valid${NC}"
else
    echo -e "${RED}✗ Docker Compose config is invalid${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}✅ All checks passed! You're ready to commit.${NC}"
