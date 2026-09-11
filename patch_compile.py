import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

target = r"""	texPath := filepath.Join(tmpDir, "resume.tex")
	if err := os.WriteFile(texPath, []byte(req.Latex), 0644); err != nil {"""

new_code = r"""	// Always run extractLatex on incoming raw LaTeX to ensure preamble is attached and sanitized
	finalLatex := extractLatex(req.Latex)
	
	texPath := filepath.Join(tmpDir, "resume.tex")
	if err := os.WriteFile(texPath, []byte(finalLatex), 0644); err != nil {"""

text = text.replace(target, new_code)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

