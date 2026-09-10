#!/bin/bash
# Start the Go API server and the React web frontend

echo "🚀 Starting CV-RADAR API and Web UI..."

# Start Go API server in background
echo "📡 Starting Go API on :8080"
go run ./cmd/server/main.go &
GO_PID=$!

# Wait for server to be ready (max 10 seconds)
echo "⏳ Waiting for API server to start..."
SERVER_READY=false
for i in {1..10}; do
  if curl -s http://localhost:8080/api/platforms > /dev/null; then
    echo "🌐 API server is ready"
    SERVER_READY=true
    break
  fi
  echo "⏳ Waiting for API server to start... ($i/10)"
  sleep 1
done

if [ "$SERVER_READY" = false ]; then
  echo "⚠️  API server did not become ready in time. Web UI may not work correctly."
fi

# Start React dev server in background
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