#!/usr/bin/env node
'use strict';

import React, { useState, useEffect } from 'react';
import { render, Box, Text, useInput, useApp, Newline } from 'ink';
import * as fs from 'fs';
import * as path from 'path';
import * as http from 'http';
import * as readline from 'readline';

// ─── Config ──────────────────────────────────────────────────────────────────

const API_BASE = process.env.ATS_API || 'http://localhost:8080';

// ─── Types ───────────────────────────────────────────────────────────────────

interface Breakdown {
  category: string;
  score: number;
  max_score: number;
  weight: number;
}

interface ATSResult {
  platform: string;
  vendor: string;
  score: number;
  score_low: number;
  score_high: number;
  grade: string;
  confidence: number;
  simulation_note: string;
  keyword_score: number;
  structure_score: number;
  experience_score: number;
  education_score: number;
  matched_keywords: string[];
  missing_keywords: string[];
  warnings: string[];
  recommendations: string[];
  auto_reject: boolean;
  breakdown: Breakdown[];
}

interface ParsedCV {
  name: string;
  email: string;
  phone: string;
  skills: string[];
  years_exp: number;
  experience: { title: string; company: string; duration: string }[];
  education: { degree: string; institution: string }[];
}

interface AnalysisResponse {
  parsed_cv: ParsedCV;
  results: ATSResult[];
}

// ─── API ─────────────────────────────────────────────────────────────────────

function apiPost(endpoint: string, body: object): Promise<any> {
  return new Promise((resolve, reject) => {
    const payload = JSON.stringify(body);
    const urlObj = new URL(API_BASE + endpoint);
    const options: http.RequestOptions = {
      hostname: urlObj.hostname,
      port: parseInt(urlObj.port || '80'),
      path: urlObj.pathname,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Content-Length': Buffer.byteLength(payload),
      },
    };
    const req = http.request(options, (res) => {
      let data = '';
      res.on('data', (chunk) => (data += chunk));
      res.on('end', () => {
        try {
          resolve(JSON.parse(data));
        } catch {
          reject(new Error(`Invalid JSON response from API`));
        }
      });
    });
    req.on('error', (e) => reject(new Error(`Cannot connect to ATS server: ${e.message}\n→ Run: ./ats-server`)));
    req.write(payload);
    req.end();
  });
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

type InkColor =
  | 'black' | 'red' | 'green' | 'yellow' | 'blue' | 'magenta' | 'cyan' | 'white' | 'gray'
  | 'redBright' | 'greenBright' | 'yellowBright' | 'blueBright' | 'magentaBright' | 'cyanBright' | 'whiteBright';

function gradeColor(grade: string): InkColor {
  const m: Record<string, InkColor> = { A: 'greenBright', B: 'green', C: 'yellowBright', D: 'yellow', F: 'redBright' };
  return m[grade] ?? 'white';
}

function scoreColor(score: number): InkColor {
  if (score >= 80) return 'greenBright';
  if (score >= 65) return 'green';
  if (score >= 50) return 'yellowBright';
  if (score >= 35) return 'yellow';
  return 'redBright';
}

function confColor(conf: number): InkColor {
  if (conf >= 65) return 'greenBright';
  if (conf >= 45) return 'yellowBright';
  return 'redBright';
}

function bar(score: number, width = 22): string {
  const filled = Math.round((Math.min(score, 100) / 100) * width);
  return '█'.repeat(filled) + '░'.repeat(width - filled);
}

function wrap(text: string, maxWidth: number): string[] {
  const words = text.split(' ');
  const lines: string[] = [];
  let current = '';
  for (const word of words) {
    if ((current + (current ? ' ' : '') + word).length > maxWidth) {
      if (current) lines.push(current);
      current = word;
    } else {
      current = current ? current + ' ' + word : word;
    }
  }
  if (current) lines.push(current);
  return lines;
}

// ─── Components ──────────────────────────────────────────────────────────────

