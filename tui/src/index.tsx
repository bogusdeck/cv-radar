#!/usr/bin/env node
import React, { useState, useEffect } from 'react';
import { render, Box, Text, useInput, useApp, Newline, Spacer } from 'ink';
import * as fs from 'fs';
import * as path from 'path';
import * as readline from 'readline';

const API = process.env.ATS_API || 'http://localhost:8080';

// ── Types ────────────────────────────────────────────────────────────────────

interface ATSResult {
  platform: string;
  vendor: string;
  score: number;
  grade: string;
  keyword_score: number;
  structure_score: number;
  experience_score: number;
  matched_keywords: string[];
  missing_keywords: string[];
  warnings: string[];
  recommendations: string[];
  auto_reject: boolean;
  breakdown: { category: string; score: number; max_score: number; weight: number }[];
}

interface AnalysisResponse {
  parsed_cv: { name: string; email: string; skills: string[]; years_exp: number };
  results: ATSResult[];
}

// ── Helpers ───────────────────────────────────────────────────────────────────

function gradeColor(grade: string): string {
  const map: Record<string, string> = { A: 'green', B: 'greenBright', C: 'yellow', D: 'red', F: 'redBright' };
  return map[grade] ?? 'white';
}

function scoreBar(score: number, width = 20): string {
  const filled = Math.round((score / 100) * width);
  return '█'.repeat(filled) + '░'.repeat(width - filled);
}

function scoreBarColor(score: number): string {
  if (score >= 80) return 'green';
  if (score >= 60) return 'yellow';
  if (score >= 40) return 'redBright';
  return 'red';
}

async function fetchText(url: string, body: object): Promise<AnalysisResponse> {
  const http = url.startsWith('https') ? require('https') : require('http');
  return new Promise((resolve, reject) => {
    const payload = JSON.stringify(body);
    const urlObj = new URL(url);
    const options = {
      hostname: urlObj.hostname,
      port: urlObj.port || 80,
      path: urlObj.pathname,
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Content-Length': Buffer.byteLength(payload) },
    };
    const req = http.request(options, (res: any) => {
      let data = '';
      res.on('data', (chunk: string) => (data += chunk));
      res.on('end', () => {
        try { resolve(JSON.parse(data)); }
        catch { reject(new Error('Invalid JSON from API')); }
      });
    });
    req.on('error', reject);
    req.write(payload);
    req.end();
  });
}

// ── Components ────────────────────────────────────────────────────────────────

function Header() {
  return (
    <Box flexDirection="column" marginBottom={1}>
      <Box>
        <Text bold color="magenta">╔══════════════════════════════════════════════╗</Text>
      </Box>
      <Box>
        <Text bold color="magenta">║  </Text>
        <Text bold color="white">🎯  ATS SCANNER</Text>
        <Text color="gray">  · 6 Platform CV Analyzer        </Text>
        <Text bold color="magenta">║</Text>
      </Box>
      <Box>
        <Text bold color="magenta">╚══════════════════════════════════════════════╝</Text>
      </Box>
    </Box>
  );
}

function ResultRow({ result }: { result: ATSResult }) {
  const bar = scoreBar(result.score);
  const barColor = scoreBarColor(result.score) as any;
  const gradeCol = gradeColor(result.grade) as any;

  return (
    <Box flexDirection="column" marginBottom={1}>
      <Box>
        <Text bold color="white">{result.platform.padEnd(16)}</Text>
        <Text color={gradeCol} bold> [{result.grade}] </Text>
        <Text color={barColor}>{bar}</Text>
        <Text color="white" bold> {String(Math.round(result.score)).padStart(3)}%</Text>
        {result.auto_reject && <Text color="red" bold>  ⛔ AUTO-REJECT RISK</Text>}
      </Box>
      <Box marginLeft={16}>
        <Text color="gray">by {result.vendor.padEnd(14)}</Text>
        <Text color="green">✓ {result.matched_keywords.length} matched  </Text>
        <Text color="red">✗ {result.missing_keywords.length} missing</Text>
      </Box>
    </Box>
  );
}

