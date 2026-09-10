import { useState, useRef } from 'react'
import DitherBackground from './DitherBackground'
import InputDitherBackground from './InputDitherBackground'
import './index.css'

const API = 'http://localhost:8080'

const PLATFORMS = [
  { name: 'All Platforms', value: '' },
  { name: 'Workday', value: 'Workday', emoji: '🔵' },
  { name: 'Taleo', value: 'Taleo', emoji: '🟠' },
  { name: 'iCIMS', value: 'iCIMS', emoji: '🟢' },
  { name: 'Greenhouse', value: 'Greenhouse', emoji: '🌿' },
  { name: 'Lever', value: 'Lever', emoji: '🟣' },
  { name: 'SuccessFactors', value: 'SuccessFactors', emoji: '🔴' },
]

const PLATFORM_EMOJIS: Record<string, string> = {
  Workday: '🔵', Taleo: '🟠', iCIMS: '🟢',
  Greenhouse: '🌿', Lever: '🟣', SuccessFactors: '🔴',
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
          <div className="font-PressStart text-[14px] text-white">{PLATFORM_EMOJIS[result.platform]} {result.platform}</div>
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

function DetailPanel({ result, onClose }: { result: ATSResult; onClose: () => void }) {
  return (
    <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4 overflow-y-auto">
      <div className="boxxy bg-black p-6 w-full max-w-4xl border-4 border-white max-h-[90vh] overflow-y-auto">
        <div className="flex justify-between items-start border-b-4 border-white pb-4 mb-4">
          <div>
            <h2 className="font-PressStart text-xl text-yellow-400">{PLATFORM_EMOJIS[result.platform]} {result.platform} Analysis</h2>
            <div className="font-PressStart text-[10px] text-gray-400 mt-2">Score: {Math.round(result.score)} | Grade: {result.grade}</div>
          </div>
          <button onClick={onClose} className="nes-btn is-error">X</button>
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

export default function App() {
  const [cvText, setCvText] = useState('')
  const [jdText, setJdText] = useState('')
  const [platform, setPlatform] = useState(PLATFORMS[1].value)
  const [uploading, setUploading] = useState(false)
  const [loading, setLoading] = useState(false)
  const [response, setResponse] = useState<AnalysisResponse | null>(null)
  const [error, setError] = useState('')
  const [selectedResult, setSelectedResult] = useState<ATSResult | null>(null)
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
                {p.emoji} {p.name}
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Input Section */}
      <div className="relative h-fit w-full flex flex-col items-center justify-center shrink-0 overflow-hidden py-2">
        <div className="absolute inset-0 z-[-1] pointer-events-none" style={{ background: 'linear-gradient(90deg, rgba(7, 11, 26, 0.96) 0%, rgba(7, 11, 26, 0.82) 38%, rgba(7, 11, 26, 0.55) 70%, rgba(7, 11, 26, 0.75) 100%)' }} />
        <div className="absolute inset-0 z-[-1] pointer-events-none" style={{ background: 'radial-gradient(ellipse at center, transparent 20%, rgba(0, 0, 0, 0.25) 65%, rgba(0, 0, 0, 0.7) 100%)' }} />

        <div className="max-w-[95%] 2xl:max-w-[1400px] mx-auto w-full z-10">
          <div className="terminal-window h-[85vh] min-h-[600px] max-h-[1000px] w-full shadow-[12px_12px_0_rgba(0,0,0,1)] flex flex-col relative">
            <div className="terminal-header shrink-0 relative z-20">
              <span className="terminal-dot" style={{ background: '#ff5f56' }} />
              <span className="terminal-dot" style={{ background: '#ffbd2e' }} />
              <span className="terminal-dot" style={{ background: '#27c93f' }} />
              <span className="ml-4 font-PressStart text-[10px] text-gray-300">cv_radar.exe</span>
            </div>
            
            <div className="flex-1 relative overflow-hidden flex flex-col">
              <InputDitherBackground />
              
              <div className="relative z-10 p-4 md:p-8 flex-1 flex flex-col gap-4 md:gap-6 h-full">
                <h2 className="text-center font-Gumball text-2xl md:text-4xl text-yellow-400 tracking-widest shrink-0">SCAN YOUR CV</h2>
                
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
                <select 
                  className="bg-[#1a1a1a] text-white font-PressStart text-[10px] border-2 border-[#333] p-4 focus:outline-none focus:border-yellow-400 cursor-pointer" 
                  value={platform} 
                  onChange={e => setPlatform(e.target.value)}
                >
                  {PLATFORMS.map(p => <option key={p.value} value={p.value}>{p.name}</option>)}
                </select>
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

      {selectedResult && <DetailPanel result={selectedResult} onClose={() => setSelectedResult(null)} />}
    </div>
  )
}
