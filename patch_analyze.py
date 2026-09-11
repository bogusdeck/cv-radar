import re

with open("cmd/server/main.go", "r") as f:
    text = f.read()

target = r"latexDocClassRegex := regexp.MustCompile(`(?m)^\s*\\documentclass`)"
new_code = r"latexDocClassRegex := regexp.MustCompile(`(?m)^\s*\\(documentclass|begin\{document\})`)"

text = text.replace(target, new_code)

with open("cmd/server/main.go", "w") as f:
    f.write(text)

