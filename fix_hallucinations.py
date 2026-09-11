import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

target = r"""	output = strings.ReplaceAll(output, "\\resumeProjectEnd", "\\resumeSubHeadingListEnd")"""

new_code = r"""	output = strings.ReplaceAll(output, "\\resumeProjectEnd", "\\resumeSubHeadingListEnd")
	output = strings.ReplaceAll(output, "\\resumeProjectItem", "\\resumeItem")
	output = strings.ReplaceAll(output, "\\resumeSkillItem", "\\resumeItem")
	output = strings.ReplaceAll(output, "\\resumeExperienceItem", "\\resumeItem")"""

text = text.replace(target, new_code)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