function DetailView({ result, onBack }: { result: ATSResult; onBack: () => void }) {
  useInput((input, key) => {
    if (input === 'b' || input === 'q' || key.escape) onBack();
  });

  return (
    <Box flexDirection="column">
      <Header />
      <Box marginBottom={1}>
        <Text bold color="cyan">━━━ {result.platform} ({result.vendor}) ━━━</Text>
      </Box>

      {/* Score breakdown */}
      <Box marginBottom={1} flexDirection="column">
        <Text bold color="gray">SCORE BREAKDOWN:</Text>
        {result.breakdown.map(b => (
          <Box key={b.category}>
            <Text color="gray">{b.category.padEnd(35)}</Text>
            <Text color={scoreBarColor(b.score) as any}>{scoreBar(b.score, 15)}</Text>
            <Text color="white" bold> {Math.round(b.score).toString().padStart(3)}%</Text>
            <Text color="gray"> (weight: {Math.round(b.weight * 100)}%)</Text>
          </Box>
        ))}
      </Box>

      {/* Warnings */}
      {result.warnings.length > 0 && (
        <Box flexDirection="column" marginBottom={1}>
          <Text bold color="red">⚠  WARNINGS:</Text>
          {result.warnings.map((w, i) => (
            <Box key={i} marginLeft={2}>
              <Text color="redBright">• {w}</Text>
            </Box>
          ))}
        </Box>
      )}

      {/* Recommendations */}
      {result.recommendations.length > 0 && (
        <Box flexDirection="column" marginBottom={1}>
          <Text bold color="yellow">💡 RECOMMENDATIONS:</Text>
          {result.recommendations.map((r, i) => (
            <Box key={i} marginLeft={2}>
              <Text color="yellowBright">→ {r}</Text>
            </Box>
          ))}
        </Box>
      )}

      {/* Missing keywords */}
      {result.missing_keywords.length > 0 && (
        <Box flexDirection="column" marginBottom={1}>
          <Text bold color="red">✗ MISSING KEYWORDS ({result.missing_keywords.length}):</Text>
          <Box marginLeft={2} flexWrap="wrap">
            <Text color="redBright">{result.missing_keywords.join('  ·  ')}</Text>
          </Box>
        </Box>
      )}

      {/* Matched keywords */}
      {result.matched_keywords.length > 0 && (
        <Box flexDirection="column" marginBottom={1}>
          <Text bold color="green">✓ MATCHED KEYWORDS ({result.matched_keywords.length}):</Text>
          <Box marginLeft={2} flexWrap="wrap">
            <Text color="greenBright">{result.matched_keywords.join('  ·  ')}</Text>
          </Box>
        </Box>
      )}

      <Newline />
      <Text color="gray">Press <Text color="white" bold>b</Text> to go back · <Text color="white" bold>q</Text> to quit</Text>
    </Box>
  );
}

