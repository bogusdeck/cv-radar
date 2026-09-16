---
name: cv-optimizer
description: >-
  ATS Resume Scanner & AI CV Optimizer — scan CV, score against job description,
  launch interactive TUI dashboard, and auto-rewrite CVs for Workday, Taleo,
  Greenhouse, and iCIMS.
---

# CV Optimizer — ATS Resume Optimizer

This skill turns the agent into a full ATS scoring and CV optimization pipeline. When invoked via `/cv-optimizer` or "optimize my CV", follow the steps below in order.

## Step 1 — Launch the Interactive TUI Dashboard

Immediately launch the visual TUI dashboard in a focused Terminal window on screen:

```bash
osascript -e 'tell application "Terminal" to activate' -e 'tell application "Terminal" to do script "cv-tui"' 2>/dev/null || cv-tui open || cv-tui
```

The TUI will:
- Load the Job Description & CV
- Let the user pick the target ATS platform (Workday, Taleo, Greenhouse, iCIMS)
- Score the CV and show total score, grade (A/B/C/D/F), and missing keywords

---

## Step 2 — Read Workspace CV and JD

Read both files in workspace:
- CV: `cv.tex` (or `cv_test.txt` / any `.tex` or `.md` CV file)
- JD: `jd.txt` (or `jd_test.txt` / any `.txt` JD file)

If either file is missing, ask the user to paste the content directly into the chat.

---

## Step 3 — Analyze CV Keywords & ATS Fit

Compare CV text against Job Description:
1. Extract required hard skills, job title, and tool keywords from the JD.
2. Cross-check against the CV to compile matched vs missing keywords.
3. Apply platform scoring rules:
   - **Workday**: Exact keyword matches required.
   - **Taleo**: Keyword repetition & density.
   - **Greenhouse**: Natural phrasing with measurable achievements.
   - **iCIMS**: Direct job title alignment.

Display the score, grade, matched keywords, and missing keywords in clean Markdown tables.

---

## Step 4 — Optimize the CV (Headless AI Script)

Run the optimization script to rewrite the CV:

```bash
./optimize.sh cv.tex jd.txt Workday opencode
```

The optimized LaTeX output will be saved to `optimized_cv.md`.

---

## Step 5 — Compile to PDF

Compile `optimized_cv.md` into a PDF via `tectonic`:

```bash
tectonic optimized_cv.md -o output/Optimized_Resume.pdf
```

---

## LaTeX Formatting Rules (CRITICAL)

When generating or editing LaTeX:
- Use EXACT structural commands from the original CV. Do NOT invent commands.
- `\resumeSubheading{Title}{Dates}{Company}{Location}` — exactly 4 args, no `\begin{}`.
- ALL `\resumeSubheading` must be inside `\resumeSubHeadingListStart` / `\resumeSubHeadingListEnd`.
- ALL `\resumeItem` must be inside `\resumeItemListStart` / `\resumeItemListEnd`.
- Escape `%` and `$` in plain text. Do NOT escape `&` inside tabular environments.
