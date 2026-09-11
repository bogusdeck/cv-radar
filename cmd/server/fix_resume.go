package main


import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const Preamble = `%-------------------------
% Resume in Latex
%------------------------
\documentclass[letterpaper,10.8pt]{article}
\usepackage{latexsym}
\usepackage[empty]{fullpage}
\usepackage{titlesec}
\usepackage{marvosym}
\usepackage[usenames,dvipsnames]{color}
\usepackage{verbatim}
\usepackage{enumitem}
\usepackage[hidelinks]{hyperref}
\usepackage{fancyhdr}
\usepackage[english]{babel}
\usepackage{tabularx}
\pagestyle{fancy}
\fancyhf{}
\fancyfoot{}
\renewcommand{\headrulewidth}{0pt}
\renewcommand{\footrulewidth}{0pt}
\addtolength{\oddsidemargin}{-0.5in}
\addtolength{\evensidemargin}{-0.5in}
\addtolength{\textwidth}{1in}
\addtolength{\topmargin}{-0.5in}
\addtolength{\textheight}{1.1in}
\urlstyle{same}
\raggedbottom
\raggedright
\setlength{\tabcolsep}{0in}
\titleformat{\section}{
  \vspace{-6pt}\bfseries\scshape\raggedright\large
}{}{0em}{}[\color{black}\titlerule \vspace{-5pt}]

%-------------------------
% Custom commands
\newcommand{\resumeItem}[1]{
  \item\small{
    {#1 \vspace{-2pt}}
  }
}
\newcommand{\resumeSubheading}[4]{
  \vspace{-2pt}\item
    \begin{tabular*}{0.97\textwidth}[t]{l@{\extracolsep{\fill}}r}
      \textbf{#1} & #2 \\
      \textit{\small#3} & \textit{\small #4} \\
    \end{tabular*}\vspace{-7pt}
}
\newcommand{\resumeSubSubheading}[2]{
    \item
    \begin{tabular*}{0.97\textwidth}{l@{\extracolsep{\fill}}r}
      \textit{\small#1} & \textit{\small #2} \\
    \end{tabular*}\vspace{-7pt}
}
\newcommand{\resumeProjectHeading}[2]{
    \item
    \begin{tabular*}{0.97\textwidth}{l@{\extracolsep{\fill}}r}
      \small#1 & #2 \\
    \end{tabular*}\vspace{-7pt}
}
\newcommand{\resumeSubItem}[1]{\resumeItem{#1}\vspace{-4pt}}
\renewcommand\labelitemii{$\vcenter{\hbox{\tiny$\bullet$}}$}
\newcommand{\resumeSubHeadingListStart}{\begin{itemize}[leftmargin=0.15in, label={}]}
\newcommand{\resumeSubHeadingListEnd}{\end{itemize}}
\newcommand{\resumeItemListStart}{\begin{itemize}[topsep=2pt, itemsep=0pt, parsep=0pt]}
\newcommand{\resumeItemListEnd}{\end{itemize}\vspace{-5pt}}

`

