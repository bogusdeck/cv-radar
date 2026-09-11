import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

state_repl = """  const [fixing, setFixing] = useState(false);
  const [fixMessage, setFixMessage] = useState('FIX RESUME');

  useEffect(() => {
    if (!fixing) {
      setFixMessage('FIX RESUME');
      return;
    }
    const msgs = ['BOOTING AGENT...', 'REWRITING CV...', 'BEATING ATS...', 'COMPILING PDF...'];
    let i = 0;
    setFixMessage(msgs[0]);
    const timer = setInterval(() => {
      i = (i + 1) % msgs.length;
      setFixMessage(msgs[i]);
    }, 1500);
    return () => clearInterval(timer);
  }, [fixing]);"""

text = text.replace("  const [fixing, setFixing] = useState(false);", state_repl)
text = text.replace("{fixing ? 'FIXING...' : 'FIX RESUME'}", "{fixMessage}")

with open("web/src/App.tsx", "w") as f:
    f.write(text)

