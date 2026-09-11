import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

text = text.replace("separated by `|`", "separated by '|'")

with open("web/src/App.tsx", "w") as f:
    f.write(text)

