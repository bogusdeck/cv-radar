package ats

import (
	"fmt"
	"strings"

	"github.com/bogusdeck/ats-scanner/internal/models"
)

// ICIMSEngine simulates iCIMS behavior:
// - Semantic/ML-based via TF-IDF cosine similarity approximation
// - Most forgiving, grammar-based NLP parser
// - Weight: Semantic(50%) Grammar/Structure(20%) Format(30%)
type ICIMSEngine struct{}

func (e *ICIMSEngine) Name() string   { return "iCIMS" }
func (e *ICIMSEngine) Vendor() string { return "iCIMS" }

func (e *ICIMSEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractKeywords(jd)

	// iCIMS uses semantic matching — use TF-IDF similarity
	semanticScore := tfidfSimilarity(cv.RawText, jd) * 100
	if semanticScore > 100 {
		semanticScore = 100
	}

	// Also do a soft keyword match (matched = bonus)
	matched, missing := exactMatch(cv.RawText, keywords)
	kwScore := scoreToPercent(float64(len(matched)) / max1(float64(len(keywords))))

	// Grammar/Structure score: iCIMS rewards well-structured CVs
	structureScore := calcStructureScore(cv)

	// Skills overlap
	skillScore := skillsOverlap(cv.Skills, jd) * 100

	total := semanticScore*0.50 + structureScore*0.20 + skillScore*0.30
	if total > 100 {
		total = 100
	}

	var warnings []string
	if len(cv.Skills) < 5 {
		warnings = append(warnings, "Skills section appears thin — iCIMS rewards rich skills sections")
	}

	var recs []string
	recs = append(recs, "iCIMS is forgiving with synonyms — focus on natural language describing your work")
	recs = append(recs, buildKeywordRecs(missing, "semantic")...)
	if structureScore < 60 {
		recs = append(recs, "Improve CV structure: use clear section headers (Summary, Experience, Skills, Education)")
	}

	return models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		KeywordScore:    round(kwScore),
		StructureScore:  round(structureScore),
		ExperienceScore: round(skillScore),
		EducationScore:  round(calcEduScore(cv)),
		MatchedKeywords: matched,
		MissingKeywords: missing,
		Warnings:        warnings,
		Recommendations: recs,
		AutoReject:      false, // iCIMS doesn't auto-reject
		Breakdown: []models.ScoreBreakdown{
			{Category: "Semantic Similarity (TF-IDF)", Score: semanticScore, MaxScore: 100, Weight: 0.50},
			{Category: "Grammar & Structure", Score: structureScore, MaxScore: 100, Weight: 0.20},
			{Category: "Skills Coverage", Score: skillScore, MaxScore: 100, Weight: 0.30},
		},
	}
}

// ─── Greenhouse ────────────────────────────────────────────────────────────

// GreenhouseEngine simulates Greenhouse behavior:
// - LLM-style broad semantic matching
// - Scorecard simulation (rates each JD requirement)
// - No auto-rejection, encourages human review
// - Weight: Semantic(70%) Structure(30%)
type GreenhouseEngine struct{}

func (e *GreenhouseEngine) Name() string   { return "Greenhouse" }
func (e *GreenhouseEngine) Vendor() string { return "Greenhouse" }

func (e *GreenhouseEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractKeywords(jd)
	matched, missing := exactMatch(cv.RawText, keywords)

	// Greenhouse: broad semantic — use high-tolerance similarity
	semanticScore := tfidfSimilarity(cv.RawText, jd) * 130 // boosted for lenient matching
	if semanticScore > 100 {
		semanticScore = 100
	}

	structureScore := calcStructureScore(cv)
	expScore := calcExpScore(cv.YearsExp, jd)

	total := semanticScore*0.70 + structureScore*0.30

	var recs []string
	recs = append(recs, "Greenhouse emphasizes storytelling — use strong action verbs and quantified achievements")
	recs = append(recs, "Scorecards are filled by humans — make your experience bullet points compelling narratives")
	recs = append(recs, buildKeywordRecs(missing, "semantic")...)

	return models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		KeywordScore:    round(scoreToPercent(float64(len(matched)) / max1(float64(len(keywords))))),
		StructureScore:  round(structureScore),
		ExperienceScore: round(expScore),
		EducationScore:  round(calcEduScore(cv)),
		MatchedKeywords: matched,
		MissingKeywords: missing,
		Warnings:        nil,
		Recommendations: recs,
		AutoReject:      false,
		Breakdown: []models.ScoreBreakdown{
			{Category: "Broad Semantic Match", Score: semanticScore, MaxScore: 100, Weight: 0.70},
			{Category: "CV Structure & Clarity", Score: structureScore, MaxScore: 100, Weight: 0.30},
		},
	}
}

