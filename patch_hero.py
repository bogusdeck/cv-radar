import re

with open("web/src/App.tsx", "r") as f:
    content = f.read()

# The hero section starts like this:
target = """      {/* Hero Section */}
      <div className="relative mb-12 min-h-screen flex flex-col items-center justify-center px-4 overflow-hidden">
        
        {/* Top Left Logo */}"""

replacement = """      {/* Hero Section */}
      <div className="relative mb-12 min-h-screen flex flex-col items-center justify-center px-4 overflow-hidden">
        
        <DitherBackground />
        <div className="scanlines absolute inset-0 w-full h-full z-10 mix-blend-overlay opacity-80 pointer-events-none"></div>
        <div className="absolute inset-0 w-full h-full z-10 bg-gradient-to-b from-[#0a0a0a]/40 via-transparent to-[#0a0a0a]/90 pointer-events-none"></div>
        
        {/* Top Left Logo */}"""

if target in content:
    content = content.replace(target, replacement)
else:
    # try generic replace if whitespace changed
    print("Could not find exact hero target")

with open("web/src/App.tsx", "w") as f:
    f.write(content)
