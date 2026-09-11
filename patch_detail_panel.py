import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

# We need to completely rewrite DetailPanel to use the new manual flow.
# Let's extract the DetailPanel function
start_idx = text.find('function DetailPanel')
end_idx = text.find('function PixelSelect', start_idx)

if end_idx == -1:
    end_idx = text.find('function SettingsModal', start_idx)

detail_panel = text[start_idx:end_idx]

new_detail_panel = r"""function DetailPanel({ result, onClose, cvText, jdText }: { result: ATSResult; onClose: () => void; cvText: string; jdText: string }) {
  const [compiling, setCompiling] = useState(false)
  const [errorMsg, setErrorMsg] = useState('');
  const [pastedLatex, setPastedLatex] = useState('');
  
  const promptText = `You are an expert ATS (Applicant Tracking System) optimizer.
Your goal is to rewrite the provided CV so that it perfectly matches the provided Job Description (JD) and scores 100% on the ${result.platform} ATS platform.

Guidelines:
- Incorporate keywords seamlessly into the Summary, Experience, Projects, and Skills sections.
- Highlight relevant experience. DO NOT hallucinate fake jobs, degrees, or metrics.
- MATCH ORIGINAL LENGTH: If the uploaded CV is a single page, keep the output single page. Adjust verbosity accordingly.
- CRITICAL: You must use the EXACT LaTeX structural commands provided in the CV TEXT below. Only modify the actual textual content (bullet points, summary, skills) to optimize for the JD.
- CRITICAL: Do NOT invent or hallucinate new LaTeX commands like \\resumeSummary or \\resumeSkills! Use standard text for the summary and standard \\section{} commands.
- CRITICAL FORMAT EXAMPLE:
\\section{Experience}
  \\resumeSubHeadingListStart
    \\resumeSubheading{Title}{Dates}{Company}{Location}
      \\resumeItemListStart
        \\resumeItem{Bullet point 1}
      \\resumeItemListEnd
  \\resumeSubHeadingListEnd
- CRITICAL: \\resumeSubheading is a COMMAND, not an environment! DO NOT use \\begin{resumeSubheading}. You MUST provide exactly 4 arguments in curly braces: {Title}{Dates}{Company}{Location}. DO NOT use \\\\ or | to combine them!
- CRITICAL: ALL \\resumeSubheading and \\resumeProjectHeading commands MUST be wrapped inside \\resumeSubHeadingListStart and \\resumeSubHeadingListEnd!
- CRITICAL: ALL \\resumeItem commands MUST be wrapped inside \\resumeItemListStart and \\resumeItemListEnd!
- CRITICAL: Do NOT insert random \\\\ line breaks in normal text or between arguments. LaTeX handles wrapping automatically.
- CRITICAL: Ensure all plain-text special characters like % and $ are escaped. Do NOT escape the alignment ampersands (&) inside the tabular* environments used by the custom commands!
- CRITICAL: Use only standard ASCII characters. Do not use unicode characters like approx, use ~ instead.

OUTPUT FORMAT:
Do NOT output the preamble (\\documentclass, \\usepackage, etc).
Start your output EXACTLY with \\begin{document} and end with \\end{document}.
Do NOT wrap the output in markdown code blocks.

--- CV TEXT (USE THIS EXACT LATEX STRUCTURE) ---
${cvText}

--- JOB DESCRIPTION ---
${jdText}`;

  const copyPrompt = () => {
    navigator.clipboard.writeText(promptText);
    alert("Prompt copied to clipboard! Paste it into ChatGPT, Claude, or Gemini.");
  };

  const handleCompile = async () => {
    if (!pastedLatex.trim()) {
      setErrorMsg("Please paste the LaTeX output from the AI first.");
      return;
    }
    
    setCompiling(true);
    try {
      const compileRes = await fetch(`${API}/api/fix-resume/compile`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ latex: pastedLatex })
      });
      if (!compileRes.ok) throw new Error(await compileRes.text());
      const blob = await compileRes.blob();

      try {
        const handle = await (window as any).showSaveFilePicker({
          suggestedName: `Optimized_${result.platform}_Resume.pdf`,
          types: [{ description: 'PDF Document', accept: { 'application/pdf': ['.pdf'] } }]
        });
        const writable = await handle.createWritable();
        await writable.write(blob);
        await writable.close();
      } catch (err: any) {
        if (err.name !== 'AbortError') {
          const url = window.URL.createObjectURL(blob);
          const a = document.createElement('a');
          a.href = url;
          a.download = `Optimized_${result.platform}_Resume.pdf`;
          a.click();
        }
      }
      
      onClose();
    } catch (e: any) {
      console.error("Full Error:", e.message);
      setErrorMsg('Failed to compile the CV. Ensure the AI output starts with \\begin{document} and ends with \\end{document}.');
    } finally {
      setCompiling(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4 overflow-y-auto">
      <div className="bg-[#1a1a1a] border-4 border-[#333] p-6 max-w-4xl w-full max-h-[90vh] overflow-y-auto pixel-border relative shadow-[8px_8px_0_rgba(0,0,0,1)]">
        
        <div className="flex justify-between items-start mb-6 border-b-4 border-[#333] pb-4">
          <div>
            <h2 className="font-PressStart text-xl text-yellow-400 mb-2 drop-shadow-[2px_2px_0_#000]">{result.platform} Analysis</h2>
            <div className="flex items-center gap-4">
              <span className={`font-PressStart text-sm ${result.score > 70 ? 'text-green-400' : result.score > 40 ? 'text-yellow-400' : 'text-red-400'}`}>
                SCORE: {result.score}/100
              </span>
            </div>
          </div>
          <button onClick={onClose} className="nes-btn is-error text-[10px]">X</button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
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
        
        <div className="mt-8 border-t-4 border-[#333] pt-6">
          <h2 className="font-PressStart text-[14px] text-yellow-400 mb-4 drop-shadow-[2px_2px_0_#000]">AI MANUAL WORKFLOW</h2>
          
          <div className="flex flex-col md:flex-row gap-6">
            <div className="flex-1 pixel-border p-4 bg-[#111]">
              <h3 className="font-PressStart text-[10px] text-blue-400 mb-4">STEP 1: Copy AI Prompt</h3>
              <p className="font-mono text-xs text-gray-400 mb-4">Click the button below to copy the heavily optimized prompt. Paste it into Claude, ChatGPT, or Gemini.</p>
              <button 
                onClick={copyPrompt}
                className="nes-btn is-primary w-full text-[10px]"
              >
                COPY PROMPT
              </button>
            </div>
            
            <div className="flex-1 pixel-border p-4 bg-[#111] flex flex-col">
              <h3 className="font-PressStart text-[10px] text-purple-400 mb-4">STEP 2: Paste LaTeX & Compile</h3>
              <textarea 
                className="flex-1 bg-[#1a1a1a] text-green-400 font-mono text-xs border-2 border-[#333] p-2 mb-4 focus:outline-none focus:border-purple-400 resize-none min-h-[100px]"
                placeholder="Paste the LaTeX output from the AI here..."
                value={pastedLatex}
                onChange={(e) => setPastedLatex(e.target.value)}
              />
              <button 
                onClick={handleCompile}
                disabled={compiling}
                className={`nes-btn w-full text-[10px] ${compiling ? 'is-disabled' : 'is-success'}`}
              >
                {compiling ? 'COMPILING...' : 'COMPILE PDF'}
              </button>
            </div>
          </div>
        </div>

        {errorMsg && (
          <div className="fixed inset-0 bg-black/90 flex items-center justify-center z-[100]">
            <div className="pixel-border p-8 max-w-md w-full bg-[#111] text-center border-red-500 shadow-[8px_8px_0_#ef4444]">
              <h2 className="font-PressStart text-red-500 text-[16px] mb-4 drop-shadow-[2px_2px_0_#000]">COMPILER ERROR</h2>
              <p className="font-mono text-gray-300 text-sm mb-8 leading-relaxed">{errorMsg}</p>
              <button 
                onClick={() => setErrorMsg('')}
                className="nes-btn is-error w-full font-PressStart text-[10px]"
              >
                ACKNOWLEDGE
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
"""

text = text.replace(detail_panel, new_detail_panel)

# Also remove llmConfig prop being passed to DetailPanel from the main App component
text = text.replace('llmConfig={llmConfig}', '')

# Since SettingsModal is now obsolete, let's remove the settings button
text = text.replace('{showSettings && <SettingsModal config={llmConfig} setConfig={setLlmConfig} onClose={() => setShowSettings(false)} />}', '')
# Keep the API constants but we don't need SettingsModal anymore

with open("web/src/App.tsx", "w") as f:
    f.write(text)

