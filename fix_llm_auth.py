import re

with open("cmd/server/llm.go", "r") as f:
    text = f.read()

target = r"""	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}"""

new_code = r"""	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("OLLAMA_API_KEY")
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}"""

text = text.replace(target, new_code)

if '"os"' not in text:
    text = text.replace('"net/http"', '"net/http"\n\t"os"')

with open("cmd/server/llm.go", "w") as f:
    f.write(text)

