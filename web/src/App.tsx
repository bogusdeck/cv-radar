import { useState, useRef, useEffect } from 'react'
import DitherBackground from './DitherBackground'
import InputDitherBackground from './InputDitherBackground'
import './index.css'

const API = 'http://localhost:8080'

const PLATFORMS = [
  { name: 'All Platforms', value: '' },
  { name: 'Workday', value: 'Workday', color: '#3b82f6' },
  { name: 'Taleo', value: 'Taleo', color: '#f97316' },
  { name: 'iCIMS', value: 'iCIMS', color: '#22c55e' },
  { name: 'Greenhouse', value: 'Greenhouse', color: '#84cc16' },
  { name: 'Lever', value: 'Lever', color: '#a855f7' },
  { name: 'SuccessFactors', value: 'SuccessFactors', color: '#ef4444' },
]

const PLATFORM_COLORS: Record<string, string> = {
  Workday: '#3b82f6', Taleo: '#f97316', iCIMS: '#22c55e',
  Greenhouse: '#84cc16', Lever: '#a855f7', SuccessFactors: '#ef4444',
}

interface ATSResult {
  platform: string
  vendor: string
  score: number
  score_low: number
  score_high: number
  grade: string
  confidence: number
  simulation_note: string
  matched_keywords: string[]
  missing_keywords: string[]
  warnings: string[]
  recommendations: string[]
  auto_reject: boolean
  breakdown: { category: string; score: number; max_score: number; weight: number }[]
}

interface AnalysisResponse {
  results: ATSResult[]
}

function ScoreCircle({ score, grade }: { score: number; grade: string }) {
  const color = score >= 80 ? '#22c55e' : score >= 65 ? '#84cc16' : score >= 50 ? '#eab308' : score >= 35 ? '#f97316' : '#ef4444'
  return (
    <div className="flex flex-col items-center justify-center p-2 boxxy" style={{ backgroundColor: color }}>
      <span className="font-PressStart text-black text-2xl">{Math.round(score)}</span>
      <span className="font-PressStart text-black text-[8px] mt-1">GRADE {grade}</span>
    </div>
  )
}

function ResultCard({ result, onClick }: { result: ATSResult; onClick: () => void }) {
  const topMatched = (result.matched_keywords || []).slice(0, 3)
  const topMissing = (result.missing_keywords || []).slice(0, 3)

  return (
    <div className="boxxy p-4 cursor-pointer hover:bg-gray-900 transition-colors" onClick={onClick}>
      <div className="flex justify-between items-start mb-4">
        <div>
          <div className="font-PressStart text-[14px] text-white flex items-center gap-3"><span className="inline-block w-3 h-3" style={{ backgroundColor: PLATFORM_COLORS[result.platform] }}></span>{result.platform}</div>
          <div className="font-PressStart text-[8px] text-gray-400 mt-2">by {result.vendor}</div>
          {result.auto_reject && <div className="font-PressStart text-[8px] text-red-500 mt-2 animate-pulse">⛔ AUTO-REJECT</div>}
        </div>
        <ScoreCircle score={result.score} grade={result.grade} />
      </div>

      <div className="w-full h-4 boxxy mb-4" style={{ backgroundColor: '#111' }}>
        <div className="h-full" style={{ width: `${result.score}%`, backgroundColor: result.score >= 80 ? '#22c55e' : '#eab308' }} />
      </div>

      <div className="flex flex-col gap-2">
        <div className="font-PressStart text-[8px] text-green-400">
          {topMatched.length > 0 ? `✓ ${topMatched.join(', ')}` : 'No keywords matched'}
        </div>
        <div className="font-PressStart text-[8px] text-red-400">
          {topMissing.length > 0 ? `✗ ${topMissing.join(', ')}` : 'Perfect!'}
        </div>
      </div>
    </div>
  )
}

