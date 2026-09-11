package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/bogusdeck/ats-scanner/internal/ats"
	"github.com/bogusdeck/ats-scanner/internal/models"
	"github.com/bogusdeck/ats-scanner/internal/parser"
)

func main() {
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/api/platforms", handlePlatforms)
	r.POST("/api/analyze", handleAnalyze)
	r.POST("/api/parse", handleParse)
	r.POST("/api/upload", handleUpload)
	r.POST("/api/fix-resume/generate", handleFixResumeGenerate)
	r.POST("/api/fix-resume/compile", handleFixResumeCompile)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	fmt.Printf("🚀 CV-RADAR API running on http://localhost:%s\n", port)
	r.Run(":" + port)
}

// GET /api/platforms — list all supported ATS platforms
func handlePlatforms(c *gin.Context) {
	type Platform struct {
		Name     string `json:"name"`
		Vendor   string `json:"vendor"`
		Strategy string `json:"strategy"`
		Behavior string `json:"behavior"`
	}
	platforms := []Platform{
		{"Workday", "Workday", "Exact + HiredScore AI", "Strict parser, skips headers/footers, penalizes creative formats"},
		{"Taleo", "Oracle", "Literal exact match", "Strictest keyword matching, auto-reject via Req Rank"},
		{"iCIMS", "iCIMS", "Semantic (ML-based)", "Role Fit AI, grammar-based NLP parser, most forgiving"},
		{"Greenhouse", "Greenhouse", "Semantic (LLM-based)", "No auto-scoring by design, human review with scorecards"},
		{"Lever", "Employ", "Stemming-based", "No ranking, search-dependent, abbreviation-blind"},
		{"SuccessFactors", "SAP", "Taxonomy normalization", "Textkernel parser, Joule AI skills matching"},
	}
	c.JSON(http.StatusOK, platforms)
}

// POST /api/analyze — analyze CV against JD
func handleAnalyze(c *gin.Context) {
	var req models.AnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	if req.CVText == "" || req.JDText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cv_text and jd_text are required"})
		return
	}

	// Auto-detect and strip LaTeX if needed
	cvText := req.CVText
	latexDocClassRegex := regexp.MustCompile(`(?m)^\s*\\(documentclass|begin\{document\})`)
	if latexDocClassRegex.MatchString(cvText) {
		cvText = parser.StripLatex(cvText)
	}
	cv := parser.ParseText(cvText)

	var results []models.ATSResult
	if req.Platform != "" {
		engine := ats.Get(req.Platform)
		if engine == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown platform: " + req.Platform})
			return
		}
		results = []models.ATSResult{engine.Analyze(cv, req.JDText)}
	} else {
		for _, engine := range ats.All() {
			results = append(results, engine.Analyze(cv, req.JDText))
		}
	}

	c.JSON(http.StatusOK, models.AnalysisResponse{
		ParsedCV: cv,
		Results:  results,
	})
}

// POST /api/parse — parse CV to structured JSON (no JD needed)
func handleParse(c *gin.Context) {
	var req struct {
		CVText string `json:"cv_text"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.CVText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cv_text is required"})
		return
	}
	cv := parser.ParseText(req.CVText)
	c.JSON(http.StatusOK, cv)
}

// POST /api/upload — upload a PDF or .tex file, returns extracted text
func handleUpload(c *gin.Context) {
	// Limit upload size to 10 MB
	const maxUploadSize = 10 << 20 // 10 MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload failed: " + err.Error()})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	tmp, err := os.CreateTemp("", "cv-*"+ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create temp file"})
		return
	}
	defer os.Remove(tmp.Name())

	if _, err = io.Copy(tmp, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save file"})
		return
	}
	tmp.Close()

	var text string
	switch ext {
	case ".pdf":
		text, err = parser.ExtractPDF(tmp.Name())
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "PDF parse failed: " + err.Error()})
			return
		}
	case ".tex":
		raw, err := os.ReadFile(tmp.Name())
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "could not read tex file"})
			return
		}
		text = parser.StripLatex(string(raw))
	default:
		raw, err := os.ReadFile(tmp.Name())
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "could not read file"})
			return
		}
		text = string(raw)
	}

	c.JSON(http.StatusOK, gin.H{"text": text})
}