function Header({ subtitle }: { subtitle?: string }) {
  return (
    <Box flexDirection="column" marginBottom={1}>
      <Text bold color="magenta">{'╔' + '═'.repeat(62) + '╗'}</Text>
      <Box>
        <Text bold color="magenta">{'║  '}</Text>
        <Text bold color="white">🎯  ATS SCANNER</Text>
        <Text color="gray">  ·  6-Platform CV Analyzer  ·  React Ink TUI</Text>
        <Text bold color="magenta">{'  ║'}</Text>
      </Box>
      {subtitle && (
        <Box>
          <Text bold color="magenta">{'║  '}</Text>
          <Text color="cyan">{subtitle.padEnd(59)}</Text>
          <Text bold color="magenta">{'║'}</Text>
        </Box>
      )}
      <Text bold color="magenta">{'╚' + '═'.repeat(62) + '╝'}</Text>
    </Box>
  );
}

function Divider({ label }: { label?: string }) {
  if (!label) return <Text color="gray">{'─'.repeat(64)}</Text>;
  const padded = `─── ${label} `;
  return <Text color="gray">{padded + '─'.repeat(Math.max(0, 64 - padded.length))}</Text>;
}

function Spinner({ msg }: { msg: string }) {
  const frames = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'];
  const [frame, setFrame] = useState(0);
  useEffect(() => {
    const t = setInterval(() => setFrame((f) => (f + 1) % frames.length), 80);
    return () => clearInterval(t);
  }, []);
  return (
    <Box>
      <Text color="cyan" bold>{frames[frame]} </Text>
      <Text color="white">{msg}</Text>
    </Box>
  );
}

// ── Overview screen ───────────────────────────────────────────────────────────

function Overview({
  data,
  selected,
  onSelect,
}: {
  data: AnalysisResponse;
  selected: number;
  onSelect: (r: ATSResult) => void;
}) {
  const { exit } = useApp();
  const [sel, setSel] = useState(selected);

  useInput((input, key) => {
    if (key.upArrow) setSel((s) => Math.max(0, s - 1));
    if (key.downArrow) setSel((s) => Math.min(data.results.length - 1, s + 1));
    if (key.return || input === 'l' || input === 'o') onSelect(data.results[sel]);
    if (input === 'q') exit();
  });

  const avg = Math.round(data.results.reduce((s, r) => s + r.score, 0) / data.results.length);
  const best = data.results.reduce((a, b) => (a.score > b.score ? a : b));
  const worst = data.results.reduce((a, b) => (a.score < b.score ? a : b));
  const rejCount = data.results.filter((r) => r.auto_reject).length;
  const cv = data.parsed_cv;

  return (
    <Box flexDirection="column">
      <Header subtitle={`Candidate: ${cv.name || 'Unknown'}  ·  ${cv.years_exp}yr exp  ·  ${cv.skills.length} skills`} />

      {/* Stats row */}
      <Box marginBottom={1} gap={3}>
        <Box flexDirection="column" alignItems="center">
          <Text color="gray" dimColor>AVG SCORE</Text>
          <Text bold color={scoreColor(avg)}>{avg}%</Text>
        </Box>
        <Box flexDirection="column" alignItems="center">
          <Text color="gray" dimColor>BEST MATCH</Text>
          <Text bold color="greenBright">{best.platform}</Text>
        </Box>
        <Box flexDirection="column" alignItems="center">
          <Text color="gray" dimColor>TOUGHEST</Text>
          <Text bold color="redBright">{worst.platform}</Text>
        </Box>
        <Box flexDirection="column" alignItems="center">
          <Text color="gray" dimColor>AUTO-REJECT</Text>
          <Text bold color={rejCount > 0 ? 'redBright' : 'greenBright'}>
            {rejCount > 0 ? `⛔ ${rejCount}` : '✓ None'}
          </Text>
        </Box>
        <Box flexDirection="column" alignItems="center">
          <Text color="gray" dimColor>EMAIL</Text>
          <Text color="cyan">{cv.email || '—'}</Text>
        </Box>
      </Box>

      <Divider label="PLATFORM SCORES" />
      <Newline />

      {/* Results list */}
      {data.results.map((r, i) => {
        const isSel = i === sel;
        const gc = gradeColor(r.grade);
        const sc = scoreColor(r.score);
        const cc = confColor(r.confidence);
        const lowHigh = `${Math.round(r.score_low ?? r.score - 10)}–${Math.round(r.score_high ?? r.score + 10)}`;

        return (
          <Box key={r.platform} flexDirection="column" marginBottom={1}>
            <Box>
              <Text color={isSel ? 'cyan' : 'gray'} bold>{isSel ? '▶ ' : '  '}</Text>
              <Text bold color="white">{r.platform.padEnd(15)}</Text>
              <Text bold color={gc}> [{r.grade}] </Text>
              <Text color={sc}>{bar(r.score)}</Text>
              <Text bold color={sc}> {String(Math.round(r.score)).padStart(3)}%</Text>
              <Text color="gray">  ~{Math.round(r.confidence)}% conf </Text>
              <Text color={cc} dimColor>({lowHigh})</Text>
              {r.auto_reject && <Text color="redBright" bold>  ⛔ AUTO-REJECT</Text>}
            </Box>
            <Box marginLeft={2}>
              <Text color="gray" dimColor>{'by ' + r.vendor.padEnd(12)}</Text>
              <Text color="greenBright">  ✓ {(r.matched_keywords ?? []).length} matched</Text>
              <Text color="redBright">  ✗ {(r.missing_keywords ?? []).length} missing</Text>
              <Text color="gray" dimColor>  kw:{Math.round(r.keyword_score ?? 0)}% str:{Math.round(r.structure_score ?? 0)}% exp:{Math.round(r.experience_score ?? 0)}%</Text>
            </Box>
          </Box>
        );
      })}

      <Newline />
      <Divider />
      <Box gap={2} marginTop={1}>
        <Text color="gray"><Text color="cyan" bold>↑↓</Text> navigate</Text>
        <Text color="gray"><Text color="cyan" bold>Enter</Text> details</Text>
        <Text color="gray"><Text color="cyan" bold>q</Text> quit</Text>
      </Box>
    </Box>
  );
}

