import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

target = r"- CRITICAL: ALL \resumeSubheading and \resumeProjectHeading commands MUST be wrapped inside \resumeSubHeadingListStart and \resumeSubHeadingListEnd!\n- CRITICAL: ALL \resumeItem commands MUST be wrapped inside \resumeItemListStart and \resumeItemListEnd!"

new_rule = r"""- CRITICAL FORMAT EXAMPLE:
\section{Experience}
  \resumeSubHeadingListStart
    \resumeSubheading{Title}{Dates}{Company}{Location}
      \resumeItemListStart
        \resumeItem{Bullet point 1}
      \resumeItemListEnd
  \resumeSubHeadingListEnd"""

text = text.replace(target, new_rule)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

