with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

text = text.replace("esumeSubheading", r"\resumeSubheading")
text = text.replace("DO NOT use \\ or |", r"DO NOT use \\ or |")
text = text.replace("random \\ line breaks", r"random \\ line breaks")

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

