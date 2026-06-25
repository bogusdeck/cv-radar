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
	return re.FindAllString(strings.ToLower(text), -1)
}

// tokenSet returns a map of unique tokens
func tokenSet(text string) map[string]bool {
	set := map[string]bool{}
	for _, t := range tokenize(text) {
		set[t] = true
	}
	return set
}

// ── Weighted JD Keyword ───────────────────────────────────────────────────────

// WeightedKeyword is a keyword from the JD with an importance weight
type WeightedKeyword struct {
	Word   string
	Weight float64 // 1.0 = required, 0.5 = nice-to-have, 0.8 = inferred
}

// extractWeightedKeywords parses a JD and assigns weights based on which section the keyword appears in.
// Required/must-have keywords count more than nice-to-have.
func extractWeightedKeywords(jd string) []WeightedKeyword {
	stop := map[string]bool{
		"and": true, "or": true, "the": true, "a": true, "an": true, "in": true,
		"to": true, "of": true, "for": true, "with": true, "is": true, "are": true,
		"you": true, "will": true, "we": true, "our": true, "your": true, "that": true,
		"this": true, "be": true, "have": true, "on": true, "at": true, "by": true,
		"as": true, "can": true, "must": true, "should": true, "may": true,
		"also": true, "etc": true, "able": true, "work": true, "team": true,
		"role": true, "job": true, "position": true, "company": true, "applicant": true,
		"candidate": true, "years": true, "year": true, "experience": true,
		"strong": true, "good": true, "great": true, "plus": true,
		"proficiency": true, "familiar": true, "knowledge": true, "understanding": true,
	}

	techPhrases := []string{
		"machine learning", "deep learning", "natural language processing", "computer vision",
		"react native", "node.js", "next.js", "rest api", "graphql", "ci/cd",
		"django rest framework", "large language model", "vector search",
		"prompt engineering", "retrieval augmented generation",
		"openai api", "aws lambda", "async queues",
		"tailwind css", "full stack", "full-stack", "microservices",
		"django rest", "object relational", "test driven", "event driven",
		"software development", "software engineer", "backend engineer",
		"frontend engineer", "distributed systems", "system design",
	}

	// Detect JD sections to set keyword weight
	// Sections like "Requirements", "Must have" → weight 1.0
	// Sections like "Nice to have", "Preferred", "Bonus" → weight 0.5
	type section struct {
		start int
		end   int
		w     float64
	}
	lower := strings.ToLower(jd)
	var sections []section

	requiredHeaders := []string{
		"requirements", "required", "must have", "must-have",
		"what you need", "qualifications", "responsibilities",
		"what we need", "core skills", "skills required",
	}
	niceHeaders := []string{
		"nice to have", "nice-to-have", "preferred", "bonus", "ideally",
		"good to have", "additional", "plus", "not required but",
	}

	sectionBoundaryRe := regexp.MustCompile(`(?i)^(requirements?|responsibilities|qualifications?|preferred|nice.to.have|must.have|what you.ll do|bonus|additional skills?)[:\s]*$`)
	lines := strings.Split(jd, "\n")
	lineStarts := make([]int, len(lines))
	pos := 0
	for i, l := range lines {
		lineStarts[i] = pos
		pos += len(l) + 1
	}

	currentWeight := 0.9 // default: assume most things are required
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if sectionBoundaryRe.MatchString(trimmed) {
			lineL := strings.ToLower(trimmed)
			newW := 0.9
			for _, h := range niceHeaders {
				if strings.Contains(lineL, h) {
					newW = 0.45
					break
				}
			}
			for _, h := range requiredHeaders {
				if strings.Contains(lineL, h) {
					newW = 1.0
					break
				}
			}
			if i > 0 {
				sections = append(sections, section{lineStarts[i], len(jd), newW})
			}
			currentWeight = newW
		}
	}
	_ = currentWeight
	_ = sections

	// Helper: get weight for a position in the JD
	getWeight := func(idx int) float64 {
		w := 0.9
		for _, s := range sections {
			if idx >= s.start {
				w = s.w
			}
		}
		return w
	}

	found := map[string]WeightedKeyword{}

	// Multi-word tech phrases first
	for _, phrase := range techPhrases {
		if idx := strings.Index(lower, phrase); idx != -1 {
			found[phrase] = WeightedKeyword{Word: phrase, Weight: getWeight(idx)}
		}
	}

	// Single-word tokens
	wordRe := regexp.MustCompile(`[a-zA-Z0-9#+.\-]+`)
	wordMatches := wordRe.FindAllStringIndex(lower, -1)
	for _, loc := range wordMatches {
		tok := lower[loc[0]:loc[1]]
		if len(tok) <= 2 || stop[tok] {
			continue
		}
		// Skip words already covered by a phrase
		alreadyCovered := false
		for phrase := range found {
			if strings.Contains(phrase, tok) {
				alreadyCovered = true
				break
			}
		}
		if alreadyCovered {
			continue
		}
		w := getWeight(loc[0])
		if existing, ok := found[tok]; ok {
			if w > existing.Weight {
				found[tok] = WeightedKeyword{Word: tok, Weight: w}
			}
		} else {
			found[tok] = WeightedKeyword{Word: tok, Weight: w}
		}
	}

	result := make([]WeightedKeyword, 0, len(found))
	for _, kw := range found {
		result = append(result, kw)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Weight != result[j].Weight {
			return result[i].Weight > result[j].Weight
		}
		return result[i].Word < result[j].Word
	})
	return result
}

