import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

target = r"""	// Remove common problematic unicode characters"""

new_code = r"""	// Escape unescaped percent signs (which act as comments and eat closing braces)
	output = regexp.MustCompile(`([^\\])%`).ReplaceAllString(output, `$1\%`)
	// Also escape unescaped ampersands safely outside tabular (if any sneak through)
	// Actually, we skip ampersands because they are used in tabular. We only escape %
	
	// Remove common problematic unicode characters"""

text = text.replace(target, new_code)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

