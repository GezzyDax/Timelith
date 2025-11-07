#!/bin/bash

echo "🛑 Stopping Timelith..."

docker compose down

echo "✅ Timelith stopped."
echo ""
echo "💡 To remove all data: docker compose down -v"
