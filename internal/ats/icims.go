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
	matched, missing, kwScore := weightedStemmedMatch(cv.RawText, keywords)

	semanticBonus := tfidfSimilarity(cv.RawText, jd) * 40
	if semanticBonus > 40 {
		semanticBonus = 40
	}

	structureScore := calcStructureScore(cv)
	skillScore := skillsOverlap(cv.Skills, jd) * 100
	expScore := calcExpScore(cv.YearsExp, jd)
	eduScore := calcEduScore(cv)

	base := kwScore*0.30 + skillScore*0.20 + structureScore*0.20 + expScore*0.20 + eduScore*0.10
	total := base + semanticBonus*0.3
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

	return withInterval(models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		Confidence:      40,
		SimulationNote:  "iCIMS Role Fit is RELATIVE, not absolute — it ranks you within the applicant pool for that specific job. The same CV can be Tier 1 for one job and Tier 3 for another depending on who else applied. No percentage output exists. Our score is an approximation of keyword/skill alignment only.",
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
	})
}

// ─── Greenhouse ──────────────────────────────────────────────────────────────

// GreenhouseEngine — Greenhouse (human-first)
// LLM-style broad semantic matching; no auto-scoring; scorecards filled by humans
// Confidence: 50% — Greenhouse intentionally has no auto-scoring algorithm
type GreenhouseEngine struct{}

func (e *GreenhouseEngine) Name() string   { return "Greenhouse" }
func (e *GreenhouseEngine) Vendor() string { return "Greenhouse" }

func (e *GreenhouseEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractWeightedKeywords(jd)
	matched, missing, kwScore := weightedStemmedMatch(cv.RawText, keywords)

	semanticScore := tfidfSimilarity(cv.RawText, jd) * 150
	if semanticScore > 100 {
		semanticScore = 100
	}

	structureScore := calcStructureScore(cv)
	expScore := calcExpScore(cv.YearsExp, jd)

	narrativeBonus := 0.0
	for _, job := range cv.Experience {
		if len(job.Description) >= 3 {
			narrativeBonus += 5.0
		}
	}
	if narrativeBonus > 20 {
		narrativeBonus = 20
	}

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

	return withInterval(models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		Confidence:      30,
		SimulationNote:  "Greenhouse has ZERO automated resume scoring — confirmed by Greenhouse documentation and recruiters. Scorecards are filled by human interviewers rating candidates as 'Definitely Not / No / Yes / Strong Yes'. Auto-rejection ONLY happens via knockout questions. This score simulates how searchable your CV is to a Greenhouse recruiter, NOT any automated ranking.",
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
	})
}

// ─── Lever ───────────────────────────────────────────────────────────────────

// LeverEngine — Lever (by Employ)
// Stemming-based matching; abbreviation-blind; no auto-ranking (search-dependent)
// Confidence: 65%
type LeverEngine struct{}

func (e *LeverEngine) Name() string   { return "Lever" }
func (e *LeverEngine) Vendor() string { return "Employ" }

func (e *LeverEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractWeightedKeywords(jd)
	matched, missing, kwScore := weightedStemmedMatch(cv.RawText, keywords)

	abbrevWarnings := checkAbbreviations(cv.RawText, jd)
	abbrevPenalty := float64(len(abbrevWarnings)) * 3.5
	kwScore = kwScore - abbrevPenalty
	if kwScore < 0 {
		kwScore = 0
	}

	titleScore := calcTitleScore(cv, jd)
	expScore := calcExpScore(cv.YearsExp, jd)

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

	return withInterval(models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		Confidence:      55,
		SimulationNote:  "Lever is a CRM-style ATS with no auto-scoring. Word stemming is officially confirmed (manage = managing = management). Abbreviation blindness is confirmed (SEO ≠ Search Engine Optimization). Boolean + fuzzy search supported. Score reflects how likely a recruiter keyword search would surface your CV.",
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
	})
}

// ─── SuccessFactors ──────────────────────────────────────────────────────────

// SuccessFactorsEngine — SAP SuccessFactors + Joule AI + Textkernel
// Skills taxonomy normalization; Joule AI infers skills from descriptions
// Confidence: 68%
type SuccessFactorsEngine struct{}

func (e *SuccessFactorsEngine) Name() string   { return "SuccessFactors" }
func (e *SuccessFactorsEngine) Vendor() string { return "SAP" }

func (e *SuccessFactorsEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	normalizedCV := normalizeTaxonomy(cv.RawText)
	normalizedJD := normalizeTaxonomy(jd)

	keywords := extractWeightedKeywords(normalizedJD)
	matched, missing, kwScore := weightedExactMatch(normalizedCV, keywords)

	expScore := calcExpScore(cv.YearsExp, normalizedJD)
	eduScore := calcEduScore(cv)

	jouleBonus := inferSkillsFromDesc(cv, jd) * 0.3
	if jouleBonus > 25 {
		jouleBonus = 25
	}

	base := kwScore*0.45 + expScore*0.30 + eduScore*0.15
	total := base + jouleBonus*0.10
	if total > 100 {
		total = 100
	}

	var recs []string
	recs = append(recs, "SuccessFactors normalizes synonyms — 'Python', 'Python3', and 'py' are treated as equivalent")
	recs = append(recs, "Include a dedicated Technical Skills section; Joule AI specifically parses it for taxonomy matching")
	recs = append(recs, buildKeywordRecs(missing, keywords, "taxonomy-normalized")...)

	return withInterval(models.ATSResult{
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
	})
}

// ─── Shared helpers ──────────────────────────────────────────────────────────

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
		"js":       "javascript", "nodejs": "node.js",
		"ts":       "typescript",
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
