import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

# 1. Update PLATFORMS and PLATFORM_EMOJIS
text = text.replace("emoji: '🔵'", "color: '#3b82f6'")
text = text.replace("emoji: '🟠'", "color: '#f97316'")
text = text.replace("emoji: '🟢'", "color: '#22c55e'")
text = text.replace("emoji: '🌿'", "color: '#84cc16'")
text = text.replace("emoji: '🟣'", "color: '#a855f7'")
text = text.replace("emoji: '🔴'", "color: '#ef4444'")

text = text.replace(
    "const PLATFORM_EMOJIS: Record<string, string> = {\n  Workday: '🔵', Taleo: '🟠', iCIMS: '🟢',\n  Greenhouse: '🌿', Lever: '🟣', SuccessFactors: '🔴',\n}",
    "const PLATFORM_COLORS: Record<string, string> = {\n  Workday: '#3b82f6', Taleo: '#f97316', iCIMS: '#22c55e',\n  Greenhouse: '#84cc16', Lever: '#a855f7', SuccessFactors: '#ef4444',\n}"
)

# 2. Update PixelSelect component
text = text.replace(
    "options: {name: string, value: string, emoji?: string}[]",
    "options: {name: string, value: string, color?: string}[]"
)
text = text.replace(
    "<span>{selectedOption.emoji ? selectedOption.emoji + ' ' : ''}{selectedOption.name}</span>",
    "<span className=\"flex items-center gap-2\">{selectedOption.color && <span className=\"inline-block w-2 h-2\" style={{ backgroundColor: selectedOption.color }}></span>}{selectedOption.name}</span>"
)
text = text.replace(
    "{p.emoji ? p.emoji + ' ' : ''}{p.name}",
    "<span className=\"flex items-center gap-2\">{p.color && <span className=\"inline-block w-2 h-2\" style={{ backgroundColor: p.color }}></span>}{p.name}</span>"
)

# 3. Update ResultCard
text = text.replace(
    "<div className=\"font-PressStart text-[14px] text-white\">{PLATFORM_EMOJIS[result.platform]} {result.platform}</div>",
    "<div className=\"font-PressStart text-[14px] text-white flex items-center gap-3\"><span className=\"inline-block w-3 h-3\" style={{ backgroundColor: PLATFORM_COLORS[result.platform] }}></span>{result.platform}</div>"
)

# 4. Update DetailPanel
text = text.replace(
    "<h2 className=\"font-PressStart text-xl text-yellow-400\">{PLATFORM_EMOJIS[result.platform]} {result.platform} Analysis</h2>",
    "<h2 className=\"font-PressStart text-xl text-yellow-400 flex items-center gap-3\"><span className=\"inline-block w-4 h-4\" style={{ backgroundColor: PLATFORM_COLORS[result.platform] }}></span>{result.platform} Analysis</h2>"
)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

