import re

with open("web/src/App.tsx", "r") as f:
    content = f.read()

# I want to wrap everything inside the p-4 md:p-8 div inside <div className="relative z-10 flex flex-col h-full gap-4 md:gap-6">
# Let's just do a clean replacement

new_input_section = """
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
              <DitherBackground />
              
              <div className="relative z-10 p-4 md:p-8 flex-1 flex flex-col gap-4 md:gap-6 h-full">
                <h2 className="text-center font-Gumball text-2xl md:text-4xl text-yellow-400 drop-shadow-[4px_4px_0_#000] tracking-widest shrink-0">SCAN YOUR CV</h2>
                
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
"""

# replace everything from {/* Input Section */} up to <main className="max-w-[95%] 2xl:max-w-[1400px] mx-auto flex flex-col px-2 md:px-8 pb-12 pt-12">
start_marker = "{/* Input Section */}"
end_marker = "<main className="

start_idx = content.find(start_marker)
end_idx = content.find(end_marker)

if start_idx != -1 and end_idx != -1:
    new_content = content[:start_idx] + new_input_section.strip() + "\n\n      " + content[end_idx:]
    with open("web/src/App.tsx", "w") as f:
        f.write(new_content)
else:
    print("Could not find markers")
