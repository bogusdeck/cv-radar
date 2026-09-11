with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

text = text.replace(r"\n", "\n")

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