// ─── Lever ─────────────────────────────────────────────────────────────────

// LeverEngine simulates Lever behavior:
// - Stemming-based matching
// - Abbreviation-blind (JS ≠ JavaScript)
// - No auto-ranking, search-dependent
// - Weight: Stemmed-match(65%) Title(35%)
type LeverEngine struct{}

func (e *LeverEngine) Name() string   { return "Lever" }
func (e *LeverEngine) Vendor() string { return "Employ" }

func (e *LeverEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractKeywords(jd)
	matched, missing := stemmedMatch(cv.RawText, keywords)

	kwScore := scoreToPercent(float64(len(matched)) / max1(float64(len(keywords))))
	titleScore := calcTitleScore(cv, jd)

	total := kwScore*0.65 + titleScore*0.35

	// Abbreviation warnings
	abbrevWarnings := checkAbbreviations(cv.RawText, jd)
	var recs []string
	recs = append(recs, "Lever is abbreviation-blind — spell out 'JavaScript' not 'JS', 'PostgreSQL' not 'PG'")
	recs = append(recs, buildKeywordRecs(missing, "stemmed")...)
	if len(abbrevWarnings) > 0 {
		recs = append(recs, fmt.Sprintf("Detected possible abbreviations that Lever may miss: %s", strings.Join(abbrevWarnings, ", ")))
	}

	return models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		KeywordScore:    round(kwScore),
		StructureScore:  75, // Lever doesn't penalize structure much
		ExperienceScore: round(calcExpScore(cv.YearsExp, jd)),
		EducationScore:  round(calcEduScore(cv)),
		MatchedKeywords: matched,
		MissingKeywords: missing,
		Warnings:        abbrevWarnings,
		Recommendations: recs,
		AutoReject:      false,
		Breakdown: []models.ScoreBreakdown{
			{Category: "Stemmed Keyword Match", Score: kwScore, MaxScore: 100, Weight: 0.65},
			{Category: "Title Alignment", Score: titleScore, MaxScore: 100, Weight: 0.35},
		},
	}
}

// ─── SuccessFactors ────────────────────────────────────────────────────────

// SuccessFactorsEngine simulates SAP SuccessFactors + Joule AI behavior:
// - Skills taxonomy normalization (Python ≈ Python3)
// - Joule AI infers skills from experience descriptions
// - Weight: Skills-taxonomy(50%) Experience(30%) Education(20%)
type SuccessFactorsEngine struct{}

func (e *SuccessFactorsEngine) Name() string   { return "SuccessFactors" }
func (e *SuccessFactorsEngine) Vendor() string { return "SAP" }

