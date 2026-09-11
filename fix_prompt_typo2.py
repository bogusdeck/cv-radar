import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

# Fix the specific line
target_line = re.search(r'- CRITICAL: When using.*', text).group(0)
new_line = r"- CRITICAL: When using \resumeSubheading, you MUST provide exactly 4 arguments in curly braces: {Title}{Dates}{Company}{Location}. DO NOT use \\ or | to combine them!"
text = text.replace(target_line, new_line)

# Also fix the previous newline replacement issue just in case
text = text.replace('\r\n', '\n').replace('\r', '')

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

