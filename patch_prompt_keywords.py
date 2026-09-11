import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

target = r"- Incorporate keywords seamlessly into the Summary, Experience, Projects, and Skills sections."
new_rule = r"""- CRITICAL KEYWORD INJECTION: You MUST thoroughly analyze the JD for all required hard skills, tools, and keywords (e.g. Kubernetes, Databases, architecture, etc.) and FORCEFULLY inject them into the CV's Skills, Summary, and Experience sections. Do not leave ANY required JD keywords out!
- Incorporate these keywords seamlessly so they read naturally."""

text = text.replace(target, new_rule)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

