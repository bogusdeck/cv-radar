import re

# Update Go server port in main.go (if it's hardcoded anywhere, but it reads from PORT env var)
# The default is 8080.
with open("cmd/server/main.go", "r") as f:
    text = f.read()
text = text.replace('port = "8080"', 'port = "8081"')
with open("cmd/server/main.go", "w") as f:
    f.write(text)

# Update App.tsx
with open("web/src/App.tsx", "r") as f:
    text = f.read()
text = text.replace("http://localhost:8080", "http://localhost:8081")
with open("web/src/App.tsx", "w") as f:
    f.write(text)

# Update start.sh
with open("start.sh", "r") as f:
    text = f.read()
text = text.replace("8080", "8081")
with open("start.sh", "w") as f:
    f.write(text)