function DetailPanel({ result, onClose, cvText, jdText, llmConfig }: { result: ATSResult; onClose: () => void; cvText: string; jdText: string; llmConfig: any }) {
  const [fixing, setFixing] = useState(false);
  const [fixMessage, setFixMessage] = useState('FIX RESUME');


  return (
    <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4 overflow-y-auto">
      <div className="boxxy bg-black p-6 w-full max-w-4xl border-4 border-white max-h-[90vh] overflow-y-auto">
        <div className="flex justify-between items-start border-b-4 border-white pb-4 mb-4">
          <div>
            <h2 className="font-PressStart text-xl text-yellow-400 flex items-center gap-3"><span className="inline-block w-4 h-4" style={{ backgroundColor: PLATFORM_COLORS[result.platform] }}></span>{result.platform} Analysis</h2>
            <div className="font-PressStart text-[10px] text-gray-400 mt-2">Score: {Math.round(result.score)} | Grade: {result.grade}</div>
          </div>
          <div className="flex gap-8 items-start">
            {result.score < 80 && (
              <button 
                onClick={async () => {
                  setFixing(true);
                  try {
                    setFixMessage('BOOTING AGENT...');
                    await new Promise(r => setTimeout(r, 600));

                    setFixMessage('REWRITING CV...');
                    const generateRes = await fetch(`${API}/api/fix-resume/generate`, {
                      method: 'POST',
                      headers: { 'Content-Type': 'application/json' },
                      body: JSON.stringify({ cv_text: cvText, jd_text: jdText, platform: result.platform, provider: llmConfig.provider, model: llmConfig.model, api_key: llmConfig.apiKey, base_url: llmConfig.baseUrl })
                    });
                    if (!generateRes.ok) throw new Error(await generateRes.text());
                    const { latex } = await generateRes.json();

                    setFixMessage('COMPILING PDF...');
                    const compileRes = await fetch(`${API}/api/fix-resume/compile`, {
                      method: 'POST',
                      headers: { 'Content-Type': 'application/json' },
                      body: JSON.stringify({ latex })
                    });
                    if (!compileRes.ok) throw new Error(await compileRes.text());
                    
                    const blob = await compileRes.blob();
                    try {
                      if ('showSaveFilePicker' in window) {
                        const handle = await (window as any).showSaveFilePicker({
                          suggestedName: `Optimized_Resume_${result.platform}.pdf`,
                          types: [{
                            description: 'PDF Document',
                            accept: {'application/pdf': ['.pdf']},
                          }],
                        });
                        const writable = await handle.createWritable();
                        await writable.write(blob);
                        await writable.close();
                      } else {
                        throw new Error('Fallback');
                      }
                    } catch (err: any) {
                      if (err.name !== 'AbortError') {
                        const url = URL.createObjectURL(blob);
                        const a = document.createElement('a');
                        a.href = url;
                        a.download = `Optimized_Resume_${result.platform}.pdf`;
                        a.click();
                      }
                    }

                    setFixMessage('SUCCESS!');
                    await new Promise(r => setTimeout(r, 1000));
                  } catch (e: any) {
                    alert('Fix failed: ' + e.message);
                  } finally {
                    setFixing(false);
                    setFixMessage('FIX RESUME');
                  }
                }}
                disabled={fixing}
                className="nes-btn is-warning text-[10px]"
              >
                {fixMessage}
              </button>
            )}
            <button onClick={onClose} className="nes-btn is-error text-[10px]">X</button>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="pixel-border p-4">
            <h3 className="font-PressStart text-[12px] text-green-400 mb-4">Matched ({result.matched_keywords?.length || 0})</h3>
            <div className="flex flex-wrap gap-2">
              {result.matched_keywords?.map(kw => <span key={kw} className="font-PressStart text-[8px] bg-green-900 border border-green-400 px-2 py-1">{kw}</span>)}
            </div>
          </div>
          
          <div className="pixel-border p-4">
            <h3 className="font-PressStart text-[12px] text-red-400 mb-4">Missing ({result.missing_keywords?.length || 0})</h3>
            <div className="flex flex-wrap gap-2">
              {result.missing_keywords?.map(kw => <span key={kw} className="font-PressStart text-[8px] bg-red-900 border border-red-400 px-2 py-1">{kw}</span>)}
            </div>
          </div>
        </div>
        
        {result.recommendations && result.recommendations.length > 0 && (
           <div className="mt-6 pixel-border p-4 border-yellow-400">
             <h3 className="font-PressStart text-[12px] text-yellow-400 mb-4">Recommendations</h3>
             <ul className="flex flex-col gap-4">
               {result.recommendations.map((r, i) => (
                 <li key={i} className="font-PressStart text-[8px] leading-relaxed break-words whitespace-pre-wrap">{r}</li>
               ))}
             </ul>
           </div>
        )}
      </div>
    </div>
  )
}