func (e *SuccessFactorsEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractKeywords(jd)

	// Taxonomy normalization — expand synonyms before matching
	normalizedCV := normalizeTaxonomy(cv.RawText)
	normalizedJD := normalizeTaxonomy(jd)

	normKeywords := extractKeywords(normalizedJD)
	matched, missing := exactMatch(normalizedCV, normKeywords)

	kwScore := scoreToPercent(float64(len(matched)) / max1(float64(len(normKeywords))))
	expScore := calcExpScore(cv.YearsExp, jd)
	eduScore := calcEduScore(cv)

	// Joule AI: infer additional skills from descriptions
	inferredScore := inferSkillsFromDesc(cv, jd)

	total := kwScore*0.50 + expScore*0.30 + eduScore*0.20
	total = (total + inferredScore) / 2.0 // blend with inferred

	var recs []string
	recs = append(recs, "SuccessFactors normalizes skill synonyms — 'Python', 'Python3', 'py' are treated as equivalent")
	recs = append(recs, "Include a dedicated Skills section; Joule AI specifically parses it")
	recs = append(recs, buildKeywordRecs(missing, "taxonomy-normalized")...)

	_ = keywords // original for reference

	return models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		KeywordScore:    round(kwScore),
		StructureScore:  round(inferredScore),
		ExperienceScore: round(expScore),
		EducationScore:  round(eduScore),
		MatchedKeywords: matched,
		MissingKeywords: missing,
		Warnings:        nil,
		Recommendations: recs,
		AutoReject:      false,
		Breakdown: []models.ScoreBreakdown{
			{Category: "Skills Taxonomy Match", Score: kwScore, MaxScore: 100, Weight: 0.50},
			{Category: "Experience Match", Score: expScore, MaxScore: 100, Weight: 0.30},
			{Category: "Education Match", Score: eduScore, MaxScore: 100, Weight: 0.20},
		},
	}
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func calcStructureScore(cv models.ParsedCV) float64 {
	score := 0.0
	if cv.Summary != "" {
		score += 20
	}
	if len(cv.Experience) > 0 {
		score += 30
	}
	if len(cv.Skills) > 3 {
		score += 25
	}
	if len(cv.Education) > 0 {
		score += 15
	}
	if len(cv.Projects) > 0 {
		score += 10
	}
	return score
}

func max1(f float64) float64 {
	if f < 1 {
		return 1
	}
	return f
}

// checkAbbreviations warns about common abbreviations Lever may miss
func checkAbbreviations(cvText, jd string) []string {
	abbrevMap := map[string]string{
		"js": "javascript", "ts": "typescript", "py": "python",
		"pg": "postgresql", "k8s": "kubernetes", "ml": "machine learning",
		"ai": "artificial intelligence", "api": "application programming interface",
		"db": "database", "fe": "frontend", "be": "backend",
	}
	var found []string
	cvLower := strings.ToLower(cvText)
	for abbrev, full := range abbrevMap {
		if strings.Contains(cvLower, abbrev) && strings.Contains(strings.ToLower(jd), full) {
			found = append(found, fmt.Sprintf("%s→%s", abbrev, full))
		}
	}
	return found
}

// normalizeTaxonomy replaces skill synonyms with canonical forms
func normalizeTaxonomy(text string) string {
	synonyms := map[string]string{
		"python3": "python", "py":          "python",
		"js":      "javascript", "nodejs":  "node.js",
		"ts":      "typescript",
		"postgres": "postgresql", "pg":     "postgresql",
		"k8s":    "kubernetes", "kube":     "kubernetes",
		"ml":     "machine learning", "ai": "artificial intelligence",
		"llm":    "large language model",
		"golang": "go",
		"react native": "react",
	}
	lower := strings.ToLower(text)
	for syn, canonical := range synonyms {
		lower = strings.ReplaceAll(lower, syn, canonical)
	}
	return lower
}

// inferSkillsFromDesc gives bonus score when skills are evident from description even if not listed
func inferSkillsFromDesc(cv models.ParsedCV, jd string) float64 {
	jdLower := strings.ToLower(jd)
	inferBonus := 0.0
	totalChecks := 0.0

	for _, job := range cv.Experience {
		for _, desc := range job.Description {
			descLower := strings.ToLower(desc)
			// If JD mentions a tech and the description mentions related concepts, give credit
			if strings.Contains(jdLower, "api") && (strings.Contains(descLower, "endpoint") || strings.Contains(descLower, "rest")) {
				inferBonus++
			}
			if strings.Contains(jdLower, "database") && (strings.Contains(descLower, "query") || strings.Contains(descLower, "sql")) {
				inferBonus++
			}
			if strings.Contains(jdLower, "cloud") && (strings.Contains(descLower, "aws") || strings.Contains(descLower, "deploy")) {
				inferBonus++
			}
			totalChecks += 3
		}
	}

	if totalChecks == 0 {
		return 50
	}
	return scoreToPercent(inferBonus / totalChecks)
}
