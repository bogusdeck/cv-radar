#!/bin/bash
# Start the Go API server and the TUI (using npx ts-node)

echo "🚀 Starting CV-RADAR API and TUI..."

# Start Go API server in background
echo "📡 Starting Go API on :8080"
go run ./cmd/server/main.go &
GO_PID=$!

# Wait for server to be ready (max 30 seconds)
echo "⏳ Waiting for API server to start..."
SERVER_READY=false
for i in {1..30}; do
  if curl -sf http://localhost:8080/api/platforms > /dev/null; then
    echo "🌐 API server is ready"
    SERVER_READY=true
    break
  fi
  echo "⏳ Waiting for API server to start... ($i/30)"
  sleep 1
done

if [ "$SERVER_READY" = false ]; then
  echo "⚠️  API server did not become ready in time. TUI may not work correctly."
fi

# Start TUI using npx tsx
echo "💻 Starting TUI"
cd tui && npx tsx src/index.tsx "$@"

# When TUI exits, kill the server
echo ""
echo "🛑 Stopping API server..."
kill $GO_PID 2>/dev/null
wait $GO_PID 2>/dev/null
echo "🛑 TUI and API stopped."