function ResultsList({ data }: { data: AnalysisResponse }) {
  const [selected, setSelected] = useState(0);
  const [viewing, setViewing] = useState<ATSResult | null>(null);
  const { exit } = useApp();

  useInput((input, key) => {
    if (viewing) return;
    if (key.upArrow) setSelected(s => Math.max(0, s - 1));
    if (key.downArrow) setSelected(s => Math.min(data.results.length - 1, s + 1));
    if (key.return || input === 'l') setViewing(data.results[selected]);
    if (input === 'q') exit();
  });

  if (viewing) {
    return <DetailView result={viewing} onBack={() => setViewing(null)} />;
  }

  const avg = Math.round(data.results.reduce((s, r) => s + r.score, 0) / data.results.length);
  const autoRejectCount = data.results.filter(r => r.auto_reject).length;

  return (
    <Box flexDirection="column">
      <Header />

      {/* Parsed CV summary */}
      <Box marginBottom={1}>
        <Text color="gray">Candidate: </Text>
        <Text bold color="white">{data.parsed_cv.name || 'Unknown'}</Text>
        <Text color="gray">  ·  Exp: </Text>
        <Text bold color="cyan">{data.parsed_cv.years_exp}yr</Text>
        <Text color="gray">  ·  Skills: </Text>
        <Text bold color="cyan">{data.parsed_cv.skills.length}</Text>
        <Text color="gray">  ·  Avg Score: </Text>
        <Text bold color={scoreBarColor(avg) as any}>{avg}%</Text>
        {autoRejectCount > 0 && (
          <>
            <Text color="gray">  ·  </Text>
            <Text bold color="red">⛔ {autoRejectCount} auto-reject risk{autoRejectCount > 1 ? 's' : ''}</Text>
          </>
        )}
      </Box>

      <Text bold color="gray">{'─'.repeat(70)}</Text>
      <Newline />

      {data.results.map((r, i) => (
        <Box key={r.platform} flexDirection="row">
          <Text color={i === selected ? 'cyan' : 'gray'}>{i === selected ? '▶ ' : '  '}</Text>
          <ResultRow result={r} />
        </Box>
      ))}

      <Newline />
      <Text bold color="gray">{'─'.repeat(70)}</Text>
      <Text color="gray">
        <Text color="cyan" bold>↑↓</Text> navigate  ·  <Text color="cyan" bold>Enter/l</Text> view details  ·  <Text color="cyan" bold>q</Text> quit
      </Text>
    </Box>
  );
}

function Loading({ message }: { message: string }) {
  const [dots, setDots] = useState('');
  useEffect(() => {
    const t = setInterval(() => setDots(d => d.length >= 3 ? '' : d + '.'), 300);
    return () => clearInterval(t);
  }, []);
  return (
    <Box flexDirection="column">
      <Header />
      <Text color="cyan">⏳ {message}{dots}</Text>
    </Box>
  );
}

// ── Main App ──────────────────────────────────────────────────────────────────

function App() {
  const [data, setData] = useState<AnalysisResponse | null>(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);
  const [loadingMsg, setLoadingMsg] = useState('Analyzing CV');
  const { exit } = useApp();

  useEffect(() => {
    const cvFile = process.argv[2];
    const jdFile = process.argv[3];

    if (!cvFile || !jdFile) {
      setError('Usage: npx ts-node src/index.tsx <cv-file> <jd-file>');
      setLoading(false);
      return;
    }

    const cvPath = path.resolve(cvFile);
    const jdPath = path.resolve(jdFile);

    if (!fs.existsSync(cvPath)) { setError(`CV file not found: ${cvPath}`); setLoading(false); return; }
    if (!fs.existsSync(jdPath)) { setError(`JD file not found: ${jdPath}`); setLoading(false); return; }

    const cvText = fs.readFileSync(cvPath, 'utf-8');
    const jdText = fs.readFileSync(jdPath, 'utf-8');

    setLoadingMsg(`Analyzing ${path.basename(cvFile)} against all 6 ATS platforms`);

    fetchText(`${API}/api/analyze`, { cv_text: cvText, jd_text: jdText })
      .then(result => { setData(result); setLoading(false); })
      .catch(err => { setError(`API Error: ${err.message}\n\nMake sure the Go server is running: ./ats-server`); setLoading(false); });
  }, []);

  useInput((input) => {
    if (!loading && (input === 'q')) exit();
  });

  if (loading) return <Loading message={loadingMsg} />;
  if (error) return (
    <Box flexDirection="column">
      <Header />
      <Text color="red" bold>✗ Error</Text>
      <Text color="redBright">{error}</Text>
      <Newline />
      <Text color="gray">Press q to quit</Text>
    </Box>
  );
  if (data) return <ResultsList data={data} />;
  return null;
}

// ── Entry ─────────────────────────────────────────────────────────────────────

render(<App />);
