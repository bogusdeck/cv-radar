import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

# 1. Add cvText and jdText to DetailPanel props
text = text.replace(
    "function DetailPanel({ result, onClose }: { result: ATSResult; onClose: () => void }) {",
    "function DetailPanel({ result, onClose, cvText, jdText }: { result: ATSResult; onClose: () => void; cvText: string; jdText: string }) {"
)

# 2. Add fix resume state and button inside DetailPanel
# We need to find the X button line and add the FIX RESUME button next to it (if grade == F or score < 80)
# But wait, we can just put the button in the header.
new_header = """          <div>
            <h2 className="font-PressStart text-xl text-yellow-400">{PLATFORM_EMOJIS[result.platform]} {result.platform} Analysis</h2>
            <div className="font-PressStart text-[10px] text-gray-400 mt-2">Score: {Math.round(result.score)} | Grade: {result.grade}</div>
          </div>
          <div className="flex gap-4 items-start">
            {result.score < 80 && (
              <button 
                onClick={async () => {
                  setFixing(true);
                  try {
                    const res = await fetch(`${API}/api/fix-resume`, {
                      method: 'POST',
                      headers: { 'Content-Type': 'application/json' },
                      body: JSON.stringify({ cv_text: cvText, jd_text: jdText, platform: result.platform })
                    });
                    if (!res.ok) throw new Error(await res.text());
                    const blob = await res.blob();
                    const url = URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.href = url;
                    a.download = `Optimized_Resume_${result.platform}.pdf`;
                    a.click();
                  } catch (e: any) {
                    alert('Fix failed: ' + e.message);
                  } finally {
                    setFixing(false);
                  }
                }}
                disabled={fixing}
                className="nes-btn is-warning text-[10px]"
              >
                {fixing ? 'FIXING...' : 'FIX RESUME'}
              </button>
            )}
            <button onClick={onClose} className="nes-btn is-error text-[10px]">X</button>
          </div>"""

# add const [fixing, setFixing] = useState(false); inside DetailPanel
text = text.replace(
    "function DetailPanel({ result, onClose, cvText, jdText }: { result: ATSResult; onClose: () => void; cvText: string; jdText: string }) {\n",
    "function DetailPanel({ result, onClose, cvText, jdText }: { result: ATSResult; onClose: () => void; cvText: string; jdText: string }) {\n  const [fixing, setFixing] = useState(false);\n"
)

# replace header
old_header_regex = r"          <div>\s*<h2 className=\"font-PressStart text-xl text-yellow-400\">\{PLATFORM_EMOJIS\[result\.platform\]\} \{result\.platform\} Analysis<\/h2>\s*<div className=\"font-PressStart text-\[10px\] text-gray-400 mt-2\">Score: \{Math\.round\(result\.score\)\} \| Grade: \{result\.grade\}<\/div>\s*<\/div>\s*<button onClick=\{onClose\} className=\"nes-btn is-error\">X<\/button>"

text = re.sub(old_header_regex, new_header, text)

# 3. Pass cvText and jdText from App
text = text.replace(
    "<DetailPanel result={selectedResult} onClose={() => setSelectedResult(null)} />",
    "<DetailPanel result={selectedResult} onClose={() => setSelectedResult(null)} cvText={cvText} jdText={jdText} />"
)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

