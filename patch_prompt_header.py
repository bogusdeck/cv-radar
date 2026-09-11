import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

target = r"- MATCH ORIGINAL LENGTH: If the uploaded CV is a single page, keep the output single page. Adjust verbosity accordingly."
new_rule = r"""- MATCH ORIGINAL LENGTH: You MUST strictly ensure the final output fits on the exact same number of pages as the original CV (e.g. if the original is 1 page, the output MUST be 1 page). Be concise, consolidate bullet points, and do NOT add fluff.
- CRITICAL: Format the personal details (email, links, location, etc.) in the header on a SINGLE LINE separated by `|` (e.g. email@val.com | github.com/val | Delhi, India) to save vertical space."""

text = text.replace(target, new_rule)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

