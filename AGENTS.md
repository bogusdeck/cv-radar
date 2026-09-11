# CV Optimizer — ATS Resume Scanner

This repository contains a full ATS CV scanning and optimization pipeline.

## Quick Start

When the user says `/cv-optimizer`, `optimize my CV`, or `scan my resume`:
1. Read the `cv-optimizer` skill from `.agents/skills/cv-optimizer/SKILL.md`
2. Follow the step-by-step workflow defined there

## Headless AI Commands

This project supports headless execution through any AI CLI:

```bash
./optimize.sh cv.tex jd.txt Workday claude     # Claude Code
./optimize.sh cv.tex jd.txt Workday agy        # Antigravity
./optimize.sh cv.tex jd.txt Workday opencode   # OpenCode
```

## TUI Dashboard

```bash
./cv-tui    # Launch the interactive ATS score dashboard
```

## Go API Server

```bash
./start.sh  # Start the full stack (API on :8085, web on :5173)
```
