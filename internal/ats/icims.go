package ats

import (
	"fmt"
	"strings"

	"github.com/bogusdeck/ats-scanner/internal/models"
)

// ICIMSEngine — iCIMS (most forgiving)
// Semantic/ML-based, grammar-aware, synonym-tolerant
// Confidence: 60% — Role Fit AI is proprietary ML, we approximate with TF-IDF
type ICIMSEngine struct{}

func (e *ICIMSEngine) Name() string   { return "iCIMS" }
func (e *ICIMSEngine) Vendor() string { return "iCIMS" }

func (e *ICIMSEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractWeightedKeywords(jd)
	// iCIMS uses semantic + synonym matching — use stemmed weighted match
	matched, missing, kwScore := weightedStemmedMatch(cv.RawText, keywords)

	// Semantic similarity bonus — iCIMS rewards contextual relevance
	semanticBonus := tfidfSimilarity(cv.RawText, jd) * 40 // up to 40pt bonus
	if semanticBonus > 40 {
		semanticBonus = 40
	}

	structureScore := calcStructureScore(cv)
	skillScore := skillsOverlap(cv.Skills, jd) * 100
	expScore := calcExpScore(cv.YearsExp, jd)
	eduScore := calcEduScore(cv)

	// iCIMS: semantic(40%) + keyword(30%) + structure(20%) + education(10%)
	base := kwScore*0.30 + skillScore*0.20 + structureScore*0.20 + expScore*0.20 + eduScore*0.10
	total := base + semanticBonus*0.3 // semantic is a bonus, not dominant
	if total > 100 {
		total = 100
	}

	var warnings []string
	if len(cv.Skills) < 5 {
		warnings = append(warnings, "Skills section appears thin — iCIMS rewards rich, explicit skills sections")
	}

	var recs []string
	recs = append(recs, "iCIMS is forgiving with synonyms — focus on natural language and thorough job descriptions")
	recs = append(recs, buildKeywordRecs(missing, keywords, "semantic/stemmed")...)
	if structureScore < 60 {
		recs = append(recs, "Improve CV structure: use clear section headers (Summary, Experience, Skills, Education)")
	}

	return models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		Confidence:      60,
		SimulationNote:  "iCIMS Role Fit AI uses proprietary ML. We approximate with TF-IDF + stemming. Score may vary ±15 points vs real iCIMS.",
		KeywordScore:    round(kwScore),
		StructureScore:  round(structureScore),
		ExperienceScore: round(expScore),
		EducationScore:  round(eduScore),
		MatchedKeywords: matched,
		MissingKeywords: missing,
		Warnings:        warnings,
		Recommendations: recs,
		AutoReject:      false,
		Breakdown: []models.ScoreBreakdown{
			{Category: "Keyword Match (Stemmed)", Score: kwScore, MaxScore: 100, Weight: 0.30},
			{Category: "Skills Coverage", Score: skillScore, MaxScore: 100, Weight: 0.20},
			{Category: "CV Structure & Grammar", Score: structureScore, MaxScore: 100, Weight: 0.20},
			{Category: "Experience Match", Score: expScore, MaxScore: 100, Weight: 0.20},
			{Category: "Education", Score: eduScore, MaxScore: 100, Weight: 0.10},
		},
	}
}

// ─── Greenhouse ──────────────────────────────────────────────────────────────

// GreenhouseEngine — Greenhouse (human-first)
// LLM-style broad semantic matching; no auto-scoring; scorecards filled by humans
// This is the most generous platform — it rarely auto-rejects
// Confidence: 50% — Greenhouse intentionally has no auto-scoring algorithm
type GreenhouseEngine struct{}

func (e *GreenhouseEngine) Name() string   { return "Greenhouse" }
func (e *GreenhouseEngine) Vendor() string { return "Greenhouse" }

