import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

# 1. Change default platform
text = text.replace("useState(PLATFORMS[1].value)", "useState(PLATFORMS[0].value)")

# 2. Add PixelSelect component
pixel_select_code = """
function PixelSelect({ options, value, onChange }: { options: {name: string, value: string, emoji?: string}[], value: string, onChange: (val: string) => void }) {
  const [isOpen, setIsOpen] = useState(false);
  const selectedOption = options.find(o => o.value === value) || options[0];
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className="relative" ref={dropdownRef}>
      <div 
        className="bg-[#1a1a1a] text-white font-PressStart text-[10px] border-2 border-[#333] p-4 cursor-pointer hover:border-yellow-400 flex justify-between items-center min-w-[220px]"
        onClick={() => setIsOpen(!isOpen)}
      >
        <span>{selectedOption.emoji ? selectedOption.emoji + ' ' : ''}{selectedOption.name}</span>
        <span className="ml-4 text-yellow-400">{isOpen ? '▲' : '▼'}</span>
      </div>
      
      {isOpen && (
        <div className="absolute top-full left-0 w-full mt-1 bg-[#1a1a1a] border-2 border-[#333] z-50 max-h-[300px] overflow-y-auto shadow-[4px_4px_0_rgba(0,0,0,1)]">
          {options.map(p => (
            <div 
              key={p.value} 
              className={`p-3 font-PressStart text-[10px] cursor-pointer hover:bg-[#333] hover:text-yellow-400 ${value === p.value ? 'bg-[#222] text-yellow-400' : 'text-white'}`}
              onClick={() => {
                onChange(p.value);
                setIsOpen(false);
              }}
            >
              {p.emoji ? p.emoji + ' ' : ''}{p.name}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
"""

# Insert PixelSelect before export default function App
text = text.replace("export default function App() {", pixel_select_code + "\nexport default function App() {")

# 3. Replace <select> with <PixelSelect>
# Let's find the select block and replace it
select_block = r'<select\s*className="bg-\[#1a1a1a\] text-white font-PressStart text-\[10px\] border-2 border-\[#333\] p-4 focus:outline-none focus:border-yellow-400 cursor-pointer"\s*value=\{platform\}\s*onChange=\{e => setPlatform\(e.target.value\)\}\s*>\s*\{PLATFORMS\.map\(p => <option key=\{p\.value\} value=\{p\.value\}>\{p\.name\}</option>\)\}\s*</select>'

text = re.sub(select_block, "<PixelSelect options={PLATFORMS} value={platform} onChange={setPlatform} />", text, flags=re.DOTALL)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

