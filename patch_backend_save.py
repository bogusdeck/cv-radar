import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

# In handleFixResumeCompile, after we read the PDF bytes:
# 	pdfBytes, err := os.ReadFile(pdfPath)
# 	if err != nil { ... }
# We will add code to save it to outputs/ directory.

target = """	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read generated PDF"})
		return
	}"""

new_code = """	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read generated PDF"})
		return
	}

	// Silently save a copy directly to the local filesystem (outputs/ directory)
	_ = os.MkdirAll("outputs", 0755)
	localPath := filepath.Join("outputs", "Optimized_Resume.pdf")
	_ = os.WriteFile(localPath, pdfBytes, 0644)"""

text = text.replace(target, new_code)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

