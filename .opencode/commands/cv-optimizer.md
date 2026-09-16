---
description: ATS Resume Scanner & AI CV Optimizer — launch interactive TUI dashboard and auto-rewrite CV
---

You are an expert ATS (Applicant Tracking System) CV analyzer and resume optimizer.

When this command is triggered:
1. **Launch Interactive TUI**: Run `osascript -e 'tell application "Terminal" to do script "cv-tui"' 2>/dev/null || cv-tui open || cv-tui` in the terminal to immediately open the visual ATS score dashboard in an interactive terminal window for the user.
2. **Analyze Workspace Documents**: Read `cv.tex` (or `cv_test.txt` / any `.tex` or `.md` CV file) and `jd.txt` (or `jd_test.txt` / any `.txt` JD file) in the workspace.
3. **Perform Keyword & ATS Audit**: Extract required skills, detect missing keywords, and calculate the target ATS score & grade (Workday, Taleo, Greenhouse, iCIMS).
4. **Headless Optimization**: Execute `./optimize.sh <path_to_cv> <path_to_jd> <platform> opencode` to auto-rewrite the resume into `optimized_cv.md`.
5. **PDF Compilation**: Compile `optimized_cv.md` via `tectonic optimized_cv.md -o output/Optimized_Resume.pdf`.