// ── Detail screen ─────────────────────────────────────────────────────────────

type DetailTab = 'breakdown' | 'keywords' | 'advice';

function Detail({ result, onBack }: { result: ATSResult; onBack: () => void }) {
  const [tab, setTab] = useState<DetailTab>('breakdown');
  const tabs: DetailTab[] = ['breakdown', 'keywords', 'advice'];

  useInput((input, key) => {
    if (input === 'b' || key.escape) onBack();
    if (input === '1') setTab('breakdown');
    if (input === '2') setTab('keywords');
    if (input === '3') setTab('advice');
    if (key.tab) setTab((t) => tabs[(tabs.indexOf(t) + 1) % tabs.length]);
  });

  const gc = gradeColor(result.grade);
  const cc = confColor(result.confidence);
  const low = Math.round(result.score_low ?? result.score - 10);
  const high = Math.round(result.score_high ?? result.score + 10);

  return (
    <Box flexDirection="column">
      <Header subtitle={`${result.platform}  ·  by ${result.vendor}`} />

      {/* Score header */}
      <Box marginBottom={1} gap={4}>
        <Box flexDirection="column" alignItems="center">
          <Text bold color={gradeColor(result.grade)} >{`  ${Math.round(result.score)}%  `}</Text>
          <Text bold color={gc}>Grade {result.grade}</Text>
        </Box>
        <Box flexDirection="column">
          <Box>
            <Text color="gray">Confidence: </Text>
            <Text bold color={cc}>~{Math.round(result.confidence)}%</Text>
          </Box>
          <Box>
            <Text color="gray">Score range: </Text>
            <Text color="yellowBright">{low}–{high}%</Text>
            <Text color="gray" dimColor>  (likely real score)</Text>
          </Box>
          {result.auto_reject && (
            <Text bold color="redBright">⛔ AUTO-REJECT RISK</Text>
          )}
        </Box>
        <Box flexDirection="column">
          <Text color={scoreColor(result.score)}>{bar(result.score, 30)}</Text>
        </Box>
      </Box>

      {/* Simulation note */}
      <Box
        marginBottom={1}
        paddingLeft={1}
        borderStyle="single"
        borderColor="yellow"
      >
        <Box flexDirection="column">
          <Text color="yellow" bold>⚠ Accuracy Note</Text>
          {wrap(result.simulation_note || '', 56).map((line, i) => (
            <Text key={i} color="gray" dimColor>{line}</Text>
          ))}
        </Box>
      </Box>

      {/* Tabs */}
      <Box marginBottom={1} gap={1}>
        {tabs.map((t, i) => (
          <Box key={t} paddingX={2}>
            <Text bold color={tab === t ? 'cyan' : 'gray'}>
              {`[${i + 1}] ${t.toUpperCase()}`}
            </Text>
          </Box>
        ))}
        <Text color="gray" dimColor>  Tab to switch</Text>
      </Box>
      <Divider />
      <Newline />

      {/* Breakdown tab */}
      {tab === 'breakdown' && (
        <Box flexDirection="column">
          {(result.breakdown ?? []).map((b) => (
            <Box key={b.category} marginBottom={1}>
              <Text color="gray">{b.category.padEnd(36)}</Text>
              <Text color={scoreColor(b.score)}>{bar(b.score, 16)}</Text>
              <Text bold color={scoreColor(b.score)}> {String(Math.round(b.score)).padStart(3)}%</Text>
              <Text color="gray" dimColor>  wt:{Math.round(b.weight * 100)}%</Text>
            </Box>
          ))}
        </Box>
      )}

      {/* Keywords tab */}
      {tab === 'keywords' && (
        <Box flexDirection="column" gap={1}>
          <Box flexDirection="column">
            <Text bold color="greenBright">✓ MATCHED ({(result.matched_keywords ?? []).length})</Text>
            <Box flexWrap="wrap" marginLeft={2}>
              {(result.matched_keywords ?? []).map((kw) => (
                <Text key={kw} color="green"> {kw}  </Text>
              ))}
            </Box>
          </Box>
          <Newline />
          <Box flexDirection="column">
            <Text bold color="redBright">✗ MISSING ({(result.missing_keywords ?? []).length})</Text>
            <Box flexWrap="wrap" marginLeft={2}>
              {(result.missing_keywords ?? []).map((kw) => (
                <Text key={kw} color="red"> {kw}  </Text>
              ))}
            </Box>
          </Box>
        </Box>
      )}

      {/* Advice tab */}
      {tab === 'advice' && (
        <Box flexDirection="column" gap={1}>
          {(result.warnings ?? []).length > 0 && (
            <Box flexDirection="column">
              <Text bold color="redBright">⚠ WARNINGS</Text>
              {(result.warnings ?? []).map((w, i) => (
                <Box key={i} marginLeft={2}>
                  {wrap(w, 58).map((line, j) => (
                    <Text key={j} color="yellow">{j === 0 ? '• ' : '  '}{line}</Text>
                  ))}
                </Box>
              ))}
            </Box>
          )}
          <Box flexDirection="column">
            <Text bold color="cyanBright">💡 RECOMMENDATIONS</Text>
            {(result.recommendations ?? []).map((rec, i) => (
              <Box key={i} marginLeft={2} flexDirection="column">
                {wrap(rec, 58).map((line, j) => (
                  <Text key={j} color="white">{j === 0 ? '→ ' : '  '}{line}</Text>
                ))}
              </Box>
            ))}
          </Box>
        </Box>
      )}

      <Newline />
      <Divider />
      <Box gap={2} marginTop={1}>
        <Text color="gray"><Text color="cyan" bold>1/2/3</Text> or <Text color="cyan" bold>Tab</Text> switch tabs</Text>
        <Text color="gray"><Text color="cyan" bold>b/Esc</Text> back</Text>
      </Box>
    </Box>
  );
}

