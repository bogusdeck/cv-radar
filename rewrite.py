import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

# Replace the broken prompt with backticks
prompt_broken = text[text.find('prompt := fmt.Sprintf('):text.find('req.Platform, req.CVText, req.JDText)')]

fixed_prompt = '''prompt := fmt.Sprintf(`You are an expert ATS (Applicant Tracking System) optimizer.
Your goal is to rewrite the provided CV so that it perfectly matches the provided Job Description (JD) and scores 100%%%% on the %s ATS platform.

Guidelines:
- Incorporate keywords seamlessly.
- Highlight relevant experience.
- DO NOT hallucinate fake jobs or degrees.
- Keep it professional and ATS-friendly.
- CRITICAL: Do NOT use tables (tabular), columns, or alignment environments. ATS parsers hate them.
- CRITICAL: You must explicitly escape all special characters, especially ampersands (\\&). Unescaped & characters will crash the compiler!
- CRITICAL: Use only standard ASCII characters. Do not use unicode characters like approx, use ~ instead.

OUTPUT FORMAT:
You MUST output ONLY a valid, complete, and compilable LaTeX document.
Use a standard class like \\documentclass[11pt,a4paper]{article}.
Include \\usepackage{geometry}, \\usepackage{hyperref}, etc.
Do NOT wrap the output in markdown code blocks. Output raw LaTeX text starting with \\documentclass and ending with \\end{document}.

--- CV TEXT ---
%s

--- JOB DESCRIPTION ---
%s`, '''

text = text.replace(prompt_broken, fixed_prompt)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

