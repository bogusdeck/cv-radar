import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

text = text.replace("markdown ```latex code block.", r"markdown \`\`\`latex code block.")

with open("web/src/App.tsx", "w") as f:
    f.write(text)