// ─── Input Mode ──────────────────────────────────────────────────────────────

// Read multi-line stdin until Ctrl+D
function readStdin(): Promise<string> {
  return new Promise((resolve) => {
    const rl = readline.createInterface({ input: process.stdin, terminal: false });
    const lines: string[] = [];
    rl.on('line', (l) => lines.push(l));
    rl.on('close', () => resolve(lines.join('\n')));
  });
}

// ─── Loading Screen ──────────────────────────────────────────────────────────

function LoadingScreen({ msg }: { msg: string }) {
  return (
    <Box flexDirection="column">
      <Header />
      <Box marginTop={1}>
        <Spinner msg={msg} />
      </Box>
    </Box>
  );
}

// ─── Error Screen ────────────────────────────────────────────────────────────

function ErrorScreen({ msg }: { msg: string }) {
  const { exit } = useApp();
  useInput((input) => { if (input === 'q') exit(); });
  return (
    <Box flexDirection="column">
      <Header />
      <Box marginTop={1} flexDirection="column">
        <Text bold color="redBright">✗ Error</Text>
        {msg.split('\n').map((line, i) => (
          <Text key={i} color="red">{line}</Text>
        ))}
      </Box>
      <Newline />
      <Text color="gray">Press <Text color="white" bold>q</Text> to quit</Text>
    </Box>
  );
}

