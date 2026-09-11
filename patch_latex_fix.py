import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

# 1. Update the Prompt
old_prompt = r'prompt := fmt\.Sprintf\("You are an expert ATS.*?req\.Platform, req\.CVText, req\.JDText\)'
new_prompt = """prompt := fmt.Sprintf("You are an expert ATS (Applicant Tracking System) optimizer.\\nYour goal is to rewrite the provided CV so that it perfectly matches the provided Job Description (JD) and scores 100%%%% on the %s ATS platform.\\n\\nGuidelines:\\n- Incorporate keywords seamlessly.\\n- Highlight relevant experience.\\n- DO NOT hallucinate fake jobs or degrees.\\n- Keep it professional and ATS-friendly.\\n- CRITICAL: Do NOT use tables (`tabular`), columns, or alignment environments. ATS parsers hate them.\\n- CRITICAL: You must explicitly escape all special characters, especially ampersands (\\\\&). Unescaped & characters will crash the compiler!\\n- CRITICAL: Use only standard ASCII characters. Do not use unicode characters like ≈, use ~ or approx instead.\\n\\nOUTPUT FORMAT:\\nYou MUST output ONLY a valid, complete, and compilable LaTeX document.\\nUse a standard class like \\\\documentclass[11pt,a4paper]{article}.\\nInclude \\\\usepackage{geometry}, \\\\usepackage{hyperref}, etc.\\nDo NOT wrap the output in markdown code blocks. Output raw LaTeX text starting with \\\\documentclass and ending with \\\\end{document}.\\n\\n--- CV TEXT ---\\n%s\\n\\n--- JOB DESCRIPTION ---\\n%s", req.Platform, req.CVText, req.JDText)"""

text = re.sub(old_prompt, new_prompt, text, flags=re.DOTALL)

# 2. Update extractLatex to include regex sanitization
old_extract = """func extractLatex(output string) string {
	start := strings.Index(output, "\\\\documentclass")
	if start == -1 {
		return output
	}
	endStr := "\\\\end{document}"
	end := strings.LastIndex(output, endStr)
	if end == -1 {
		return output[start:]
	}
	return output[start : end+len(endStr)]
}"""

new_extract = """import "regexp"

func extractLatex(output string) string {
	start := strings.Index(output, "\\\\documentclass")
	if start != -1 {
		endStr := "\\\\end{document}"
		end := strings.LastIndex(output, endStr)
		if end != -1 {
			output = output[start : end+len(endStr)]
		} else {
			output = output[start:]
		}
	}
	
	// Safety net: escape unescaped ampersands and percents which crash tectonic
	reAmp := regexp.MustCompile(`([^\\\\])&`)
	output = reAmp.ReplaceAllString(output, `$1\\&`)
	
	// Remove common problematic unicode characters
	output = strings.ReplaceAll(output, "≈", "~")
	output = strings.ReplaceAll(output, "–", "--") // en-dash
	output = strings.ReplaceAll(output, "—", "---") // em-dash
	
	return output
}"""

text = text.replace(old_extract, new_extract)
# Because we injected `import "regexp"`, we need to make sure we don't have duplicated import blocks that cause syntax errors.
# Actually, it's safer to add regexp to the top imports.

text = text.replace('import "regexp"\n\n', '')
text = text.replace('"path/filepath"', '"path/filepath"\n\t"regexp"')

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

