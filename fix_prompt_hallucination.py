import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

# Add a new CRITICAL rule to the prompt
target = r"- CRITICAL: \resumeSubheading is a COMMAND, not an environment!"
new_rule = r"- CRITICAL: Do NOT invent or hallucinate new LaTeX commands like \resumeSummary or \resumeSkills! Use standard text for the summary and standard \section{} commands.\n- CRITICAL: \resumeSubheading is a COMMAND, not an environment!"

text = text.replace(target, new_rule)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