// ─── Main App ────────────────────────────────────────────────────────────────

type AppState =
  | { stage: 'loading'; msg: string }
  | { stage: 'error'; msg: string }
  | { stage: 'overview'; data: AnalysisResponse }
  | { stage: 'detail'; data: AnalysisResponse; result: ATSResult };

function App() {
  const [state, setState] = useState<AppState>({ stage: 'loading', msg: 'Starting up…' });

  useEffect(() => {
    (async () => {
      const args = process.argv.slice(2);
      let cvText = '';
      let jdText = '';

      // Mode 1: ats-tui cv.txt jd.txt
      if (args.length >= 2 && !args.includes('--cv') && !args.includes('--jd') && args[0] !== '--stdin') {
        const cvPath = path.resolve(args[0]);
        const jdPath = path.resolve(args[1]);

        if (!fs.existsSync(cvPath)) {
          setState({ stage: 'error', msg: `CV file not found: ${cvPath}` });
          return;
        }
        if (!fs.existsSync(jdPath)) {
          setState({ stage: 'error', msg: `JD file not found: ${jdPath}` });
          return;
        }

        cvText = fs.readFileSync(cvPath, 'utf-8');
        jdText = fs.readFileSync(jdPath, 'utf-8');
        setState({ stage: 'loading', msg: `Parsing ${path.basename(args[0])} against ${path.basename(args[1])}…` });
      }
      // Mode 2: cat cv.txt | ats-tui --stdin "job description..."
      else if (args[0] === '--stdin' && args[1]) {
        setState({ stage: 'loading', msg: 'Reading CV from stdin…' });
        cvText = await readStdin();
        jdText = args.slice(1).join(' ');
      }
      // Mode 3: ats-tui --cv cv.txt --jd "paste JD here"
      else if (args.includes('--cv') && args.includes('--jd')) {
        const cvIdx = args.indexOf('--cv');
        const jdIdx = args.indexOf('--jd');
        const cvPath = path.resolve(args[cvIdx + 1]);
        if (!fs.existsSync(cvPath)) {
          setState({ stage: 'error', msg: `CV file not found: ${cvPath}` });
          return;
        }
        cvText = fs.readFileSync(cvPath, 'utf-8');
        jdText = args.slice(jdIdx + 1).join(' ');
        setState({ stage: 'loading', msg: `Analyzing ${path.basename(args[cvIdx + 1])}…` });
      }
      else {
        setState({
          stage: 'error',
          msg: [
            'Usage:',
            '',
            '  # Analyze files:',
            '  ats-tui <cv-file> <jd-file>',
            '',
            '  # Pipe CV, pass JD as arg:',
            '  cat cv.txt | ats-tui --stdin "job description text"',
            '',
            '  # CV file + inline JD:',
            '  ats-tui --cv cv.tex --jd "job description text"',
            '',
            '  # Set custom API server:',
            '  ATS_API=http://myserver:8080 ats-tui cv.txt jd.txt',
          ].join('\n'),
        });
        return;
      }

      setState({ stage: 'loading', msg: 'Analyzing against 6 ATS platforms…' });

      try {
        const result: AnalysisResponse = await apiPost('/api/analyze', {
          cv_text: cvText,
          jd_text: jdText,
        });
        setState({ stage: 'overview', data: result });
      } catch (e: any) {
        setState({ stage: 'error', msg: e.message });
      }
    })();
  }, []);

  if (state.stage === 'loading') return <LoadingScreen msg={state.msg} />;
  if (state.stage === 'error') return <ErrorScreen msg={state.msg} />;

  if (state.stage === 'overview') {
    return (
      <Overview
        data={state.data}
        selected={0}
        onSelect={(r) => setState({ stage: 'detail', data: state.data, result: r })}
      />
    );
  }

  if (state.stage === 'detail') {
    return (
      <Detail
        result={state.result}
        onBack={() => setState({ stage: 'overview', data: state.data })}
      />
    );
  }

  return null;
}

