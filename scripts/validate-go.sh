#!/bin/bash
# Simple Go code validation script

echo "🔍 Checking Go code..."
cd /home/user/Timelith/go-backend

echo ""
echo "✅ Go files found:"
find . -name "*.go" -type f | wc -l

echo ""
echo "📦 Checking go.mod..."
if [ -f go.mod ]; then
    echo "✅ go.mod exists"
    grep "^module" go.mod
else
    echo "❌ go.mod not found"
    exit 1
fi

echo ""
echo "🔧 Validating Go syntax (gofmt)..."
UNFORMATTED=$(gofmt -l .)
if [ -z "$UNFORMATTED" ]; then
    echo "✅ All Go files are properly formatted"
else
    echo "⚠️  These files need formatting:"
    echo "$UNFORMATTED"
fi

echo ""
echo "📝 Checking for common issues..."
echo "Duplicate imports:"
grep -r "import (" . --include="*.go" -A 10 | grep -E "^\s+\".*\"$" | sort | uniq -d || echo "✅ No duplicate imports found"

echo ""
echo "✅ Basic validation complete!"
echo ""
echo "To build in Docker, run:"
echo "  docker-compose build go-backend"
