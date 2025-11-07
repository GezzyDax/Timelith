#!/bin/bash
# Pre-commit check script - ensures code quality before committing
# Run this manually before git commit: make pre-commit
set -e

echo "╔══════════════════════════════════════════╗"
echo "║   🔍 Pre-Commit Quality Checks          ║"
echo "╔══════════════════════════════════════════╝"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Track failures
FAILED=0
TOTAL_CHECKS=0
PASSED_CHECKS=0

# Function to run a check
run_check() {
    local check_name="$1"
    local check_command="$2"
    local dir="${3:-.}"

    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))

    echo -e "${BLUE}┌─ ${check_name}${NC}"

    if (cd "$dir" && eval "$check_command") > /tmp/check_output 2>&1; then
        echo -e "${GREEN}└─ ✓ PASSED${NC}"
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
        echo ""
        return 0
    else
        echo -e "${RED}└─ ✗ FAILED${NC}"
        echo -e "${YELLOW}Output:${NC}"
        cat /tmp/check_output
        echo ""
        FAILED=1
        return 1
    fi
}

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  Warning: .env file not found${NC}"
    echo -e "Creating from .env.example..."
    cp .env.example .env
    echo -e "${GREEN}✓ Created .env file${NC}"
    echo ""
fi

echo -e "${BOLD}═══ Go Backend Checks ═══${NC}"
echo ""

# Go format check
run_check "Go Code Formatting" \
    "! gofmt -l . | grep -q ." \
    "go-backend"

# Go build
run_check "Go Build" \
    "go build -o bin/server ./cmd/server" \
    "go-backend"

# Go tests with race detection
run_check "Go Tests (with race detection)" \
    "go test -race -short ./..." \
    "go-backend"

# Go mod tidy
run_check "Go Modules Tidy" \
    "go mod tidy && git diff --exit-code go.mod go.sum" \
    "go-backend"

# Go lint (if golangci-lint is available)
if command -v golangci-lint &> /dev/null; then
    run_check "Go Linting (golangci-lint)" \
        "golangci-lint run --timeout=5m" \
        "go-backend"
else
    echo -e "${YELLOW}⚠️  golangci-lint not installed, skipping linting${NC}"
    echo -e "   Install: https://golangci-lint.run/usage/install/"
    echo ""
fi

echo -e "${BOLD}═══ Web UI Checks ═══${NC}"
echo ""

# Check if node_modules exists
if [ ! -d "web-ui/node_modules" ]; then
    echo -e "${YELLOW}⚠️  node_modules not found, installing...${NC}"
    cd web-ui && npm ci && cd ..
    echo ""
fi

# TypeScript type check
run_check "TypeScript Type Check" \
    "npx tsc --noEmit" \
    "web-ui"

# ESLint
run_check "ESLint" \
    "npm run lint" \
    "web-ui"

# Next.js build
run_check "Next.js Build" \
    "NEXT_PUBLIC_API_URL=http://localhost:8080 npm run build" \
    "web-ui"

echo -e "${BOLD}═══ Docker Checks ═══${NC}"
echo ""

# Docker Compose config validation
run_check "Docker Compose Config" \
    "docker compose config > /dev/null"

# Summary
echo ""
echo "╔══════════════════════════════════════════╗"
echo -e "║   ${BOLD}Summary${NC}                              ║"
echo "╠══════════════════════════════════════════╣"
echo -e "║  Total checks:   ${BOLD}$TOTAL_CHECKS${NC}                       ║"
echo -e "║  Passed:         ${GREEN}${BOLD}$PASSED_CHECKS${NC}                       ║"
echo -e "║  Failed:         ${RED}${BOLD}$((TOTAL_CHECKS - PASSED_CHECKS))${NC}                       ║"
echo "╚══════════════════════════════════════════╝"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}${BOLD}✅ All checks passed! Ready to commit.${NC}"
    echo ""
    echo "You can now safely run:"
    echo -e "  ${BLUE}git add .${NC}"
    echo -e "  ${BLUE}git commit -m \"your message\"${NC}"
    echo ""
    exit 0
else
    echo -e "${RED}${BOLD}❌ Some checks failed!${NC}"
    echo ""
    echo "Please fix the issues above before committing."
    echo "Run individual checks:"
    echo -e "  ${BLUE}make backend-test${NC}   - Run Go tests"
    echo -e "  ${BLUE}make backend-lint${NC}   - Run Go linter"
    echo -e "  ${BLUE}make web-lint${NC}       - Run ESLint"
    echo -e "  ${BLUE}make web-build${NC}      - Build Next.js"
    echo ""
    exit 1
fi
