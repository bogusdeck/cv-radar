---
name: cv-optimizer
description: >-
  ATS Resume Scanner & AI CV Optimizer — scan CV, score against job description,
  launch interactive TUI dashboard, and auto-rewrite CVs for Workday, Taleo,
  Greenhouse, and iCIMS.
arguments: mode
user_invocable: true
user-invocable: true
argument-hint: "[tui | scan | optimize | pdf | {JD text or file}]"
license: MIT
---

# CV Optimizer — Router & Command Center

CV Optimizer is a multi-CLI ATS resume scanning and optimization tool.

## Mode Routing

Determine the mode from `$mode` / `$ARGUMENTS`:

| Input | Mode | Action |
|-------|------|--------|
| (empty / no args) | `menu` | Show interactive command menu |
| `tui` | `tui` | Launch interactive terminal UI (`cv-tui open`) |
| `scan` | `scan` | Scan & score `cv.tex` against `jd.txt` in chat |
| `optimize` | `optimize` | Rewrite CV using `./optimize.sh` |
| `pdf` | `pdf` | Compile `optimized_cv.md` to PDF via `tectonic` |
| JD text / URL | **`auto-pipeline`** | Full pipeline: audit → score → optimize → PDF |

---

## Discovery Menu (no arguments)

If no arguments are supplied, show this menu in chat:

```
CV Optimizer -- ATS Resume Scanner & AI Optimizer

Available commands:
  /cv-optimizer tui        → Launch visual interactive TUI dashboard in Terminal window
  /cv-optimizer scan       → Run ATS keyword & layout audit in chat
  /cv-optimizer optimize   → Headless AI rewrite (Workday, Taleo, Greenhouse, iCIMS)
  /cv-optimizer pdf        → Compile optimized resume to PDF (tectonic)
  /cv-optimizer {JD/URL}   → Run FULL pipeline (audit + score + rewrite + PDF)
```

---

## Step 1 — TUI Dashboard Mode (`/cv-optimizer tui`)

When `tui` mode is selected or requested, execute:

```bash
osascript -e 'tell application "Terminal" to do script "cv-tui"' 2>/dev/null || cv-tui open || cv-tui
```

---

## Step 2 — Scan & Audit Mode (`/cv-optimizer scan` or JD input)

1. Read `cv.tex` (or `cv_test.txt` / workspace CV) and `jd.txt` (or workspace JD).
2. Extract hard skills, required experience, and job title keywords.
3. Compare CV vs JD across platforms:
   - **Workday**: Strict exact string keyword matches.
   - **Taleo**: Keyword frequency & density scoring.
   - **Greenhouse**: Impact metrics & structured achievements.
   - **iCIMS**: Direct title and section alignment.
4. Output a summary score (0-100), grade (A-F), matched keywords, and missing keywords in clean Markdown tables.

---

## Step 3 — Headless AI Rewrite (`/cv-optimizer optimize`)

Run the optimization script:

```bash
./optimize.sh cv.tex jd.txt Workday opencode
```

Save output to `optimized_cv.md`.

---

## Step 4 — PDF Compilation (`/cv-optimizer pdf`)

Compile `optimized_cv.md` via `tectonic`:

```bash
tectonic optimized_cv.md -o output/Optimized_Resume.pdf
```

---

## LaTeX Formatting Rules (CRITICAL)

- Preserve exact LaTeX commands (`\resumeSubheading`, `\resumeItem`).
- ALL `\resumeSubheading` must be inside `\resumeSubHeadingListStart` / `\resumeSubHeadingListEnd`.
- ALL `\resumeItem` must be inside `\resumeItemListStart` / `\resumeItemListEnd`.
- Escape `%` and `$` in text body.
