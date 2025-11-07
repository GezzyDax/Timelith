#!/bin/bash
# Run all tests locally
set -e

echo "🧪 Running all tests..."
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

FAILED=0

# Go tests
echo "=== Running Go Backend Tests ==="
cd go-backend
if go test -v -race -coverprofile=coverage.out ./...; then
    echo -e "${GREEN}✓ Go tests passed${NC}"
    echo ""
    echo "Coverage report:"
    go tool cover -func=coverage.out | tail -n 1
else
    echo -e "${RED}✗ Go tests failed${NC}"
    FAILED=1
fi
cd ..

echo ""
echo "=== Testing Docker Builds ==="

# Build Go backend image
echo "Building Go backend image..."
if docker build -t timelith-backend:test ./go-backend; then
    echo -e "${GREEN}✓ Backend image built${NC}"
else
    echo -e "${RED}✗ Backend image build failed${NC}"
    FAILED=1
fi

# Build Web UI image
echo ""
echo "Building Web UI image..."
if docker build -t timelith-web:test ./web-ui; then
    echo -e "${GREEN}✓ Web UI image built${NC}"
else
    echo -e "${RED}✗ Web UI image build failed${NC}"
    FAILED=1
fi

# Integration test
echo ""
echo "=== Running Integration Test ==="
# Create test .env
cat << EOF > .env.test
POSTGRES_PASSWORD=test_password
TELEGRAM_APP_ID=12345
TELEGRAM_APP_HASH=test_hash
JWT_SECRET=test-jwt-secret-key
ENCRYPTION_KEY=test-encryption-key-32-bytes-!
ENVIRONMENT=test
NEXT_PUBLIC_API_URL=http://localhost:8080
EOF

# Start services
docker compose --env-file .env.test up -d postgres redis
sleep 5

# Check health
if docker compose exec -T postgres pg_isready -U timelith > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Integration test passed${NC}"
else
    echo -e "${RED}✗ Integration test failed${NC}"
    FAILED=1
fi

# Cleanup
docker compose down -v
rm -f .env.test

echo ""
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    exit 1
fi
