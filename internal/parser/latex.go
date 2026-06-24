package parser

import (
	"regexp"
	"strings"
)

// StripLatex converts a LaTeX .tex CV file content to plain text
func StripLatex(src string) string {
	// ── Step 1: discard the preamble (everything before \begin{document}) ──
	if idx := strings.Index(src, `\begin{document}`); idx != -1 {
		src = src[idx+len(`\begin{document}`):]
	}
	// Remove \end{document}
	src = strings.ReplaceAll(src, `\end{document}`, "")

	// ── Step 2: remove comments ──
	commentRe := regexp.MustCompile(`(?m)%[^\n]*$`)
	src = commentRe.ReplaceAllString(src, "")

	// ── Step 3: replace \href{url}{text} → text ──
	hrefRe := regexp.MustCompile(`\\href\{[^}]*\}\{([^}]*)\}`)
	src = hrefRe.ReplaceAllString(src, "$1")

	// ── Step 4: replace \section{X} → \nX\n ──
	sectionRe := regexp.MustCompile(`\\(?:section|subsection)\{([^}]*)\}`)
	src = sectionRe.ReplaceAllString(src, "\n$1\n")

	// ── Step 5: keep content of formatting commands ──
	fmtRe := regexp.MustCompile(`\\(?:textbf|textit|emph|underline|small|large|Huge|huge|scshape|textsc|textrm|textsf|texttt)\{([^}]*)\}`)
	for fmtRe.MatchString(src) {
		src = fmtRe.ReplaceAllString(src, "$1")
	}

	// ── Step 6: resumeItem{…} → bullet ──
	riRe := regexp.MustCompile(`\\resumeItem\{`)
	src = riRe.ReplaceAllString(src, "• ")

	// ── Step 7: strip all remaining \command{...} (keep content) ──
	argCmd := regexp.MustCompile(`\\[a-zA-Z]+\*?\{([^{}]*)\}`)
	for argCmd.MatchString(src) {
		src = argCmd.ReplaceAllString(src, " $1 ")
	}

	// ── Step 8: strip remaining \command without args ──
	plainCmd := regexp.MustCompile(`\\[a-zA-Z@]+\*?`)
	src = plainCmd.ReplaceAllString(src, " ")

	// ── Step 9: clean special characters ──
	src = strings.ReplaceAll(src, "{", " ")
	src = strings.ReplaceAll(src, "}", " ")
	src = strings.ReplaceAll(src, "\\&", "&")
	src = strings.ReplaceAll(src, "\\%", "%")
	src = strings.ReplaceAll(src, "\\$", "$")
	src = strings.ReplaceAll(src, "$|$", " | ")
	src = strings.ReplaceAll(src, "--", "–")
	src = strings.ReplaceAll(src, "\\\\", "\n")
	src = strings.ReplaceAll(src, "$", "")

	// ── Step 10: normalize whitespace ──
	multiSpace := regexp.MustCompile(`[ \t]+`)
	src = multiSpace.ReplaceAllString(src, " ")
	multiNewline := regexp.MustCompile(`\n{3,}`)
	src = multiNewline.ReplaceAllString(src, "\n\n")

	return strings.TrimSpace(src)
}

