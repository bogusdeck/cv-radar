import { useState, useRef } from 'react'
import './index.css'

const API = 'http://localhost:8080'

const PLATFORMS = [
  { name: 'All Platforms', value: '' },
  { name: 'Workday', value: 'Workday', emoji: '🔵' },
  { name: 'Taleo (Oracle)', value: 'Taleo', emoji: '🟠' },
  { name: 'iCIMS', value: 'iCIMS', emoji: '🟢' },
  { name: 'Greenhouse', value: 'Greenhouse', emoji: '🌿' },
  { name: 'Lever', value: 'Lever', emoji: '🟣' },
  { name: 'SuccessFactors (SAP)', value: 'SuccessFactors', emoji: '🔴' },
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
  keyword_score: number
  structure_score: number
  experience_score: number
  education_score: number
  matched_keywords: string[]
  missing_keywords: string[]
  warnings: string[]
  recommendations: string[]
  auto_reject: boolean
  breakdown: { category: string; score: number; max_score: number; weight: number }[]
}

interface ParsedCV {
  name: string
  email: string
  skills: string[]
  years_exp: number
  experience: { title: string; company: string }[]
}

interface AnalysisResponse {
  parsed_cv: ParsedCV
  results: ATSResult[]
}

function ScoreCircle({ score, grade, low, high }: { score: number; grade: string; low: number; high: number }) {
  return (
    <div className={`score-circle grade-${grade}`}>
      <span className="score-number">{Math.round(score)}</span>
      <span className="score-label">{grade}</span>
      <span style={{ fontSize: '0.55rem', opacity: 0.6, marginTop: '2px', whiteSpace: 'nowrap' }}>
        {Math.round(low)}–{Math.round(high)}
      </span>
    </div>
  )
}

function ConfidencePill({ confidence }: { confidence: number }) {
  const color = confidence >= 70 ? 'var(--green)' : confidence >= 55 ? 'var(--yellow)' : 'var(--orange)'
  const label = confidence >= 70 ? 'High confidence' : confidence >= 55 ? 'Med confidence' : 'Low confidence'
  return (
    <span style={{
      fontSize: '0.6rem', fontWeight: 700, padding: '2px 7px', borderRadius: '99px',
      border: `1px solid ${color}22`, background: `${color}11`, color,
      whiteSpace: 'nowrap'
    }}>
      ~{Math.round(confidence)}% · {label}
    </span>
  )
}

function ResultCard({ result, selected, onClick }: { result: ATSResult; selected: boolean; onClick: () => void }) {
  const topMatched = (result.matched_keywords || []).slice(0, 5)
  const topMissing = (result.missing_keywords || []).slice(0, 4)

  return (
    <div
      className={`ats-result-card ${selected ? 'selected' : ''} ${result.auto_reject ? 'auto-reject' : ''}`}
      onClick={onClick}
    >
      <div className="result-header">
        <div className="platform-info">
          <div className="platform-name">{PLATFORM_EMOJIS[result.platform]} {result.platform}</div>
          <div className="vendor-name" style={{ display: 'flex', gap: '0.5rem', alignItems: 'center', marginTop: '4px' }}>
            by {result.vendor}
            <ConfidencePill confidence={result.confidence} />
          </div>
          {result.auto_reject && (
            <div className="auto-reject-badge">⛔ Auto-Reject Risk</div>
          )}
        </div>
        <ScoreCircle score={result.score} grade={result.grade} low={result.score_low} high={result.score_high} />
      </div>

      <div className="progress-bar">
        <div className="progress-fill" style={{ width: `${result.score}%` }} />
      </div>

      <div className="keyword-preview">
        {topMatched.map(kw => (
          <span key={kw} className="kw-tag kw-matched">✓ {kw}</span>
        ))}
        {topMissing.map(kw => (
          <span key={kw} className="kw-tag kw-missing">✗ {kw}</span>
        ))}
      </div>
    </div>
  )
}

