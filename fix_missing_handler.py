with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

handler_code = """
type FixResumeRequest struct {
	CVText   string `json:"cv_text"`
	JDText   string `json:"jd_text"`
	Platform string `json:"platform"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
	APIKey   string `json:"api_key"`
	BaseURL  string `json:"base_url"`
}

func handleFixResumeGenerate(c *gin.Context) {
	var req FixResumeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	prompt := fmt.Sprintf(`You are an expert ATS (Applicant Tracking System) optimizer.
Your goal is to rewrite the provided CV so that it perfectly matches the provided Job Description (JD) and scores 100%%%% on the %s ATS platform.

Guidelines:
- Incorporate keywords seamlessly into the Summary, Experience, Projects, and Skills sections.
- Highlight relevant experience. DO NOT hallucinate fake jobs, degrees, or metrics.
- MATCH ORIGINAL LENGTH: If the uploaded CV is a single page, keep the output single page. Adjust verbosity accordingly.
- CRITICAL: You must use the EXACT LaTeX structural commands provided in the CV TEXT below. Only modify the actual textual content (bullet points, summary, skills) to optimize for the JD.
- CRITICAL: When using \resumeSubheading, you MUST provide exactly 4 arguments in curly braces: {Title}{Dates}{Company}{Location}. DO NOT use \\\\ or | to combine them!
- CRITICAL: Do NOT insert random \\\\ line breaks in normal text or between arguments. LaTeX handles wrapping automatically.
- CRITICAL: Ensure all plain-text special characters like %% and $ are escaped. Do NOT escape the alignment ampersands (&) inside the tabular* environments used by the custom commands!
- CRITICAL: Use only standard ASCII characters. Do not use unicode characters like approx, use ~ instead.

OUTPUT FORMAT:
Do NOT output the preamble (\\documentclass, \\usepackage, etc).
Start your output EXACTLY with \\begin{document} and end with \\end{document}.
Do NOT wrap the output in markdown code blocks.

--- CV TEXT (USE THIS EXACT LATEX STRUCTURE) ---
%s

--- JOB DESCRIPTION ---
%s`, req.Platform, req.CVText, req.JDText)

	llmConfig := LLMConfig{
		Provider: req.Provider,
		Model:    req.Model,
		APIKey:   req.APIKey,
		BaseURL:  req.BaseURL,
	}

	latexCode, err := generateWithLLM(llmConfig, prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	latexCode = extractLatex(latexCode)

	c.JSON(http.StatusOK, gin.H{"latex": latexCode})
}

type CompileRequest struct"""

text = text.replace("type CompileRequest struct", handler_code)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

