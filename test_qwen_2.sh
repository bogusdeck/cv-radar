curl -s http://localhost:11434/api/generate -d '{
  "model": "qwen2.5:7b-instruct-q4_K_M",
  "prompt": "You must use \\resumeSubheading{Title}{Dates}{Company}{Location}. DO NOT use \\\\. Format this: Software Dev at Instahyre, Noida. Dec 2024 - Jun 2026.",
  "stream": false
}' | jq -r .response
