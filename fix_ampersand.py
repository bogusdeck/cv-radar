import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

target = r"""	// Escape unescaped percent signs (which act as comments and eat closing braces)
	output = regexp.MustCompile(`([^\\])%`).ReplaceAllString(output, `$1\%`)
	// Also escape unescaped ampersands safely outside tabular (if any sneak through)
	// Actually, we skip ampersands because they are used in tabular. We only escape %"""

new_code = r"""	// Escape unescaped percent signs (which act as comments and eat closing braces)
	output = regexp.MustCompile(`([^\\])%`).ReplaceAllString(output, `$1\%`)
	// Escape ampersands because the LLM shouldn't be generating tabular alignment tabs anyway.
	// (The structural & are safe inside the Preamble constant)
	output = regexp.MustCompile(`([^\\])&`).ReplaceAllString(output, `$1\&`)"""

text = text.replace(target, new_code)

target_unicode = r"""	output = strings.ReplaceAll(output, "≈", "~")"""
new_unicode = r"""	output = strings.ReplaceAll(output, "≈", "~")
	output = strings.ReplaceAll(output, "∼", "~") // U+223C"""

text = text.replace(target_unicode, new_unicode)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

