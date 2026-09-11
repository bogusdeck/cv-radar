import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

# We will completely rewrite fix_resume.go to have two endpoints
new_go_code = """package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

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

	prompt := fmt.Sprintf("You are an expert ATS (Applicant Tracking System) optimizer.\\nYour goal is to rewrite the provided CV so that it perfectly matches the provided Job Description (JD) and scores 100%%%% on the %s ATS platform.\\n\\nGuidelines:\\n- Incorporate keywords seamlessly.\\n- Highlight relevant experience.\\n- DO NOT hallucinate fake jobs or degrees, but adapt the phrasing and bullet points to match the JD's terminology exactly.\\n- Keep it professional and ATS-friendly (no weird columns or graphics).\\n\\nOUTPUT FORMAT:\\nYou MUST output ONLY a valid, complete, and compilable LaTeX document.\\nUse a standard class like \\\\documentclass[11pt,a4paper]{article}.\\nInclude \\\\usepackage{geometry}, \\\\usepackage{hyperref}, etc.\\nEnsure all special characters are escaped properly.\\nDo NOT wrap the output in markdown code blocks (```latex ... ```). Output raw LaTeX text starting with \\\\documentclass and ending with \\\\end{document}.\\n\\n--- CV TEXT ---\\n%s\\n\\n--- JOB DESCRIPTION ---\\n%s", req.Platform, req.CVText, req.JDText)

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

	texPath := filepath.Join(tmpDir, "resume.tex")
	if err := os.WriteFile(texPath, []byte(req.Latex), 0644); err != nil {
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

	c.Header("Content-Disposition", "attachment; filename=Optimized_Resume.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
"""

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(new_go_code)

