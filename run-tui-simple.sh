#!/bin/bash
# Simple script to run TUI - assumes API is already running

echo "🚀 Starting CV-RADAR TUI (simple version)..."
echo "Make sure the API server is running on http://localhost:8080"
echo "You can start it with: go run ./cmd/server/main.go"
echo ""

# Build TUI
echo "📦 Building TUI..."
cd tui && npm run build

# Check if build succeeded
if [ $? -ne 0 ]; then
  echo "❌ Build failed!"
  exit 1
fi

echo "✅ Build successful!"
echo ""

# Try to run with node - if it fails due to ES module issues, suggest alternative
echo "💻 Attempting to start TUI..."
cd tui && node dist/index.js "$@"

# If that fails, provide guidance
if [ $? -ne 0 ]; then
  echo ""
  echo "⚠️  TUI failed to start due to ES module compatibility issues."
  echo "This is a known issue with the 'ink' package and Node.js versions."
  echo ""
  echo "Workarounds:"
  echo "1. Use an older Node.js version (v18 or v20)"
  echo "2. Run the TUI in development mode: cd tui && npm run dev"
  echo "3. Use the static CLI mode: cd tui && node dist/index.js --help"
  echo ""
  echo "For testing, you can also use the web interface:"
  echo "   cd web && npm run dev"
  echo "   Then open http://localhost:3000"
fi