// ─── Static Output Runner ─────────────────────────────────────────────────────

function getGradeColor(grade: string): string {
  const m: Record<string, string> = { A: '\x1b[92m', B: '\x1b[32m', C: '\x1b[93m', D: '\x1b[33m', F: '\x1b[91m' };
  return m[grade] ?? '\x1b[37m';
}

function getScoreColor(score: number): string {
  if (score >= 80) return '\x1b[92m';
  if (score >= 65) return '\x1b[32m';
  if (score >= 50) return '\x1b[93m';
  if (score >= 35) return '\x1b[33m';
  return '\x1b[91m';
}

function getConfColor(conf: number): string {
  if (conf >= 65) return '\x1b[92m';
  if (conf >= 45) return '\x1b[93m';
  return '\x1b[91m';
}

function staticBar(score: number, width = 20): string {
  const filled = Math.round((Math.min(score, 100) / 100) * width);
  return '█'.repeat(filled) + '░'.repeat(width - filled);
}

function printHelp() {
  console.log([
    'Usage:',
    '',
    '  # Analyze files:',
    '  ats-tui <cv-file> <jd-file>',
    '',
    '  # Pipe CV, pass JD as arg:',
    '  cat cv.txt | ats-tui --stdin "job description text"',
    '',
    '  # CV file + inline JD:',
    '  ats-tui --cv cv.tex --jd "job description text"',
    '',
    '  # Set custom API server:',
    '  ATS_API=http://myserver:8080 ats-tui cv.txt jd.txt',
  ].join('\n'));
}