func (e *GreenhouseEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractWeightedKeywords(jd)
	matched, missing, kwScore := weightedStemmedMatch(cv.RawText, keywords)

	// Greenhouse: very broad semantic matching (lenient)
	semanticScore := tfidfSimilarity(cv.RawText, jd) * 150
	if semanticScore > 100 {
		semanticScore = 100
	}

	structureScore := calcStructureScore(cv)
	expScore := calcExpScore(cv.YearsExp, jd)

	// Greenhouse weights narrative quality heavily — a well-structured CV scores higher
	narrativeBonus := 0.0
	for _, job := range cv.Experience {
		if len(job.Description) >= 3 {
			narrativeBonus += 5.0
		}
	}
	if narrativeBonus > 20 {
		narrativeBonus = 20
	}

	// Greenhouse: semantic(50%) + structure(30%) + experience(20%)
	total := semanticScore*0.50 + structureScore*0.30 + expScore*0.20 + narrativeBonus
	if total > 100 {
		total = 100
	}

	var recs []string
	recs = append(recs, "Greenhouse emphasizes storytelling — use strong action verbs and quantified achievements")
	recs = append(recs, "Scorecards are filled by humans reviewing your CV — make bullet points compelling, not just keyword-rich")
	if len(cv.Experience) > 0 && len(cv.Experience[0].Description) < 3 {
		recs = append(recs, "Add more bullet points to each experience role — Greenhouse reviewers look for depth")
	}
	recs = append(recs, buildKeywordRecs(missing, keywords, "semantic")...)

	return models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		Confidence:      50,
		SimulationNote:  "Greenhouse has no automated scoring by design. Our score simulates recruiter scorecard likelihood based on keyword density and structure. Real outcomes depend entirely on human reviewers.",
		KeywordScore:    round(kwScore),
		StructureScore:  round(structureScore),
		ExperienceScore: round(expScore),
		EducationScore:  round(calcEduScore(cv)),
		MatchedKeywords: matched,
		MissingKeywords: missing,
		Warnings:        nil,
		Recommendations: recs,
		AutoReject:      false,
		Breakdown: []models.ScoreBreakdown{
			{Category: "Broad Semantic Match", Score: semanticScore, MaxScore: 100, Weight: 0.50},
			{Category: "CV Structure & Clarity", Score: structureScore, MaxScore: 100, Weight: 0.30},
			{Category: "Experience Depth", Score: expScore, MaxScore: 100, Weight: 0.20},
		},
	}
}

// ─── Lever ───────────────────────────────────────────────────────────────────

// LeverEngine — Lever (by Employ)
// Stemming-based matching; abbreviation-blind; no auto-ranking (search-dependent)
// Confidence: 65% — Lever's search behavior is documented, no auto-scoring algorithm
type LeverEngine struct{}

func (e *LeverEngine) Name() string   { return "Lever" }
func (e *LeverEngine) Vendor() string { return "Employ" }

func (e *LeverEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractWeightedKeywords(jd)
	// Lever: stemming-based, abbreviation-blind
	matched, missing, kwScore := weightedStemmedMatch(cv.RawText, keywords)

	// Abbreviation penalty unique to Lever
	abbrevWarnings := checkAbbreviations(cv.RawText, jd)
	abbrevPenalty := float64(len(abbrevWarnings)) * 3.5
	kwScore = kwScore - abbrevPenalty
	if kwScore < 0 {
		kwScore = 0
	}

	titleScore := calcTitleScore(cv, jd)
	expScore := calcExpScore(cv.YearsExp, jd)

	// Lever: keyword(55%) + title(30%) + exp(15%)
	// Lower than Taleo but higher than semantic platforms
	total := kwScore*0.55 + titleScore*0.30 + expScore*0.15
	if total > 100 {
		total = 100
	}

	var recs []string
	recs = append(recs, "Lever is abbreviation-blind — spell out 'JavaScript' not 'JS', 'PostgreSQL' not 'PG'")
	recs = append(recs, buildKeywordRecs(missing, keywords, "stemmed")...)
	if len(abbrevWarnings) > 0 {
		recs = append(recs, fmt.Sprintf("Detected abbreviations that Lever may miss: %s", strings.Join(abbrevWarnings, ", ")))
	}

	return models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		Confidence:      65,
		SimulationNote:  "Lever has no auto-ranking; scores depend on recruiter search queries. We simulate search relevance based on stemmed keyword density.",
		KeywordScore:    round(kwScore),
		StructureScore:  75,
		ExperienceScore: round(expScore),
		EducationScore:  round(calcEduScore(cv)),
		MatchedKeywords: matched,
		MissingKeywords: missing,
		Warnings:        abbrevWarnings,
		Recommendations: recs,
		AutoReject:      false,
		Breakdown: []models.ScoreBreakdown{
			{Category: "Stemmed Keyword Match", Score: kwScore, MaxScore: 100, Weight: 0.55},
			{Category: "Title Alignment", Score: titleScore, MaxScore: 100, Weight: 0.30},
			{Category: "Experience Years", Score: expScore, MaxScore: 100, Weight: 0.15},
		},
	}
}

// ─── SuccessFactors ──────────────────────────────────────────────────────────

// SuccessFactorsEngine — SAP SuccessFactors + Joule AI + Textkernel
// Skills taxonomy normalization; Joule AI infers skills from descriptions
// Confidence: 68% — Textkernel taxonomy is partially known; Joule AI internals are not
type SuccessFactorsEngine struct{}

func (e *SuccessFactorsEngine) Name() string   { return "SuccessFactors" }
func (e *SuccessFactorsEngine) Vendor() string { return "SAP" }

