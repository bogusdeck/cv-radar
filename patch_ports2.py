import re

with open("cmd/server/main.go", "r") as f:
    text = f.read()
text = text.replace('port = "8081"', 'port = "8085"')
with open("cmd/server/main.go", "w") as f:
    f.write(text)

with open("web/src/App.tsx", "r") as f:
    text = f.read()
text = text.replace("http://localhost:8081", "http://localhost:8085")
with open("web/src/App.tsx", "w") as f:
    f.write(text)

with open("start.sh", "r") as f:
    text = f.read()
text = text.replace("8081", "8085")
# Also need to add lsof kill for 8085 in start.sh
text = text.replace("lsof -ti:8080 | xargs kill -9 2>/dev/null", "lsof -ti:8080 | xargs kill -9 2>/dev/null\nlsof -ti:8085 | xargs kill -9 2>/dev/null")
with open("start.sh", "w") as f:
    f.write(text)