function printStaticReport(data: AnalysisResponse) {
  const cv = data.parsed_cv;
  console.log(`\n\x1b[35m╔${'═'.repeat(62)}╗`);
  console.log(`║  \x1b[1m\x1b[37m🎯  ATS SCANNER\x1b[0m\x1b[35m  ·  6-Platform CV Analyzer (Static Mode)    ║`);
  console.log(`╚${'═'.repeat(62)}╝\x1b[0m`);

  console.log(`\x1b[1mCandidate:\x1b[0m ${cv.name || 'Unknown'}  ·  \x1b[36m${cv.years_exp}yr exp\x1b[0m  ·  \x1b[36m${cv.skills.length} skills\x1b[0m`);
  if (cv.email) console.log(`\x1b[1mContact:\x1b[0m ${cv.email} | ${cv.phone || '—'}`);
  console.log(`\x1b[90m─${'─'.repeat(62)}─\x1b[0m`);

  const avg = Math.round(data.results.reduce((s, r) => s + r.score, 0) / data.results.length);
  const best = data.results.reduce((a, b) => (a.score > b.score ? a : b));
  const worst = data.results.reduce((a, b) => (a.score < b.score ? a : b));
  const rejCount = data.results.filter((r) => r.auto_reject).length;

  console.log(`\x1b[1mAVG SCORE:\x1b[0m ${getScoreColor(avg)}${avg}%\x1b[0m | \x1b[1mBEST:\x1b[0m \x1b[92m${best.platform}\x1b[0m | \x1b[1mTOUGHEST:\x1b[0m \x1b[91m${worst.platform}\x1b[0m | \x1b[1mAUTO-REJECT:\x1b[0m ${rejCount > 0 ? `\x1b[91m⛔ ${rejCount}\x1b[0m` : '\x1b[92m✓ None\x1b[0m'}`);
  console.log(`\x1b[90m─${'─'.repeat(62)}─\x1b[0m`);

  console.log(`\x1b[1m\x1b[35mPLATFORM SCORES:\x1b[0m`);
  for (const r of data.results) {
    const gc = getGradeColor(r.grade);
    const sc = getScoreColor(r.score);
    const cc = getConfColor(r.confidence);
    const lowHigh = `${Math.round(r.score_low ?? r.score - 10)}–${Math.round(r.score_high ?? r.score + 10)}`;

    console.log(`  \x1b[1m${r.platform.padEnd(15)}\x1b[0m ${gc}[${r.grade}]\x1b[0m ${sc}${staticBar(r.score)}\x1b[0m ${sc}${String(Math.round(r.score)).padStart(3)}%\x1b[0m  \x1b[90m~${Math.round(r.confidence)}% conf\x1b[0m ${cc}(${lowHigh}%)\x1b[0m${r.auto_reject ? ' \x1b[91m⛔ AUTO-REJECT\x1b[0m' : ''}`);
    console.log(`    \x1b[90mby ${r.vendor.padEnd(12)}\x1b[0m  \x1b[92m✓ ${(r.matched_keywords ?? []).length} keywords\x1b[0m  \x1b[91m✗ ${(r.missing_keywords ?? []).length} missing\x1b[0m`);
  }

  console.log(`\x1b[90m─${'─'.repeat(62)}─\x1b[0m`);
  console.log(`\x1b[1m\x1b[36mDETAILED KEYWORDS & WARNINGS PER PLATFORM:\x1b[0m`);

  for (const r of data.results) {
    console.log(`\n  \x1b[1m🎯 ${r.platform} (${r.vendor})\x1b[0m`);
    if (r.simulation_note) {
      console.log(`    \x1b[33m⚠ Accuracy Note: ${r.simulation_note}\x1b[0m`);
    }
    if ((r.missing_keywords ?? []).length > 0) {
      console.log(`    \x1b[91mMissing Keywords:\x1b[0m ${r.missing_keywords.slice(0, 10).join(', ')}${r.missing_keywords.length > 10 ? '...' : ''}`);
    }
    if ((r.warnings ?? []).length > 0) {
      console.log(`    \x1b[93mWarnings:\x1b[0m`);
      for (const w of r.warnings.slice(0, 3)) {
        console.log(`      • ${w}`);
      }
    }
    if ((r.recommendations ?? []).length > 0) {
      console.log(`    \x1b[36mRecommendations:\x1b[0m`);
      for (const rec of r.recommendations.slice(0, 3)) {
        console.log(`      → ${rec}`);
      }
    }
  }
  console.log(`\n\x1b[90mTip: Run this command in a regular terminal for a fully interactive, tabbed dashboard!\x1b[0m\n`);
}

async function runStaticCLI() {
  const args = process.argv.slice(2);
  let cvText = '';
  let jdText = '';

  if (args.length >= 2 && !args.includes('--cv') && !args.includes('--jd') && args[0] !== '--stdin') {
    const cvPath = path.resolve(args[0]);
    const jdPath = path.resolve(args[1]);

    if (!fs.existsSync(cvPath)) {
      console.error(`Error: CV file not found: ${cvPath}`);
      process.exit(1);
    }
    if (!fs.existsSync(jdPath)) {
      console.error(`Error: JD file not found: ${jdPath}`);
      process.exit(1);
    }

    cvText = fs.readFileSync(cvPath, 'utf-8');
    jdText = fs.readFileSync(jdPath, 'utf-8');
  } else if (args[0] === '--stdin' && args[1]) {
    cvText = await readStdin();
    jdText = args.slice(1).join(' ');
  } else if (args.includes('--cv') && args.includes('--jd')) {
    const cvIdx = args.indexOf('--cv');
    const jdIdx = args.indexOf('--jd');
    const cvPath = path.resolve(args[cvIdx + 1]);
    if (!fs.existsSync(cvPath)) {
      console.error(`Error: CV file not found: ${cvPath}`);
      process.exit(1);
    }
    cvText = fs.readFileSync(cvPath, 'utf-8');
    jdText = args.slice(jdIdx + 1).join(' ');
  } else {
    printHelp();
    process.exit(1);
  }

  console.log('Analyzing against 6 ATS platforms...');
  try {
    const response: AnalysisResponse = await apiPost('/api/analyze', {
      cv_text: cvText,
      jd_text: jdText,
    });
    printStaticReport(response);
  } catch (e: any) {
    console.error(`Error: ${e.message}`);
    process.exit(1);
  }
}

// ─── Entry ───────────────────────────────────────────────────────────────────

const isInteractive = process.stdin.isTTY && typeof process.stdin.setRawMode === 'function';

if (isInteractive) {
  render(<App />, { patchConsole: false });
} else {
  runStaticCLI();
}
