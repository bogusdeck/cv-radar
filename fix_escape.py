with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

text = text.replace('"\\\\documentclass"', '"\\\\documentclass"')
text = text.replace('"\\\\end{document}"', '"\\\\end{document}"')

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)
