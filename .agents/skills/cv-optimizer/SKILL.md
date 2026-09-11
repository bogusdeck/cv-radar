---
name: cv-optimizer
description: >-
  Use this skill when the user says "cv-optimizer", "/cv-optimizer", "optimize my CV", or "scan my resume".
  Launches an interactive ATS analysis and optimization workflow. Scores the CV against the job description,
  identifies missing keywords, then optionally rewrites the CV using the headless AI backend.
  Supports Workday, Taleo, Greenhouse, and iCIMS platforms.
---

# CV Optimizer — ATS Resume Optimizer

This skill turns the agent into a full ATS scoring and CV optimization pipeline. When invoked, follow the steps below in order.

## Step 1 — Launch the TUI Dashboard

Tell the user to run the TUI in a separate terminal tab to get the visual score:

```
./cv-tui
```

The TUI will:
- Load the Job Description from `jd_test.txt`
- Let the user pick the target ATS platform (Workday, Taleo, Greenhouse, iCIMS)
- Score the CV and show the total score, grade (A/B/C/D/F), and all missing keywords

## Step 2 — Read the CV and JD

Read both files now:
- CV: `cv_test.txt` (or any `.tex` / `.md` file the user points to)
- JD: `jd_test.txt` (or any `.txt` file the user points to)

If either file is missing, ask the user to paste the content directly into the chat.

## Step 3 — Run the ATS Scan (Go backend)

Use the Go binary to score the CV against the JD:

```bash
go run ./cmd/server &
# or if compiled already:
./server &
```

Then hit the API directly:

```bash
curl -X POST http://localhost:8085/api/analyze \
  -F "jd=@jd_test.txt" \
  -F "platform=Workday"
```

Print the result clearly: score, grade, matched keywords, missing keywords, recommendations.

## Step 4 — Identify Missing Keywords

From the scan result, extract the `MissingKeywords` list. These are the keywords in the JD that the CV does NOT contain. Tell the user exactly which ones are missing.

## Step 5 — Optimize the CV (Headless AI)

If the user asks to optimize, run the bash script that dispatches to the headless AI CLI:

```bash
# Using Claude Code
./optimize.sh cv_test.txt jd_test.txt Workday claude

# Using Antigravity (agy)
./optimize.sh cv_test.txt jd_test.txt Workday agy

# Using OpenCode
./optimize.sh cv_test.txt jd_test.txt Workday opencode
```

Wait for it to complete. The output will be saved to `optimized_cv.md`.

## Step 6 — Compile the Optimized CV to PDF

Once `optimized_cv.md` is ready, compile it to PDF via the Go backend:

```bash
curl -X POST http://localhost:8085/api/fix-resume/compile \
  -H "Content-Type: application/json" \
  -d "{\"latex\": \"$(cat optimized_cv.md | tr -d '\n' | sed 's/"/\\"/g')\"}"
```

The PDF will be returned as a download or saved to `output/`.

## Step 7 — Re-scan and Verify

Re-run Step 3 against the newly optimized CV to confirm the score improved. The target is 90+.

---

## ATS Platform Strategies

Apply these platform-specific strategies when generating the prompt in Step 5:

| Platform | Strategy |
|---|---|
| **Workday** | Strict exact-keyword matching. Use the EXACT terminology from the JD. Even slight rephrasing causes misses. |
| **Taleo** | High keyword density. Repeat the most critical JD keywords 2-3 times across Summary, Experience, and Skills. |
| **Greenhouse** | Semantic AI + human review. Avoid stuffing. Focus on measurable impact and natural phrasing. |
| **iCIMS** | Mixed exact + semantic. Align job titles closely to the JD title. |

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
