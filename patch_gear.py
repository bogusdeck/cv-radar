import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

# 1. Add PixelGear component
gear_comp = """
const PixelGear = () => (
  <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg" className="inline-block mr-2 -mt-1">
    <path d="M6 0h4v2h2v2h2v2h2v4h-2v2h-2v2h-2v2H6v-2H4v-2H2v-2H0V6h2V4h2V2h2V0Zm4 6H6v4h4V6Z" fillRule="evenodd" clipRule="evenodd" />
  </svg>
)
"""

text = text.replace("export default function App() {", gear_comp + "\nexport default function App() {")

# 2. Remove the AI SETUP button from the terminal header
header_btn_regex = r"<button onClick=\{.*?\} className=\"absolute right-4 top-1/2 -translate-y-1/2 text-\[10px\] font-PressStart text-gray-400 hover:text-yellow-400 transition-colors\">⚙️ AI SETUP</button>"
text = re.sub(header_btn_regex, "", text)

# 3. Add it inside the SCAN YOUR CV header
# Old: <h2 className="text-center font-Gumball text-2xl md:text-4xl text-yellow-400 tracking-widest shrink-0">SCAN YOUR CV</h2>
# New: wrap in a relative div, put button on right
scan_h2 = r"<h2 className=\"text-center font-Gumball text-2xl md:text-4xl text-yellow-400 tracking-widest shrink-0\">SCAN YOUR CV</h2>"
new_scan_h2 = """<div className="relative flex justify-center items-center shrink-0 w-full">
                  <h2 className="text-center font-Gumball text-2xl md:text-4xl text-yellow-400 tracking-widest">SCAN YOUR CV</h2>
                  <button onClick={() => setShowSettings(true)} className="absolute right-0 text-[10px] font-PressStart text-gray-400 hover:text-yellow-400 transition-colors flex items-center">
                    <PixelGear /> SETUP
                  </button>
                </div>"""
text = text.replace('<h2 className="text-center font-Gumball text-2xl md:text-4xl text-yellow-400 tracking-widest shrink-0">SCAN YOUR CV</h2>', new_scan_h2)

# Also fix the settings modal title to use PixelGear
text = text.replace("⚙️ AI AGENT SETUP", "<PixelGear /> AI AGENT SETUP")

with open("web/src/App.tsx", "w") as f:
    f.write(text)

