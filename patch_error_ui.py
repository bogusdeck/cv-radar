import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

# 1. Add errorMsg state
state_target = r"const \[fixing, setFixing\] = useState\(false\)"
state_new = "const [fixing, setFixing] = useState(false)\n  const [errorMsg, setErrorMsg] = useState('')"
text = re.sub(state_target, state_new, text)

# 2. Update the catch block in onClick
catch_target = r"\} catch \(e: any\) \{\s*alert\('Fix failed: ' \+ e\.message\);\s*\}"
catch_new = r"""} catch (e: any) {
                    console.error("Full Error:", e.message);
                    setErrorMsg('Failed to process or compile the CV. This is usually caused by an AI formatting glitch. Please try again.');
                  }"""
text = re.sub(catch_target, catch_new, text)

# 3. Add the custom error popup in the JSX before the closing </div>
popup_ui = r"""      {errorMsg && (
        <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50">
          <div className="bg-[#1a1a1a] border-4 border-red-500 p-8 max-w-md w-full text-center m-4">
            <h2 className="text-red-500 text-2xl font-bold mb-4 font-mono">SYSTEM ERROR</h2>
            <p className="text-gray-300 font-mono mb-6">{errorMsg}</p>
            <button 
              onClick={() => setErrorMsg('')}
              className="bg-red-500 text-black px-8 py-3 font-bold font-mono hover:bg-red-400 uppercase w-full tracking-wider"
            >
              Acknowledge
            </button>
          </div>
        </div>
      )}
    </div>"""

text = re.sub(r"    </div>\s*$", popup_ui + "\n", text)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

