# How to Test CV-RADAR

## Prerequisites
- Go installed (for API server)
- Node.js installed (for TUI and Web UI)
- API server must be running on port 8080 for TUI/Web to function

## Option 1: Test Web Interface (Recommended for Reliable Testing)

### Step 1: Start the API Server
In Terminal 1:
```bash
go run ./cmd/server/main.go
```
You should see: `🚀 CV-RADAR API running on http://localhost:8080`

### Step 2: Start the Web UI
In Terminal 2:
```bash
cd web && npm run dev
```
Then open: **http://localhost:3000**

The web interface provides a full-featured CV/JD analyzer with platform selection, score visualization, and detailed breakdowns.

### Step 3: Stop Services
- Terminal 1 (API): Press Ctrl+C
- Terminal 2 (Web UI): Press Ctrl+C

## Option 2: Test TUI Interface (Static CLI Mode)

Due to ES module compatibility issues with the 'ink' package in newer Node.js versions, the interactive TUI may not start reliably. However, the TUI includes a static CLI mode that works without ink.

### Step 1: Start the API Server
In Terminal 1:
```bash
go run ./cmd/server/main.go
```

### Step 2: Use TUI in Static Mode
In Terminal 2:
```bash
# Example: analyze files
cd tui && node dist/index.js ./path/to/cv.txt ./path/to/jd.txt

# Example: pipe CV and pass JD as argument
cat ./path/to/cv.txt | node dist/index.js --stdin "job description text"

# Example: CV file + inline JD
node dist/index.js --cv ./path/to/cv.txt --jd "job description text"
```

The static mode will output a formatted report to the terminal.

### Step 3: Stop API Server
Terminal 1: Press Ctrl+C

## Option 3: Using Startup Scripts

### For Web + API (Recommended):
```bash
./start-web.sh
```
This starts both API and Web UI in background (Ctrl+C to stop both).

### Original Script (API + Web):
```bash
./start.sh
```

## Manual Testing Steps

1. **Start API Server:**
   ```bash
   go run ./cmd/server/main.go
   ```

2. **In another terminal, test with sample data:**
   ```bash
   # Create test CV and JD files
   echo "John Doe
   Email: john@example.com
   Phone: 555-1234
   
   Experience:
   - Software Engineer at Tech Corp (2020-Present)
     • Developed web applications using React and Node.js
     • Worked with PostgreSQL and AWS
   
   Skills: JavaScript, React, Node.js, PostgreSQL, AWS, Git
   
   Education:
   - BS Computer Science, University Example, 2020" > test_cv.txt
   
   cp sample_jd.txt test_jd.txt
   
   # Test API directly
   curl -X POST http://localhost:8080/api/analyze \
     -H "Content-Type: application/json" \
     -d "{\"cv_text\":\"$(cat test_cv.txt)\",\"jd_text\":\"$(cat test_jd.txt)\"}"
   ```

## Troubleshooting

### TUI Interactive Mode Issues
If you attempt to run `cd tui && npm run dev` or `cd tui && ts-node src/index.tsx` and see errors like:
```
Error [ERR_REQUIRE_ASYNC_MODULE]: require() cannot be used on an ESM graph with top-level await
```
This is due to the 'ink' package being an ES module and Node.js version incompatibility. Use the web interface or static CLI mode instead.

### API Server Not Responding
- Ensure you ran `go run ./cmd/server/main.go` and see the startup message
- Check if port 8080 is free: `lsof -i :8080`
- Try explicitly: `go run ./cmd/server/main.go -port=8081` then use `ATS_API=http://localhost:8081` env var

## File Structure
- API Server: `./cmd/server/main.go`
- TUI Source: `./tui/src/index.tsx`  
- Web Source: `./web/src/`
- Shared Code: `./internal/` (parsers, ATS engines, models)
- Startup Scripts: `start.sh`, `start-tui.sh`, `start-web.sh`
- Documentation: `README.md`, `TESTING.md` (this file)

## Sample Data
- `sample_jd.txt`: Contains a sample job description for testing
- You can create your own CV text files to test with

## Production Builds
For distribution:
- TUI: `cd tui && npm run build` → outputs to `dist/`
- Web: `cd web && npm run build` → outputs to `dist/` 
- API: `go build -o ats-server ./cmd/server/` → produces standalone binary

Given the ES module challenges with the TUI's interactive mode, the web interface provides the most reliable and feature-rich testing experience. The TUI's static CLI mode remains functional for quick terminal-based analysis without the interactive dashboard.