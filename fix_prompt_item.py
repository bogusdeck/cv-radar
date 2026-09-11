import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

target = r"- CRITICAL: \resumeSubheading is a COMMAND, not an environment!"
new_rule = r"- CRITICAL: ALL \resumeSubheading and \resumeProjectHeading commands MUST be wrapped inside \resumeSubHeadingListStart and \resumeSubHeadingListEnd!\n- CRITICAL: \resumeSubheading is a COMMAND, not an environment!"

text = text.replace(target, new_rule)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