function DetailPanel({ result, onClose }: { result: ATSResult; onClose: () => void }) {
  return (
    <div className="detail-panel">
      <div className="detail-header">
        <div>
          <div className="detail-title">{PLATFORM_EMOJIS[result.platform]} {result.platform} Deep Analysis</div>
          <div style={{ color: 'var(--text-muted)', fontSize: '0.8rem', marginTop: '0.25rem', display: 'flex', flexWrap: 'wrap', gap: '0.5rem', alignItems: 'center' }}>
            <span>by {result.vendor}</span>
            <span>·</span>
            <strong style={{ color: 'var(--text-primary)' }}>{result.score}/100</strong>
            <span>·</span>
            <strong style={{ color: `var(--grade-${result.grade})` }}>Grade {result.grade}</strong>
            <span>·</span>
            <ConfidencePill confidence={result.confidence} />
          </div>
          {result.simulation_note && (
            <div style={{
              marginTop: '0.5rem', padding: '0.5rem 0.75rem',
              background: 'rgba(234,179,8,0.06)', border: '1px solid rgba(234,179,8,0.2)',
              borderRadius: '8px', fontSize: '0.75rem', color: '#fde68a', lineHeight: 1.5
            }}>
              ⚠️ <strong>Simulation accuracy:</strong> {result.simulation_note}
            </div>
          )}
        </div>
        <button
          onClick={onClose}
          style={{
            background: 'var(--bg-surface)', border: '1px solid var(--border)',
            color: 'var(--text-secondary)', borderRadius: '8px',
            padding: '0.5rem 1rem', cursor: 'pointer', fontSize: '0.8rem', fontWeight: 600
          }}
        >✕ Close</button>
      </div>

      {/* Breakdown */}
      <div className="detail-section">
        <div className="detail-section-title">Score Breakdown</div>
        <div className="breakdown-grid">
          {(result.breakdown || []).map(b => (
            <div key={b.category} className="breakdown-item">
              <div className="breakdown-category">{b.category}</div>
              <div className="breakdown-score" style={{ color: scoreColor(b.score) }}>
                {Math.round(b.score)}<span style={{ fontSize: '0.8rem', opacity: 0.6 }}>/100</span>
              </div>
              <div className="breakdown-weight">Weight: {Math.round(b.weight * 100)}%</div>
            </div>
          ))}
        </div>
      </div>

      {/* Warnings */}
      {result.warnings && result.warnings.length > 0 && (
        <div className="detail-section">
          <div className="detail-section-title">⚠️ Warnings</div>
          <div className="recommendations-list">
            {result.warnings.map((w, i) => (
              <div key={i} className="warning-item">
                <span className="rec-icon">⚠️</span>
                {w}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Recommendations */}
      {result.recommendations && result.recommendations.length > 0 && (
        <div className="detail-section">
          <div className="detail-section-title">💡 Recommendations</div>
          <div className="recommendations-list">
            {result.recommendations.map((r, i) => (
              <div key={i} className="rec-item">
                <span className="rec-icon">→</span>
                {r}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Keywords */}
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
        {result.matched_keywords && result.matched_keywords.length > 0 && (
          <div className="detail-section">
            <div className="detail-section-title">✅ Matched Keywords ({result.matched_keywords.length})</div>
            <div className="keywords-list">
              {result.matched_keywords.map(kw => (
                <span key={kw} className="kw-tag kw-matched">{kw}</span>
              ))}
            </div>
          </div>
        )}
        {result.missing_keywords && result.missing_keywords.length > 0 && (
          <div className="detail-section">
            <div className="detail-section-title">❌ Missing Keywords ({result.missing_keywords.length})</div>
            <div className="keywords-list">
              {result.missing_keywords.map(kw => (
                <span key={kw} className="kw-tag kw-missing">{kw}</span>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

function scoreColor(score: number) {
  if (score >= 80) return 'var(--grade-a)'
  if (score >= 65) return 'var(--grade-b)'
  if (score >= 50) return 'var(--grade-c)'
  if (score >= 35) return 'var(--grade-d)'
  return 'var(--grade-f)'
}

export default function App() {
  const [cvText, setCvText] = useState('')
  const [jdText, setJdText] = useState('')
  const [platform, setPlatform] = useState('')
  const [loading, setLoading] = useState(false)
  const [response, setResponse] = useState<AnalysisResponse | null>(null)
  const [selectedResult, setSelectedResult] = useState<ATSResult | null>(null)
  const [error, setError] = useState('')
  const [uploading, setUploading] = useState(false)
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
      setError('Failed to upload file — is the server running?')
    } finally {
      setUploading(false)
      e.target.value = ''
    }
  }

  const handleAnalyze = async () => {
    if (!cvText.trim() || !jdText.trim()) {
      setError('Please paste both your CV and the Job Description.')
      return
    }
    setError('')
    setLoading(true)
    setSelectedResult(null)
    try {
      const res = await fetch(`${API}/api/analyze`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ cv_text: cvText, jd_text: jdText, platform }),
      })
      if (!res.ok) throw new Error(await res.text())
      const data: AnalysisResponse = await res.json()
      setResponse(data)
    } catch (e: unknown) {
      setError(`Analysis failed: ${e instanceof Error ? e.message : 'Server error'}. Make sure the Go server is running on port 8080.`)
    } finally {
      setLoading(false)
    }
  }

  const avgScore = response
    ? Math.round(response.results.reduce((s, r) => s + r.score, 0) / response.results.length)
    : null

  const bestPlatform = response
    ? response.results.reduce((a, b) => a.score > b.score ? a : b)
    : null

  const worstPlatform = response
    ? response.results.reduce((a, b) => a.score < b.score ? a : b)
    : null

  return (
    <div className="app-wrapper">
      <header className="header">
        <div className="header-logo">
          <div className="logo-icon">🎯</div>
          ATS Scanner
        </div>
        <span className="header-badge">Beta</span>
        <div className="header-spacer" />
        <span style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>
          6 ATS Platforms · Go Backend
        </span>
      </header>

      <main className="main-content">
        {/* Hero */}
        <div className="hero">
          <h1>Beat Every ATS.<br />Get the Interview.</h1>
          <p>Analyze your CV against 6 real ATS platforms — Workday, Taleo, iCIMS, Greenhouse, Lever, and SuccessFactors — with platform-specific scoring algorithms.</p>
          <div className="ats-chips">
            {PLATFORMS.slice(1).map(p => (
              <div key={p.value} className="ats-chip">{p.emoji} {p.name}</div>
            ))}
          </div>
        </div>

        {/* Input Panel */}
        <div className="card input-panel">
          <div className="card-title">Input</div>
          <div className="input-grid">
            <div className="textarea-wrapper">
              <div className="textarea-label">
                📄 CV / Resume
                <input
                  ref={fileRef}
                  type="file"
                  accept=".pdf,.tex,.txt"
                  style={{ display: 'none' }}
                  onChange={handleUpload}
                />
                <button
                  className="file-upload-btn"
                  onClick={() => fileRef.current?.click()}
                  disabled={uploading}
                >
                  {uploading ? '⏳ Uploading…' : '⬆ Upload PDF / .tex'}
                </button>
              </div>
              <textarea
                placeholder="Paste your CV text here, or upload a PDF / .tex file above…"
                value={cvText}
                onChange={e => setCvText(e.target.value)}
                rows={10}
              />
            </div>
            <div className="textarea-wrapper">
              <div className="textarea-label">📋 Job Description</div>
              <textarea
                placeholder="Paste the full job description here…"
                value={jdText}
                onChange={e => setJdText(e.target.value)}
                rows={10}
              />
            </div>
          </div>

          <div className="action-row">
            <select
              className="platform-select"
              value={platform}
              onChange={e => setPlatform(e.target.value)}
            >
              {PLATFORMS.map(p => (
                <option key={p.value} value={p.value}>{p.name}</option>
              ))}
            </select>
            <button
              className="analyze-btn"
              onClick={handleAnalyze}
              disabled={loading}
            >
              {loading ? <><div className="spinner" /> Analyzing…</> : '⚡ Analyze CV'}
            </button>
          </div>

          {error && (
            <div style={{
              marginTop: '1rem', padding: '0.75rem 1rem',
              background: 'rgba(239,68,68,0.08)', border: '1px solid rgba(239,68,68,0.25)',
              borderRadius: 'var(--radius-sm)', fontSize: '0.82rem', color: '#fca5a5'
            }}>{error}</div>
          )}
        </div>

        {/* Summary Stats */}
        {response && avgScore !== null && (
          <div className="summary-bar">
            <div className="summary-stat">
              <div className="summary-stat-label">Avg Score</div>
              <div className="summary-stat-value" style={{ color: scoreColor(avgScore) }}>{avgScore}</div>
            </div>
            <div className="summary-stat">
              <div className="summary-stat-label">Best Platform</div>
              <div className="summary-stat-value" style={{ fontSize: '1rem', marginTop: '0.25rem' }}>
                {PLATFORM_EMOJIS[bestPlatform!.platform]} {bestPlatform!.platform}
              </div>
            </div>
            <div className="summary-stat">
              <div className="summary-stat-label">Toughest Platform</div>
              <div className="summary-stat-value" style={{ fontSize: '1rem', marginTop: '0.25rem' }}>
                {PLATFORM_EMOJIS[worstPlatform!.platform]} {worstPlatform!.platform}
              </div>
            </div>
            <div className="summary-stat">
              <div className="summary-stat-label">Auto-Reject Risk</div>
              <div className="summary-stat-value" style={{ color: response.results.some(r => r.auto_reject) ? 'var(--red)' : 'var(--green)' }}>
                {response.results.filter(r => r.auto_reject).length > 0
                  ? `${response.results.filter(r => r.auto_reject).length} platform${response.results.filter(r => r.auto_reject).length > 1 ? 's' : ''}`
                  : 'None ✓'}
              </div>
            </div>
            <div className="summary-stat">
              <div className="summary-stat-label">Detected Name</div>
              <div className="summary-stat-value" style={{ fontSize: '0.9rem', marginTop: '0.25rem' }}>
                {response.parsed_cv.name || '—'}
              </div>
            </div>
            <div className="summary-stat">
              <div className="summary-stat-label">Years Experience</div>
              <div className="summary-stat-value">{response.parsed_cv.years_exp || '—'}</div>
            </div>
          </div>
        )}

        {/* Results grid */}
        {response && (
          <div className="results-panel">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem', flexWrap: 'wrap', gap: '0.5rem' }}>
              <div className="card-title" style={{ marginBottom: 0 }}>Platform Results</div>
              <div style={{
                padding: '0.4rem 0.8rem', background: 'rgba(234,179,8,0.06)',
                border: '1px solid rgba(234,179,8,0.25)', borderRadius: '8px',
                fontSize: '0.72rem', color: '#fde68a', display: 'flex', alignItems: 'center', gap: '0.4rem'
              }}>
                ⚠️ <strong>Simulation</strong> — scores are estimates, not real ATS outputs. Score ranges shown on each card. Click a card for accuracy details.
              </div>
            </div>
            <div className="results-grid">
              {response.results.map(r => (
                <ResultCard
                  key={r.platform}
                  result={r}
                  selected={selectedResult?.platform === r.platform}
                  onClick={() => setSelectedResult(selectedResult?.platform === r.platform ? null : r)}
                />
              ))}
            </div>
          </div>
        )}

        {/* Detail panel */}
        {selectedResult && (
          <DetailPanel result={selectedResult} onClose={() => setSelectedResult(null)} />
        )}

        {/* Parsed skills */}
        {response && response.parsed_cv.skills.length > 0 && (
          <div className="parsed-cv">
            <div className="card-title">Detected Skills from your CV ({response.parsed_cv.skills.length})</div>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.4rem', marginTop: '0.75rem' }}>
              {response.parsed_cv.skills.map(s => (
                <span key={s} className="skill-tag">{s}</span>
              ))}
            </div>
          </div>
        )}

        {/* Empty state */}
        {!response && !loading && (
          <div className="empty-state">
            <div className="empty-icon">🎯</div>
            <p>Paste your CV and a job description, then click <strong>Analyze CV</strong> to get platform-specific ATS scores.</p>
          </div>
        )}
      </main>

      <footer className="footer">
        ATS Scanner · Built with Go + React · 6 platforms · Open Source
      </footer>
    </div>
  )
}
