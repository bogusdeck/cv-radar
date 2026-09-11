with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

text = text.replace('"bytes"\n', "")
text = text.replace('"encoding/json"\n', "")
text = text.replace('"io"\n', "")

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)
