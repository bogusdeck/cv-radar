import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

# Update FixResumeRequest struct
struct_regex = r"type FixResumeRequest struct \{.*?Platform string `json:\"platform\"`\n\}"
new_struct = """type FixResumeRequest struct {
	CVText   string `json:"cv_text"`
	JDText   string `json:"jd_text"`
	Platform string `json:"platform"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
	APIKey   string `json:"api_key"`
	BaseURL  string `json:"base_url"`
}"""
text = re.sub(struct_regex, new_struct, text, flags=re.DOTALL)

# Find prompt := ... logic and the API call, replace with generateWithLLM
api_call_regex = r"	reqBody, _ := json\.Marshal\(map.*?latexCode := geminiResp\.Candidates\[0\]\.Content\.Parts\[0\]\.Text"
new_api_call = """	llmConfig := LLMConfig{
		Provider: req.Provider,
		Model:    req.Model,
		APIKey:   req.APIKey,
		BaseURL:  req.BaseURL,
	}

	latexCode, err := generateWithLLM(llmConfig, prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}"""
text = re.sub(api_call_regex, new_api_call, text, flags=re.DOTALL)

# Remove old GEMINI_API_KEY check
key_check_regex = r"\s*apiKey := os\.Getenv\(\"GEMINI_API_KEY\"\).*?return\s*\}"
text = re.sub(key_check_regex, "", text, flags=re.DOTALL)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

