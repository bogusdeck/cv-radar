with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

old_rule = "- CRITICAL: Use only standard ASCII characters. Do not use unicode characters like approx, use ~ instead."
new_rules = """- CRITICAL: Use only standard ASCII characters. Do not use unicode characters like approx, use ~ instead.
- CRITICAL: When using \resumeSubheading, you MUST provide exactly 4 arguments in curly braces: {Title}{Dates}{Company}{Location}. DO NOT use \\\\ or | to combine them!
- CRITICAL: Do NOT insert random \\\\ line breaks in normal text or between arguments. LaTeX handles wrapping automatically."""

text = text.replace(old_rule, new_rules)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)
