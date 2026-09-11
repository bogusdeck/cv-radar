import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

# 1. Remove the useEffect timer
timer_regex = r"  useEffect\(\(\) => \{\n    if \(\!fixing\) \{\n      setFixMessage\('FIX RESUME'\);\n      return;\n    \}\n    const msgs = \['BOOTING AGENT\.\.\.', 'REWRITING CV\.\.\.', 'BEATING ATS\.\.\.', 'COMPILING PDF\.\.\.'\];\n    let i = 0;\n    setFixMessage\(msgs\[0\]\);\n    const timer = setInterval\(\(\) => \{\n      i = \(i \+ 1\) \% msgs\.length;\n      setFixMessage\(msgs\[i\]\);\n    \}, 1500\);\n    return \(\) => clearInterval\(timer\);\n  \}, \[fixing\]\);"
text = re.sub(timer_regex, "", text)

# 2. Rewrite the onClick handler
onclick_regex = r"onClick=\{async \(\) => \{\s*setFixing\(true\);\s*try \{\s*const res = await fetch\(`\$\{API\}/api/fix-resume`, \{.*?a\.click\(\);\s*\} catch \(e: any\) \{\s*alert\('Fix failed: ' \+ e\.message\);\s*\} finally \{\s*setFixing\(false\);\s*\}\s*\}\}"

new_onclick = """onClick={async () => {
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
                    const url = URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.href = url;
                    a.download = `Optimized_Resume_${result.platform}.pdf`;
                    a.click();

                    setFixMessage('SUCCESS!');
                    await new Promise(r => setTimeout(r, 1000));
                  } catch (e: any) {
                    alert('Fix failed: ' + e.message);
                  } finally {
                    setFixing(false);
                    setFixMessage('FIX RESUME');
                  }
                }}"""

text = re.sub(onclick_regex, new_onclick, text, flags=re.DOTALL)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

