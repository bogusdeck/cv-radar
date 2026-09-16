<p align="center"><pre>
  ██████╗██╗   ██╗      ██████╗  █████╗ ██████╗  █████╗ ██████╗ 
 ██╔════╝██║   ██║      ██╔══██╗██╔══██╗██╔══██╗██╔══██╗██╔══██╗
 ██║     ██║   ██║█████╗██████╔╝███████║██║  ██║███████║██████╔╝
 ██║     ╚██╗ ██╔╝╚════╝██╔══██╗██╔══██║██║  ██║██╔══██║██╔══██╗
 ╚██████╗ ╚████╔╝       ██║  ██║██║  ██║██████╔╝██║  ██║██║  ██║
  ╚═════╝  ╚═══╝        ╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝
</pre></p>

<p align="center"><strong>ATS Resume Scanner & AI CV Optimizer (CLI & Agent Skill)</strong></p>


<p align="center">
  Score your CV against a Job Description. Find missing keywords. Auto-rewrite it for Workday, Taleo, Greenhouse, or iCIMS — headlessly, right inside your AI coding CLI.
</p>

<p align="center">
  <a href="https://github.com/bogusdeck/cv-radar/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go" alt="Go">
  <img src="https://img.shields.io/badge/works_with-Claude_Code-000?logo=anthropic" alt="Claude Code">
  <img src="https://img.shields.io/badge/works_with-Codex-412991?logo=openai" alt="Codex">
  <img src="https://img.shields.io/badge/works_with-Antigravity-FF75B5" alt="Antigravity">
</p>

---

## Quick Start — One Command

```bash
curl -fsSL https://raw.githubusercontent.com/bogusdeck/cv-radar/main/install.sh | bash
```

Or manually:

```bash
git clone https://github.com/bogusdeck/cv-radar.git
cd cv-radar
./install.sh
```

Then open your AI CLI in **any** folder:

```bash
agy         # Antigravity
claude      # Claude Code
opencode    # OpenCode
codex       # Codex
```

And say:

```
/cv-optimizer
```

The agent reads the built-in skill and walks you through the full workflow.

---

## What It Does

| Feature | Description |
|---|---|
| **ATS Scoring Engine** | Offline Go engine scoring your CV against a JD for Workday, Taleo, Greenhouse, and iCIMS |
| **Missing Keyword Detection** | Finds every required JD keyword missing from your CV |
| **TUI Dashboard** | Interactive terminal UI (`cv-tui`) to browse your score, grade, and recommendations |
| **AI Skill** | Global skill auto-discovered by Claude Code, Codex, and Antigravity — just say `/cv-optimizer` |
| **Headless AI Optimization** | Dispatches `claude -p` / `agy -p` / `opencode run` in the background to rewrite your CV |

---

## Installation

### Prerequisites

| Tool | Required | Notes |
|---|---|---|
| `go` 1.21+ | ✅ Yes | [Download](https://go.dev/dl) |
| `git` | ✅ Yes | [Download](https://git-scm.com) |
| `tectonic` | For PDF compilation | [Install](https://tectonic-typesetting.github.io/book/latest/installation) |

---

## Usage

### 1. Inside your AI Coding CLI (Recommended)

Open any folder in Claude Code, Antigravity, Codex, or OpenCode:

```bash
agy        # or: claude / opencode / codex
```

Then just say:

```
/cv-optimizer
```

The AI will:
- Read your CV and Job Description
- Score your CV against the target ATS
- Show missing keywords
- Rewrite your CV headlessly using the AI backend
- Output a compilable LaTeX file

---

### 2. TUI Dashboard

Run from anywhere in your terminal:

```bash
cv-tui
```

- Arrow keys to select ATS platform
- See your score, grade, and missing keywords
- Press `O` to trigger headless AI optimization without leaving the TUI

---

### 3. Headless AI Optimization (CLI)

```bash
# Using Claude Code
./optimize.sh cv.tex jd.txt Workday claude

# Using Antigravity
./optimize.sh cv.tex jd.txt Workday agy

# Using OpenCode
./optimize.sh cv.tex jd.txt Workday opencode

# Using Codex
./optimize.sh cv.tex jd.txt Workday codex
```

The optimized LaTeX is saved to `optimized_cv.md`, ready to compile via `tectonic`.

---

## ATS Platform Support

| Platform | Scoring Strategy |
|---|---|
| **Workday** | Strict exact-keyword matching. Penalizes missing required skills heavily. |
| **Taleo** | Keyword density scoring. Repeating JD keywords boosts score. |
| **Greenhouse** | Semantic AI + human review. Rewards natural phrasing and measurable impact. |
| **iCIMS** | Mixed exact + semantic. Title alignment matters most. |

---

## Project Structure

```
.agents/skills/cv-optimizer/   ← AI skill (auto-discovered by Antigravity)
.claude/skills/cv-optimizer/   ← Claude Code skill
.codex-plugin/plugin.json      ← Codex plugin registration
cmd/tui/                       ← Go Bubbletea TUI dashboard (`cv-tui`)
internal/ats/                  ← ATS scoring engines
internal/parser/               ← CV/PDF/LaTeX parser
internal/models/               ← Shared data models
optimize.sh                    ← Headless AI dispatch script
install.sh                     ← One-command installer
```

---

## Built With

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Bubble Tea](https://img.shields.io/badge/Bubble_Tea-FF75B5?style=flat&logo=go&logoColor=white)
![Tectonic](https://img.shields.io/badge/Tectonic-LaTeX-008080?style=flat)

---

## Author

Built by [Tanish Vashisth](https://github.com/bogusdeck)
