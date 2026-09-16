---
name: cv-optimizer
description: >-
  Use this skill when the user says "cv-optimizer", "/cv-optimizer", "optimize my CV", or "scan my resume".
  Launches the interactive cv-tui dashboard and ATS optimization workflow. Scores the CV against the job description,
  identifies missing keywords, then optionally rewrites the CV using the headless AI backend.
  Supports Workday, Taleo, Greenhouse, and iCIMS platforms.
---

# CV Optimizer — ATS Resume Optimizer

This skill turns the agent into a full ATS scoring and CV optimization pipeline. When invoked, follow the steps below in order.

## Step 1 — Launch the TUI Dashboard

Immediately execute `cv-tui` (or `./cv-tui` if in local directory) in the terminal to launch the visual TUI dashboard:

```bash
cv-tui
```

The TUI will:
- Load the Job Description & CV
- Let the user pick the target ATS platform (Workday, Taleo, Greenhouse, iCIMS)
- Score the CV and show the total score, grade (A/B/C/D/F), and missing keywords

## Step 2 — Read the CV and JD

Read both files:
- CV: `cv.tex` (or `cv_test.txt` / any `.tex` / `.md` file the user points to)
- JD: `jd.txt` (or `jd_test.txt` / any `.txt` file the user points to)

If either file is missing, ask the user to paste the content directly into the chat.

## Step 3 — Analyze CV Keywords & ATS Fit

Compare the CV text against the Job Description:
1. Extract required hard skills, job title, and tool keywords from the JD.
2. Cross-check against the CV to compile matched vs missing keywords.
3. Apply platform scoring rules:
   - **Workday**: Exact keyword matches required.
   - **Taleo**: Keyword repetition & density.
   - **Greenhouse**: Natural phrasing with measurable achievements.
   - **iCIMS**: Direct job title alignment.

## Step 4 — Optimize the CV (Headless AI Script)

Run the optimization script to rewrite the CV:

```bash
# Using Antigravity (agy)
./optimize.sh cv.tex jd.txt Workday agy

# Using Claude Code
./optimize.sh cv.tex jd.txt Workday claude

# Using OpenCode
./optimize.sh cv.tex jd.txt Workday opencode
```

The optimized LaTeX output will be saved to `optimized_cv.md`.

## Step 5 — Compile to PDF

Compile `optimized_cv.md` into a PDF via `tectonic`:

```bash
tectonic optimized_cv.md -o output/Optimized_Resume.pdf
```

---

## LaTeX Formatting Rules (CRITICAL — never skip these)

When generating or editing LaTeX:
- Use EXACT structural commands from the original CV. Do NOT invent commands.
- `\resumeSubheading{Title}{Dates}{Company}{Location}` — exactly 4 args, no `\begin{}`.
- ALL `\resumeSubheading` must be inside `\resumeSubHeadingListStart` / `\resumeSubHeadingListEnd`.
- ALL `\resumeItem` must be inside `\resumeItemListStart` / `\resumeItemListEnd`.
- Escape `%` and `$` in plain text. Do NOT escape `&` inside tabular environments.
- Header: single line separated by `|` — `email | github | location`.
- Output ONLY `\begin{document}...\end{document}`, wrapped in a ` ```latex ` code block.
- Do NOT output the preamble.
