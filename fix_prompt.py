import re
with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()
    
start = text.find("prompt :=")
end = text.find("reqBody, _ :=")
if start != -1 and end != -1:
    new_prompt = """	prompt := fmt.Sprintf("You are an expert ATS (Applicant Tracking System) optimizer.\\nYour goal is to rewrite the provided CV so that it perfectly matches the provided Job Description (JD) and scores 100%%%% on the %s ATS platform.\\n\\nGuidelines:\\n- Incorporate keywords seamlessly.\\n- Highlight relevant experience.\\n- DO NOT hallucinate fake jobs or degrees, but adapt the phrasing and bullet points to match the JD's terminology exactly.\\n- Keep it professional and ATS-friendly (no weird columns or graphics).\\n\\nOUTPUT FORMAT:\\nYou MUST output ONLY a valid, complete, and compilable LaTeX document.\\nUse a standard class like \\\\documentclass[11pt,a4paper]{article}.\\nInclude \\\\usepackage{geometry}, \\\\usepackage{hyperref}, etc.\\nEnsure all special characters are escaped properly.\\nDo NOT wrap the output in markdown code blocks (```latex ... ```). Output raw LaTeX text starting with \\\\documentclass and ending with \\\\end{document}.\\n\\n--- CV TEXT ---\\n%s\\n\\n--- JOB DESCRIPTION ---\\n%s", req.Platform, req.CVText, req.JDText)

	"""
    text = text[:start] + new_prompt + text[end:]
    with open("cmd/server/fix_resume.go", "w") as f:
        f.write(text)
