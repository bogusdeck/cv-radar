#!/bin/bash
# Start both the Go API server and the React web frontend

echo "🚀 Starting ATS Scanner..."

# Start Go API server in background
echo "📡 Starting Go API on :8080"
cd "$(dirname "$0")"
go run ./cmd/server/main.go &
GO_PID=$!

# Start React dev server
echo "🌐 Starting React web UI on :3000"
cd web && npm run dev &
WEB_PID=$!

echo ""
echo "✅ Both servers running:"
echo "   → API:    http://localhost:8080"
echo "   → Web UI: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop both."

# Wait and cleanup on Ctrl+C
trap "kill $GO_PID $WEB_PID 2>/dev/null; echo 'Stopped.'" SIGINT SIGTERM
wait
