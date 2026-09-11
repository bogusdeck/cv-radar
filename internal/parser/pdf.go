package parser

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ExtractPDF extracts plain text from a PDF file path using the Python PyMuPDF script
func ExtractPDF(path string) (string, error) {
	// The Python script is located in the root of the project.
	// Since the Go server is run from the project root (via start.sh), we can just call it.
	
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("python3", "pdf_extractor.py", absPath)
	
	// fallback for windows
	if runtime.GOOS == "windows" {
		cmd = exec.Command("python", "pdf_extractor.py", absPath)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run python PDF extractor: %w\nStderr: %s", err, stderr.String())
	}

	return out.String(), nil
}