function PixelSelect({ options, value, onChange }: { options: {name: string, value: string, color?: string}[], value: string, onChange: (val: string) => void }) {
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
    <div className="relative z-50" ref={dropdownRef}>
      <div 
        className="bg-[#1a1a1a] text-white font-PressStart text-[10px] border-2 border-[#333] p-4 cursor-pointer hover:border-yellow-400 flex justify-between items-center min-w-[220px]"
        onClick={() => setIsOpen(!isOpen)}
      >
        <span className="flex items-center gap-2">{selectedOption.color && <span className="inline-block w-2 h-2" style={{ backgroundColor: selectedOption.color }}></span>}{selectedOption.name}</span>
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
              <span className="flex items-center gap-2">{p.color && <span className="inline-block w-2 h-2" style={{ backgroundColor: p.color }}></span>}{p.name}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}


function SettingsModal({ config, setConfig, onClose }: any) {
  return (
    <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-[100] p-4">
      <div className="boxxy bg-black p-6 w-full max-w-lg border-4 border-white shadow-[8px_8px_0_rgba(0,0,0,1)]">
        <div className="flex justify-between items-center mb-6 border-b-2 border-[#333] pb-4">
          <h2 className="font-PressStart text-sm text-yellow-400"><PixelGear /> AI AGENT SETUP</h2>
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
            placeholder={config.provider === 'ollama' ? 'qwen2.5:7b-instruct-q4_K_M' : 'gpt-4o / claude-3-5-sonnet'}
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


const PixelGear = () => (
  <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg" className="inline-block mr-2 -mt-1">
    <path d="M6 0h4v2h2v2h2v2h2v4h-2v2h-2v2h-2v2H6v-2H4v-2H2v-2H0V6h2V4h2V2h2V0Zm4 6H6v4h4V6Z" fillRule="evenodd" clipRule="evenodd" />
  </svg>
)

export default function App() {
  const [cvText, setCvText] = useState('')
  const [jdText, setJdText] = useState('')
  const [platform, setPlatform] = useState(PLATFORMS[0].value)
  const [uploading, setUploading] = useState(false)
  const [loading, setLoading] = useState(false)
  const [response, setResponse] = useState<AnalysisResponse | null>(null)
  const [error, setError] = useState('')
  const [selectedResult, setSelectedResult] = useState<ATSResult | null>(null)
  const [showSettings, setShowSettings] = useState(false)
  const [llmConfig, setLlmConfig] = useState({ provider: 'ollama', model: 'qwen2.5:7b-instruct-q4_K_M', apiKey: '', baseUrl: '' })
  const fileRef = useRef<HTMLInputElement>(null)

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    setUploading(true)
    const form = new FormData()
    form.append('file', file)
    try {
      const res = await fetch(`${API}/api/upload`, { method: 'POST', body: form })
      const data = await res.json()
      if (data.text) setCvText(data.text)
    } catch {
      setError('Failed to upload file.')
    } finally {
      setUploading(false)
      e.target.value = ''
    }
  }

  const handleAnalyze = async () => {
    if (!cvText.trim() || !jdText.trim()) {
      setError('Paste CV and JD!')
      return
    }
    setError('')
    setLoading(true)
    setSelectedResult(null)

    try {
      const res = await fetch(`${API}/api/analyze`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ cv_text: cvText, jd_text: jdText, target_platform: platform })
      })
      const data = await res.json()
      if (data.error) throw new Error(data.error)
      setResponse(data)
    } catch (err: any) {
      setError(err.message || 'Analysis failed.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="font-VT323 min-h-screen bg-[#0a0a0a] text-white">
      {/* Hero Section */}
      <div className="relative mb-12 min-h-screen flex flex-col items-center justify-center px-4 overflow-hidden">
        
        <DitherBackground />
        <div className="scanlines absolute inset-0 w-full h-full z-10 mix-blend-overlay opacity-80 pointer-events-none"></div>
        <div className="absolute inset-0 w-full h-full z-10 bg-gradient-to-b from-[#0a0a0a]/40 via-transparent to-[#0a0a0a]/90 pointer-events-none"></div>
        
        {/* Top Left Logo */}
        <div className="absolute top-4 left-4 md:top-8 md:left-8 z-20">
          <img src="/assets/images/logo.png" alt="Travel Coffer Logo" className="w-12 h-12 md:w-16 md:h-16 drop-shadow-[4px_4px_0_#000]" style={{ imageRendering: 'pixelated' }} />
        </div>

        <div className="relative z-10 w-full max-w-4xl mx-auto flex flex-col items-center">
          <h1 className="text-5xl md:text-8xl text-center text-yellow-400 font-Gumball drop-shadow-[6px_6px_0_#000] mb-8">CV-RADAR</h1>
          <p className="text-center font-PressStart text-[12px] md:text-[16px] text-white drop-shadow-[4px_4px_0_#000]">Beat Every ATS.</p>
          
          <div className="flex flex-wrap justify-center gap-8 mt-12 px-4">
            {PLATFORMS.slice(1).map(p => (
              <div key={p.value} className="boxxy bg-black/90 px-4 py-3 font-PressStart text-[10px] md:text-[12px] mx-2 my-2">
                <span className="flex items-center gap-2"><span className="inline-block w-2 h-2" style={{ backgroundColor: p.color }}></span>{p.name}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Input Section */}
      <div className="relative h-fit w-full flex flex-col items-center justify-center shrink-0 py-2">
        <div className="absolute inset-0 z-[-1] pointer-events-none" style={{ background: 'linear-gradient(90deg, rgba(7, 11, 26, 0.96) 0%, rgba(7, 11, 26, 0.82) 38%, rgba(7, 11, 26, 0.55) 70%, rgba(7, 11, 26, 0.75) 100%)' }} />
        <div className="absolute inset-0 z-[-1] pointer-events-none" style={{ background: 'radial-gradient(ellipse at center, transparent 20%, rgba(0, 0, 0, 0.25) 65%, rgba(0, 0, 0, 0.7) 100%)' }} />

        <div className="max-w-[95%] 2xl:max-w-[1400px] mx-auto w-full z-40">
          <div className="terminal-window h-[85vh] min-h-[600px] max-h-[1000px] w-full shadow-[12px_12px_0_rgba(0,0,0,1)] flex flex-col relative">
            <div className="terminal-header shrink-0 relative z-20">
              <span className="terminal-dot" style={{ background: '#ff5f56' }} />
              <span className="terminal-dot" style={{ background: '#ffbd2e' }} />
              <span className="terminal-dot" style={{ background: '#27c93f' }} />
              <span className="ml-4 font-PressStart text-[10px] text-gray-300">cv_radar.exe</span>
              
            </div>
            
            <div className="flex-1 relative flex flex-col">
              <InputDitherBackground />
              
              <div className="relative z-10 p-4 md:p-8 flex-1 flex flex-col gap-4 md:gap-6 h-full">
                <div className="relative flex justify-center items-center shrink-0 w-full">
                  <h2 className="text-center font-Gumball text-2xl md:text-4xl text-yellow-400 tracking-widest">SCAN YOUR CV</h2>
                  <button onClick={() => setShowSettings(true)} className="absolute right-0 text-[10px] font-PressStart text-gray-400 hover:text-yellow-400 transition-colors flex items-center">
                    <PixelGear /> SETUP
                  </button>
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-8 flex-1 min-h-0">
                  <div className="flex flex-col gap-2 flex-1 min-h-0">
                    <div className="flex justify-between items-end h-[24px] shrink-0">
                      <span className="font-PressStart text-[10px] md:text-[12px] text-yellow-400">CV / RESUME</span>
                      <button 
                        className="font-PressStart text-[8px] bg-white text-black px-2 py-1 hover:bg-yellow-400 hover:text-black transition-colors cursor-pointer border-b-2 border-gray-400 active:border-b-0 active:mt-[2px]" 
                        onClick={() => fileRef.current?.click()} 
                        disabled={uploading}
                      >
                        {uploading ? 'UPLOADING...' : '[ UPLOAD PDF ]'}
                      </button>
                      <input ref={fileRef} type="file" accept=".pdf,.txt" className="hidden" onChange={handleUpload} />
                    </div>
                    <textarea 
                      className="flex-1 w-full bg-[#1a1a1a]/90 text-green-400 font-VT323 text-lg md:text-xl p-4 border-2 border-[#333] focus:outline-none focus:border-yellow-400 resize-none shadow-inner"
                      placeholder="> Paste CV text here..."
                      value={cvText}
                      onChange={e => setCvText(e.target.value)}
                    />
                  </div>
                  
                  <div className="flex flex-col gap-2 flex-1 min-h-0">
                    <div className="flex justify-between items-end h-[24px] shrink-0">
                      <span className="font-PressStart text-[10px] md:text-[12px] text-yellow-400">JOB DESCRIPTION</span>
                    </div>
                    <textarea 
                      className="flex-1 w-full bg-[#1a1a1a]/90 text-blue-400 font-VT323 text-lg md:text-xl p-4 border-2 border-[#333] focus:outline-none focus:border-yellow-400 resize-none shadow-inner"
                      placeholder="> Paste JD text here..."
                    value={jdText}
                    onChange={e => setJdText(e.target.value)}
                  />
                </div>
              </div>

              <div className="flex flex-col md:flex-row justify-center items-center gap-6 mt-4 shrink-0 mb-4">
                <PixelSelect options={PLATFORMS} value={platform} onChange={setPlatform} />
                <button 
                  className="nes-btn text-xl px-10 py-4 bg-yellow-400 hover:bg-yellow-300 text-black border-4 border-white shadow-[6px_6px_0_#000] active:shadow-[2px_2px_0_#000] active:translate-y-1 transition-all" 
                  onClick={handleAnalyze} 
                  disabled={loading}
                >
                  {loading ? 'SCANNING...' : 'SCAN NOW'}
                </button>
              </div>
              
              {error && <div className="mt-4 font-PressStart text-red-500 text-[10px] text-center bg-red-900/30 p-4 border-2 border-red-500">{error}</div>}
            </div>
          </div>
        </div>
      </div>

      </div>
      <main className="max-w-[95%] 2xl:max-w-[1400px] mx-auto flex flex-col px-2 md:px-8 pb-12 pt-12">
        {/* Results Section */}
        {response && (
          <div className="mt-8">
            <h2 className="section-title text-center text-2xl drop-shadow-[4px_4px_0_#000]">ANALYSIS RESULTS</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
              {response.results.map(r => (
                <ResultCard key={r.platform} result={r} onClick={() => setSelectedResult(r)} />
              ))}
            </div>
          </div>
        )}
      </main>

      {showSettings && <SettingsModal config={llmConfig} setConfig={setLlmConfig} onClose={() => setShowSettings(false)} />}
      {selectedResult && <DetailPanel result={selectedResult} onClose={() => setSelectedResult(null)} cvText={cvText} jdText={jdText} llmConfig={llmConfig} />}
    </div>
  )
}
