import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

# 1. Update the Prompt to include the template and new rules
old_prompt_regex = r"prompt := fmt\.Sprintf\(`You are an expert ATS.*?req\.Platform, req\.CVText, req\.JDText\)"

new_prompt = """prompt := fmt.Sprintf(`You are an expert ATS (Applicant Tracking System) optimizer.
Your goal is to rewrite the provided CV so that it perfectly matches the provided Job Description (JD) and scores 100%%%% on the %s ATS platform.

Guidelines:
- Incorporate keywords seamlessly into the Summary, Experience, Projects, and Skills sections.
- Highlight relevant experience. DO NOT hallucinate fake jobs, degrees, or metrics.
- MATCH ORIGINAL LENGTH: If the uploaded CV is a single page, keep the output single page. If double, keep double. Adjust verbosity accordingly.
- CRITICAL: You must use the EXACT LaTeX template/layout provided in the CV TEXT below. Do NOT alter the preamble, custom commands, or overall structure. Only modify the actual textual content (bullet points, summary, skills) to optimize for the JD.
- CRITICAL: Ensure all plain-text special characters like %% and $ are escaped. Do NOT escape the alignment ampersands (&) inside the tabular* environments used by the custom commands!
- CRITICAL: Use only standard ASCII characters. Do not use unicode characters like approx, use ~ instead.

OUTPUT FORMAT:
You MUST output ONLY a valid, complete, and compilable LaTeX document based on the provided template.
Do NOT wrap the output in markdown code blocks. Output raw LaTeX text starting with \\documentclass and ending with \\end{document}.

--- CV TEXT (USE THIS EXACT LATEX TEMPLATE) ---
%s

--- JOB DESCRIPTION ---
%s`, req.Platform, req.CVText, req.JDText)"""

text = re.sub(old_prompt_regex, new_prompt, text, flags=re.DOTALL)

# 2. Remove the dumb ampersand regex safety net in extractLatex
re_amp = r"""	// Safety net: escape unescaped ampersands and percents which crash tectonic
	reAmp := regexp\.MustCompile\(`\(\[\^\\\\\\\\\]\)&`\)
	output = reAmp\.ReplaceAllString\(output, `\$1\\\\&`\)
"""
text = re.sub(re_amp, "", text)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