// Legacy unweighted helper — used where weighting isn't needed
func extractKeywords(jd string) []string {
	wkws := extractWeightedKeywords(jd)
	out := make([]string, len(wkws))
	for i, w := range wkws {
		out[i] = w.Word
	}
	return out
}

// ── Matching ─────────────────────────────────────────────────────────────────

// weightedExactMatch returns matched/missing lists and a weighted score (0–100)
// Keywords from required sections score more than nice-to-have ones.
func weightedExactMatch(cvText string, keywords []WeightedKeyword) (matched, missing []string, weightedScore float64) {
	lower := strings.ToLower(cvText)
	totalWeight := 0.0
	matchedWeight := 0.0

	for _, kw := range keywords {
		totalWeight += kw.Weight
		if strings.Contains(lower, strings.ToLower(kw.Word)) {
			matched = append(matched, kw.Word)
			matchedWeight += kw.Weight
		} else {
			missing = append(missing, kw.Word)
		}
	}

	if totalWeight == 0 {
		return matched, missing, 0
	}
	return matched, missing, (matchedWeight / totalWeight) * 100
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

// stemWord applies a basic Porter-like stemmer
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
		if cvStems[stemWord(kw)] {
			matched = append(matched, kw)
		} else {
			missing = append(missing, kw)
		}
	}
	return
}

// weightedStemmedMatch applies stemming + section weighting
func weightedStemmedMatch(cvText string, keywords []WeightedKeyword) (matched, missing []string, weightedScore float64) {
	cvTokens := tokenize(cvText)
	cvStems := map[string]bool{}
	for _, t := range cvTokens {
		cvStems[stemWord(t)] = true
	}

	totalWeight := 0.0
	matchedWeight := 0.0
	for _, kw := range keywords {
		totalWeight += kw.Weight
		if cvStems[stemWord(kw.Word)] {
			matched = append(matched, kw.Word)
			matchedWeight += kw.Weight
		} else {
			missing = append(missing, kw.Word)
		}
	}
	if totalWeight == 0 {
		return matched, missing, 0
	}
	return matched, missing, (matchedWeight / totalWeight) * 100
}

// tfidfSimilarity approximates cosine similarity using TF weighting
func tfidfSimilarity(cvText, jdText string) float64 {
	cvTokens := tokenize(cvText)
	jdTokens := tokenize(jdText)

	if len(cvTokens) == 0 || len(jdTokens) == 0 {
		return 0
	}

	jdTF := map[string]float64{}
	for _, t := range jdTokens {
		jdTF[t]++
	}

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

// ── Scoring Helpers ───────────────────────────────────────────────────────────

func scoreToPercent(ratio float64) float64 {
	if ratio > 1 {
		ratio = 1
	}
	if ratio < 0 {
		ratio = 0
	}
	return ratio * 100
}

// sectionAwareSkillsOverlap weights skills found in the Skills section higher
func sectionAwareSkillsOverlap(cv interface{ GetSkillsInSection() []string; GetAllSkills() []string }, jdText string) float64 {
	// simplified: use flat list but document the intent
	return 0
}

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

func detectFormatIssues(rawText string) []string {
	var issues []string
	if strings.Count(rawText, "|") > 10 {
		issues = append(issues, "Tables or columns detected — may confuse strict parsers")
	}
	if strings.Count(rawText, "\t") > 5 {
		issues = append(issues, "Tab characters detected — prefer space-based formatting")
	}
	if strings.Contains(rawText, `\begin{tabular`) || strings.Contains(rawText, `\begin{table`) {
		issues = append(issues, "LaTeX table environments found — convert to plain lists")
	}
	return issues
}
