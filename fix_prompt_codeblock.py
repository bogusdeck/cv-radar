import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

# Change the prompt instruction
target_prompt = r"Do NOT wrap the output in markdown code blocks."
new_prompt = r"CRITICAL: Wrap the entire output in a single markdown ```latex code block."
text = text.replace(target_prompt, new_prompt)

# Add stripping to handleCompile
target_compile = r"body: JSON.stringify({ latex: pastedLatex })"
new_compile = r"""body: JSON.stringify({ 
          latex: pastedLatex.replace(/^```(latex)?\n?/i, '').replace/\n?```$/g, '') 
        })"""
# Wait, let's use a simpler clean logic in JS:
new_compile_safe = r"""body: JSON.stringify({ 
          latex: pastedLatex.replace(/```latex/gi, '').replace(/```/g, '').trim()
        })"""

text = text.replace(target_compile, new_compile_safe)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

