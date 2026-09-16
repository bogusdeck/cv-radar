# Testing CV-RADAR

This guide details how to test the ATS Resume Scanner & AI CV Optimizer components.

## Prerequisites

- **Go 1.21+** (API server & TUI dashboard)
- **Node.js 18+** (Web dashboard)
- **Tectonic** (Optional - LaTeX to PDF compilation)

---

## 1. Quick Full Stack Launch

Start both the Go API server (`:8085`) and the React Web UI (`:3000`):

```bash
./start.sh
```

Open `http://localhost:5173` (or `http://localhost:3000`) in your browser to test document uploading, ATS platform analysis, and PDF compilation.

---

## 2. Testing the Go TUI Dashboard

Build and launch the interactive Bubble Tea terminal dashboard:

```bash
go run ./cmd/tui
# or run the precompiled binary:
./cv-tui
```

### Features to test in TUI:
- Platform selection via Arrow Keys (Workday, Taleo, iCIMS, Greenhouse)
- Score calculation and letter grade (A–F)
- Matched and missing keyword lists
- Pressing `O` to trigger headless AI optimization

---

## 3. Testing the Go API Endpoints Directly

Start the server:
```bash
go run ./cmd/server
```

### Test Platforms Listing:
```bash
curl -s http://localhost:8085/api/platforms | jq .
```

### Test CV Analysis:
```bash
curl -X POST http://localhost:8085/api/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "cv_text": "Software Engineer experienced in Go, React, PostgreSQL, Docker, AWS",
    "jd_text": "Looking for Senior Engineer with Go, React, Kubernetes, AWS experience",
    "platform": "Workday"
  }' | jq .
```

---

## 4. Testing Headless AI Optimization

Run the headless AI optimization script with your choice of AI CLI backend (`claude`, `agy`, or `opencode`):

```bash
./optimize.sh cv.tex jd.txt Workday agy
```

Output will be generated in `optimized_cv.md`.

---

## 5. Testing Automated Setup Script

Test the single-command setup script:

```bash
./install.sh
```