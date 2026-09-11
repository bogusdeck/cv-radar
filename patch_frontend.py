import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

# 1. Add Settings Modal Component
settings_modal = """
function SettingsModal({ config, setConfig, onClose }: any) {
  return (
    <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-[100] p-4">
      <div className="boxxy bg-black p-6 w-full max-w-lg border-4 border-white shadow-[8px_8px_0_rgba(0,0,0,1)]">
        <div className="flex justify-between items-center mb-6 border-b-2 border-[#333] pb-4">
          <h2 className="font-PressStart text-sm text-yellow-400">⚙️ AI AGENT SETUP</h2>
          <button onClick={onClose} className="nes-btn is-error text-[10px]">X</button>
        </div>
        
        <div className="flex flex-col gap-4 font-PressStart text-[10px] text-white">
          <label>Provider</label>
          <select 
            value={config.provider} 
            onChange={e => setConfig({...config, provider: e.target.value})}
            className="bg-[#1a1a1a] text-white border-2 border-[#333] p-3 focus:outline-none focus:border-yellow-400"
          >
            <option value="ollama">Local/Cloud Ollama</option>
            <option value="openai">OpenAI (ChatGPT)</option>
            <option value="anthropic">Anthropic (Claude)</option>
            <option value="gemini">Google Gemini</option>
          </select>

          <label>Model Name (Optional)</label>
          <input 
            type="text" 
            placeholder={config.provider === 'ollama' ? 'llama3' : 'gpt-4o / claude-3-5-sonnet'}
            value={config.model}
            onChange={e => setConfig({...config, model: e.target.value})}
            className="bg-[#1a1a1a] text-white border-2 border-[#333] p-3 focus:outline-none focus:border-yellow-400"
          />

          {config.provider !== 'ollama' && (
            <>
              <label>API Key</label>
              <input 
                type="password" 
                value={config.apiKey}
                onChange={e => setConfig({...config, apiKey: e.target.value})}
                className="bg-[#1a1a1a] text-white border-2 border-[#333] p-3 focus:outline-none focus:border-yellow-400"
              />
            </>
          )}

          {config.provider === 'ollama' && (
            <>
              <label>Ollama Base URL</label>
              <input 
                type="text" 
                placeholder="http://localhost:11434"
                value={config.baseUrl}
                onChange={e => setConfig({...config, baseUrl: e.target.value})}
                className="bg-[#1a1a1a] text-white border-2 border-[#333] p-3 focus:outline-none focus:border-yellow-400"
              />
            </>
          )}

          <button onClick={onClose} className="nes-btn is-success mt-4">SAVE CONFIG</button>
        </div>
      </div>
    </div>
  )
}
"""

text = text.replace("export default function App() {", settings_modal + "\nexport default function App() {")

# 2. Add state to App
state_str = "  const [showSettings, setShowSettings] = useState(false)\n  const [llmConfig, setLlmConfig] = useState({ provider: 'ollama', model: '', apiKey: '', baseUrl: '' })"
text = text.replace("  const [selectedResult, setSelectedResult] = useState<ATSResult | null>(null)", "  const [selectedResult, setSelectedResult] = useState<ATSResult | null>(null)\n" + state_str)

# 3. Add Settings button to terminal header
header_btn = """              <span className="ml-4 font-PressStart text-[10px] text-gray-300">cv_radar.exe</span>
              <button onClick={() => setShowSettings(true)} className="absolute right-4 top-1/2 -translate-y-1/2 text-[10px] font-PressStart text-gray-400 hover:text-yellow-400 transition-colors">⚙️ AI SETUP</button>"""
text = text.replace("              <span className=\"ml-4 font-PressStart text-[10px] text-gray-300\">cv_radar.exe</span>", header_btn)

# 4. Pass llmConfig to DetailPanel
text = text.replace(
    "function DetailPanel({ result, onClose, cvText, jdText }: { result: ATSResult; onClose: () => void; cvText: string; jdText: string }) {",
    "function DetailPanel({ result, onClose, cvText, jdText, llmConfig }: { result: ATSResult; onClose: () => void; cvText: string; jdText: string; llmConfig: any }) {"
)
text = text.replace(
    "<DetailPanel result={selectedResult} onClose={() => setSelectedResult(null)} cvText={cvText} jdText={jdText} />",
    "<DetailPanel result={selectedResult} onClose={() => setSelectedResult(null)} cvText={cvText} jdText={jdText} llmConfig={llmConfig} />"
)

# 5. Include llmConfig in fetch POST body
post_body = "body: JSON.stringify({ cv_text: cvText, jd_text: jdText, platform: result.platform, provider: llmConfig.provider, model: llmConfig.model, api_key: llmConfig.apiKey, base_url: llmConfig.baseUrl })"
text = text.replace(
    "body: JSON.stringify({ cv_text: cvText, jd_text: jdText, platform: result.platform })",
    post_body
)

# 6. Render SettingsModal
text = text.replace(
    "{selectedResult && <DetailPanel",
    "{showSettings && <SettingsModal config={llmConfig} setConfig={setLlmConfig} onClose={() => setShowSettings(false)} />}\n      {selectedResult && <DetailPanel"
)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