func (e *SuccessFactorsEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	// Normalize synonyms in both CV and JD before matching
	normalizedCV := normalizeTaxonomy(cv.RawText)
	normalizedJD := normalizeTaxonomy(jd)

	keywords := extractWeightedKeywords(normalizedJD)
	matched, missing, kwScore := weightedExactMatch(normalizedCV, keywords)

	expScore := calcExpScore(cv.YearsExp, normalizedJD)
	eduScore := calcEduScore(cv)

	// Joule AI inferred skills bonus
	jouleBonus := inferSkillsFromDesc(cv, jd) * 0.3 // up to 30pt contribution
	if jouleBonus > 25 {
		jouleBonus = 25
	}

	// SuccessFactors: skills taxonomy(45%) + experience(30%) + education(15%) + joule bonus(10%)
	base := kwScore*0.45 + expScore*0.30 + eduScore*0.15
	total := base + jouleBonus*0.10
	if total > 100 {
		total = 100
	}

	var recs []string
	recs = append(recs, "SuccessFactors normalizes synonyms — 'Python', 'Python3', and 'py' are treated as equivalent")
	recs = append(recs, "Include a dedicated Technical Skills section; Joule AI specifically parses it for taxonomy matching")
	recs = append(recs, buildKeywordRecs(missing, keywords, "taxonomy-normalized")...)

	return models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		Confidence:      68,
		SimulationNote:  "SAP Joule AI internals are proprietary. Textkernel taxonomy normalization is partially documented. Score may vary ±12 points vs real SuccessFactors.",
		KeywordScore:    round(kwScore),
		StructureScore:  round(jouleBonus),
		ExperienceScore: round(expScore),
		EducationScore:  round(eduScore),
		MatchedKeywords: matched,
		MissingKeywords: missing,
		Warnings:        nil,
		Recommendations: recs,
		AutoReject:      false,
		Breakdown: []models.ScoreBreakdown{
			{Category: "Skills Taxonomy Match", Score: kwScore, MaxScore: 100, Weight: 0.45},
			{Category: "Experience Match", Score: expScore, MaxScore: 100, Weight: 0.30},
			{Category: "Education", Score: eduScore, MaxScore: 100, Weight: 0.15},
			{Category: "Joule AI Skill Inference", Score: jouleBonus * 10, MaxScore: 100, Weight: 0.10},
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
		score += 25
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
	// Reward quantified achievements
	for _, job := range cv.Experience {
		for _, desc := range job.Description {
			if strings.Contains(desc, "%") || strings.ContainsAny(desc, "0123456789") {
				score += 2
				break
			}
		}
	}
	if score > 100 {
		score = 100
	}
	return score
}

func max1(f float64) float64 {
	if f < 1 {
		return 1
	}
	return f
}

func checkAbbreviations(cvText, jd string) []string {
	abbrevMap := map[string]string{
		"js": "javascript", "ts": "typescript", "py": "python",
		"pg": "postgresql", "k8s": "kubernetes", "ml": "machine learning",
		"api": "application programming interface",
		"db": "database", "fe": "frontend", "be": "backend",
	}
	var found []string
	cvLower := strings.ToLower(cvText)
	for abbrev, full := range abbrevMap {
		if strings.Contains(cvLower, " "+abbrev+" ") && strings.Contains(strings.ToLower(jd), full) {
			found = append(found, fmt.Sprintf("%s→%s", abbrev, full))
		}
	}
	return found
}

func normalizeTaxonomy(text string) string {
	synonyms := map[string]string{
		"python3": "python", "py": "python",
		"js":      "javascript", "nodejs": "node.js",
		"ts":      "typescript",
		"postgres": "postgresql", "pg": "postgresql",
		"k8s":    "kubernetes", "kube": "kubernetes",
		"ml":     "machine learning",
		"golang": "go",
	}
	lower := strings.ToLower(text)
	for syn, canonical := range synonyms {
		lower = strings.ReplaceAll(lower, " "+syn+" ", " "+canonical+" ")
	}
	return lower
}

func inferSkillsFromDesc(cv models.ParsedCV, jd string) float64 {
	jdLower := strings.ToLower(jd)
	inferBonus := 0.0
	totalChecks := 0.0

	for _, job := range cv.Experience {
		for _, desc := range job.Description {
			descLower := strings.ToLower(desc)
			if strings.Contains(jdLower, "api") && (strings.Contains(descLower, "endpoint") || strings.Contains(descLower, "rest")) {
				inferBonus++
			}
			if strings.Contains(jdLower, "database") && (strings.Contains(descLower, "query") || strings.Contains(descLower, "sql")) {
				inferBonus++
			}
			if strings.Contains(jdLower, "cloud") && (strings.Contains(descLower, "aws") || strings.Contains(descLower, "deploy")) {
				inferBonus++
			}
			if strings.Contains(jdLower, "performance") && (strings.Contains(descLower, "latency") || strings.Contains(descLower, "optimiz")) {
				inferBonus++
			}
			totalChecks += 4
		}
	}

	if totalChecks == 0 {
		return 50
	}
	return scoreToPercent(inferBonus / totalChecks)
}
