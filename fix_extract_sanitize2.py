import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

target = r"""	output = strings.ReplaceAll(output, "\\resumeProjectEnd", "\\resumeSubHeadingListEnd")"""

new_code = r"""	output = strings.ReplaceAll(output, "\\resumeProjectEnd", "\\resumeSubHeadingListEnd")
	output = regexp.MustCompile(`\\resumeSubheading\{\s*\}\{\s*\}\{\s*\}(\{\s*\})?`).ReplaceAllString(output, `\\item`)"""

text = text.replace(target, new_code)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

