package ats

import (
	"regexp"
	"sort"
	"strings"
)

// shared NLP helpers used across all ATS engines

func toLower(s string) string { return strings.ToLower(s) }

// tokenize splits text into lowercase words, stripping punctuation
func tokenize(text string) []string {
	re := regexp.MustCompile(`[a-zA-Z0-9#+.\-]+`)
	tokens := re.FindAllString(strings.ToLower(text), -1)
	return tokens
}

// tokenSet returns a map of unique tokens
func tokenSet(text string) map[string]bool {
	set := map[string]bool{}
	for _, t := range tokenize(text) {
		set[t] = true
	}
	return set
}

// extractKeywords pulls meaningful multi-word and single-word keywords from JD text
func extractKeywords(jd string) []string {
	// Common JD noise words to skip
	stop := map[string]bool{
		"and": true, "or": true, "the": true, "a": true, "an": true, "in": true,
		"to": true, "of": true, "for": true, "with": true, "is": true, "are": true,
		"you": true, "will": true, "we": true, "our": true, "your": true, "that": true,
		"this": true, "be": true, "have": true, "on": true, "at": true, "by": true,
		"as": true, "can": true, "must": true, "should": true, "may": true,
	}

	// First grab known tech phrases (multi-word)
	techPhrases := []string{
		"machine learning", "deep learning", "natural language processing", "computer vision",
		"react native", "node.js", "next.js", "rest api", "graphql", "ci/cd",
		"django rest framework", "llm", "large language model", "vector search",
		"prompt engineering", "retrieval augmented generation", "rag pipeline",
		"openai api", "aws lambda", "sql queries", "raw sql", "async queues",
		"tailwind css", "react.js", "full stack", "full-stack", "microservices",
	}

	found := map[string]bool{}
	lower := strings.ToLower(jd)
	for _, phrase := range techPhrases {
		if strings.Contains(lower, phrase) {
			found[phrase] = true
		}
	}

	// Single word keywords
	for _, tok := range tokenize(jd) {
		if len(tok) > 2 && !stop[tok] {
			found[tok] = true
		}
	}

	keys := make([]string, 0, len(found))
	for k := range found {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// exactMatch returns matched and missing keyword lists (pure exact matching)
func exactMatch(cvText string, keywords []string) (matched, missing []string) {
	lower := strings.ToLower(cvText)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			matched = append(matched, kw)
		} else {
			missing = append(missing, kw)
		}
	}
	return
}

// stemWord applies a very basic Porter-like stemmer
func stemWord(word string) string {
	suffixes := []string{"ing", "tion", "ed", "er", "ly", "ment", "ness", "s", "es"}
	w := strings.ToLower(word)
	for _, suffix := range suffixes {
		if strings.HasSuffix(w, suffix) && len(w)-len(suffix) > 3 {
			return w[:len(w)-len(suffix)]
		}
	}
	return w
}

// stemmedMatch applies stemming before matching
func stemmedMatch(cvText string, keywords []string) (matched, missing []string) {
	cvTokens := tokenize(cvText)
	cvStems := map[string]bool{}
	for _, t := range cvTokens {
		cvStems[stemWord(t)] = true
	}

	for _, kw := range keywords {
		kwStem := stemWord(kw)
		if cvStems[kwStem] {
			matched = append(matched, kw)
		} else {
			missing = append(missing, kw)
		}
	}
	return
}

// tfidfSimilarity approximates cosine similarity using TF weighting
func tfidfSimilarity(cvText, jdText string) float64 {
	cvTokens := tokenize(cvText)
	jdTokens := tokenize(jdText)

	if len(cvTokens) == 0 || len(jdTokens) == 0 {
		return 0
	}

	// Build TF for JD
	jdTF := map[string]float64{}
	for _, t := range jdTokens {
		jdTF[t]++
	}

	// Count how many JD tokens appear in CV
	cvSet := tokenSet(cvText)
	dotProduct := 0.0
	jdMag := 0.0
	for t, tf := range jdTF {
		jdMag += tf * tf
		if cvSet[t] {
			dotProduct += tf
		}
	}
	if jdMag == 0 {
		return 0
	}
	return dotProduct / jdMag
}

// scoreToPercent clamps and rounds a ratio to 0–100
func scoreToPercent(ratio float64) float64 {
	if ratio > 1 {
		ratio = 1
	}
	if ratio < 0 {
		ratio = 0
	}
	return ratio * 100
}

// skillsOverlap measures what fraction of JD skills appear in CV skills
func skillsOverlap(cvSkills []string, jdText string) float64 {
	if len(cvSkills) == 0 {
		return 0
	}
	jdLower := strings.ToLower(jdText)
	matched := 0
	for _, s := range cvSkills {
		if strings.Contains(jdLower, strings.ToLower(s)) {
			matched++
		}
	}
	return float64(matched) / float64(len(cvSkills))
}

// detectFormatIssues looks for signs of complex formatting
func detectFormatIssues(rawText string) []string {
	var issues []string
	if strings.Count(rawText, "|") > 10 {
		issues = append(issues, "Tables or columns detected — may confuse strict parsers")
	}
	if strings.Count(rawText, "\t") > 5 {
		issues = append(issues, "Tab characters detected — prefer space-based formatting")
	}
	if strings.Contains(rawText, "\\begin{tabular") || strings.Contains(rawText, "\\begin{table") {
		issues = append(issues, "LaTeX table environments found — convert to plain lists")
	}
	return issues
}