func extractLatex(output string) string {
	start := strings.Index(output, "\\begin{document}")
	if start != -1 {
		endStr := "\\end{document}"
		end := strings.LastIndex(output, endStr)
		if end != -1 {
			output = output[start : end+len(endStr)]
		} else {
			output = output[start:]
		}
	} else {
	    // If LLM still outputs \documentclass, find it
	    startDoc := strings.Index(output, "\\documentclass")
	    if startDoc != -1 {
	        return output[startDoc:]
	    }
	}
	
	// Sanitize common LLM hallucinations
	output = regexp.MustCompile(`(?s)\\resumeSummary\{(.*?)\}`).ReplaceAllString(output, `$1`)
	output = regexp.MustCompile(`(?s)\\resumeProjectDescription\{(.*?)\}`).ReplaceAllString(output, `\\resumeItem{$1}`)
	output = regexp.MustCompile(`(?s)\\resumeSkillsBullet\{(.*?)\}`).ReplaceAllString(output, `\\resumeItem{$1}`)
	output = strings.ReplaceAll(output, "\\resumeProjectStart", "\\resumeSubHeadingListStart")
	output = strings.ReplaceAll(output, "\\resumeProjectEnd", "\\resumeSubHeadingListEnd")
	output = strings.ReplaceAll(output, "\\resumeProjectItem", "\\resumeItem")
	output = strings.ReplaceAll(output, "\\resumeSkillItem", "\\resumeItem")
	output = strings.ReplaceAll(output, "\\resumeExperienceItem", "\\resumeItem")
	output = regexp.MustCompile(`\\resumeSubheading\{\s*\}\{\s*\}\{\s*\}(\{\s*\})?`).ReplaceAllString(output, `\\item`)
	
	// Escape unescaped percent signs (which act as comments and eat closing braces)
	output = regexp.MustCompile(`([^\\])%`).ReplaceAllString(output, `$1\%`)
	// Escape ampersands because the LLM shouldn't be generating tabular alignment tabs anyway.
	// (The structural & are safe inside the Preamble constant)
	output = regexp.MustCompile(`([^\\])&`).ReplaceAllString(output, `$1\&`)
	
	// Remove common problematic unicode characters
	output = strings.ReplaceAll(output, "≈", "~")
	output = strings.ReplaceAll(output, "∼", "~") // U+223C
	output = strings.ReplaceAll(output, "–", "--") // en-dash
	output = strings.ReplaceAll(output, "—", "---") // em-dash

	return Preamble + output
}


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
- CRITICAL: Do NOT invent or hallucinate new LaTeX commands like \resumeSummary or \resumeSkills! Use standard text for the summary and standard \section{} commands.
- CRITICAL FORMAT EXAMPLE:
\section{Experience}
  \resumeSubHeadingListStart
    \resumeSubheading{Title}{Dates}{Company}{Location}
      \resumeItemListStart
        \resumeItem{Bullet point 1}
      \resumeItemListEnd
  \resumeSubHeadingListEnd
- CRITICAL: \resumeSubheading is a COMMAND, not an environment! DO NOT use \begin{resumeSubheading}. You MUST provide exactly 4 arguments in curly braces: {Title}{Dates}{Company}{Location}. DO NOT use \\ or | to combine them!
- CRITICAL: Do NOT insert random \\ line breaks in normal text or between arguments. LaTeX handles wrapping automatically.
- CRITICAL: Ensure all plain-text special characters like %% and $ are escaped. Do NOT escape the alignment ampersands (&) inside the tabular* environments used by the custom commands!
- CRITICAL: Use only standard ASCII characters. Do not use unicode characters like approx, use ~ instead.

OUTPUT FORMAT:
Do NOT output the preamble (\documentclass, \usepackage, etc).
Start your output EXACTLY with \begin{document} and end with \end{document}.
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

type CompileRequest struct {
	Latex string `json:"latex"`
}

func handleFixResumeCompile(c *gin.Context) {
	var req CompileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	tmpDir, err := os.MkdirTemp("", "cv-fix-*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create temp directory"})
		return
	}
	defer os.RemoveAll(tmpDir)

	// Always run extractLatex on incoming raw LaTeX to ensure preamble is attached and sanitized
	finalLatex := extractLatex(req.Latex)
	
	texPath := filepath.Join(tmpDir, "resume.tex")
	if err := os.WriteFile(texPath, []byte(finalLatex), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not write tex file"})
		return
	}

	cmd := exec.Command("tectonic", texPath)
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Tectonic error output:", string(out))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compile LaTeX to PDF: " + string(out)})
		return
	}

	pdfPath := filepath.Join(tmpDir, "resume.pdf")
	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read generated PDF"})
		return
	}

	// Silently save a copy directly to the local filesystem (outputs/ directory)
	_ = os.MkdirAll("outputs", 0755)
	localPath := filepath.Join("outputs", "Optimized_Resume.pdf")
	_ = os.WriteFile(localPath, pdfBytes, 0644)

	c.Header("Content-Disposition", "attachment; filename=Optimized_Resume.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
