# CV-RADAR

> Analyze your CV against 6 real ATS platforms with platform-specific scoring algorithms.

## Supported Platforms

| Platform | Vendor | Strategy | Key Behavior |
|---|---|---|---|
| Workday | Workday | Exact + HiredScore AI | Strict parser, penalizes creative formats |
| Taleo | Oracle | Literal exact match | Strictest, auto-reject via Req Rank |
| iCIMS | iCIMS | Semantic (ML-based) | Most forgiving, grammar-based NLP |
| Greenhouse | Greenhouse | Semantic (LLM-based) | Human review, no auto-scoring |
| Lever | Employ | Stemming-based | Abbreviation-blind, search-dependent |
| SuccessFactors | SAP | Taxonomy normalization | Joule AI skills matching |

## Quick Start

```bash
# Clone and start everything
chmod +x start.sh
./start.sh
```

Then open **http://localhost:3000**

## Manual Start

**Go API server (port 8080):**
```bash
go run ./cmd/server/main.go
```

**React web UI (port 3000):**
```bash
cd web && npm run dev
```

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/platforms` | List all ATS platforms |
| POST | `/api/analyze` | Analyze CV vs JD (JSON body) |
| POST | `/api/parse` | Parse CV to structured JSON |
| POST | `/api/upload` | Upload PDF/.tex file → extract text |

### Analyze Request
```json
{
  "cv_text": "...",
  "jd_text": "...",
  "platform": ""  // empty = all platforms
}
```

## Project Structure

```
cv-radar/
├── cmd/server/          # Go HTTP API
├── internal/
│   ├── parser/          # PDF, LaTeX, plain text parsers
│   ├── ats/             # 6 ATS scoring engines
│   └── models/          # Shared data models
├── web/                 # React + Vite web UI
└── start.sh             # Start everything
```

## Tech Stack
- **Backend**: Go + Gin
- **Frontend**: React + TypeScript + Vite
- **CV Parsing**: PDF (ledongthuc/pdf), LaTeX (custom stripper), plain text
- **NLP**: TF-IDF cosine similarity, Porter stemmer, skills taxonomy normalization
