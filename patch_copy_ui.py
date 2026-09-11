import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

# 1. Add copied state
state_target = r"const \[pastedLatex, setPastedLatex\] = useState\(''\);"
state_new = r"const [pastedLatex, setPastedLatex] = useState('');\n  const [copied, setCopied] = useState(false);"
text = text.replace(state_target, state_new)

# 2. Update copyPrompt function
func_target = r"""  const copyPrompt = () => {
    navigator.clipboard.writeText(promptText);
    alert("Prompt copied to clipboard! Paste it into ChatGPT, Claude, or Gemini.");
  };"""
func_new = r"""  const copyPrompt = () => {
    navigator.clipboard.writeText(promptText);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };"""
text = text.replace(func_target, func_new)

# 3. Update button
btn_target = r"""              <button 
                onClick={copyPrompt}
                className="nes-btn is-primary w-full text-[10px]"
              >
                COPY PROMPT
              </button>"""
btn_new = r"""              <button 
                onClick={copyPrompt}
                className={`nes-btn w-full text-[10px] ${copied ? 'is-success' : 'is-primary'}`}
              >
                {copied ? 'PROMPT COPIED!' : 'COPY PROMPT'}
              </button>"""
text = text.replace(btn_target, btn_new)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

