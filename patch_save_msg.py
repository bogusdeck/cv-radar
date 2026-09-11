import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

target = """                    const blob = await compileRes.blob();
                    try {"""

new_code = """                    const blob = await compileRes.blob();
                    setFixMessage('CHOOSE SAVE LOCATION...');
                    try {"""

text = text.replace(target, new_code)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

