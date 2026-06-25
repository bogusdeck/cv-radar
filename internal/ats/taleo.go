package ats

import (
	"fmt"
	"strings"

	"github.com/bogusdeck/ats-scanner/internal/models"
)

// TaleoEngine — Oracle Taleo
// Strictest: literal exact string matching ONLY (no stemming, no synonyms)
// Auto-reject (Req Rank) below ~40% threshold
// Abbreviation-blind by default
// Confidence: 80% — Taleo's exact-match behavior is very well-documented
type TaleoEngine struct{}

func (e *TaleoEngine) Name() string   { return "Taleo" }
func (e *TaleoEngine) Vendor() string { return "Oracle" }

func (e *TaleoEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractWeightedKeywords(jd)
	// Taleo: no bonus for synonyms, only exact literal match counts
	matched, missing, kwScore := weightedExactMatch(cv.RawText, keywords)

	// Extra penalty: Taleo checks for full-form vs abbreviation mismatches
	abbrevPenalty := calcAbbrevPenalty(cv.RawText, jd)
	kwScore = kwScore - abbrevPenalty
	if kwScore < 0 {
		kwScore = 0
	}

	titleScore := calcTitleScore(cv, jd)
	expScore := calcExpScore(cv.YearsExp, jd)

	// Taleo Req Rank: weighted sum, strict thresholds
	total := kwScore*0.60 + titleScore*0.25 + expScore*0.15

	// Taleo auto-rejects earlier than other platforms
	autoReject := total < 40

	var warnings []string
	if autoReject {
		warnings = append(warnings, fmt.Sprintf("Score %.0f%% is below Taleo's Req Rank threshold (~40%%) — likely auto-rejected before human review", total))
	}
	if abbrevPenalty > 0 {
		warnings = append(warnings, fmt.Sprintf("Abbreviation mismatch detected (–%.0f%% penalty): Taleo requires exact spelling", abbrevPenalty))
	}

	var recs []string
	recs = append(recs, buildKeywordRecs(missing, keywords, "exact literal")...)
	recs = append(recs, "Taleo does NOT support synonyms — mirror exact phrasing from the JD")
	recs = append(recs, "Spell out full forms: 'JavaScript' not 'JS', 'PostgreSQL' not 'Postgres', 'Application Programming Interface' not 'API' if JD uses full form")

	return withInterval(models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		Confidence:      62,
		SimulationNote:  "Legacy Taleo uses exact-phrase indexing (well-documented). But there is NO universal Req Rank threshold — it is employer-configured. Oracle Recruiting Cloud (modern Taleo) uses ML scoring 0-5 per dimension. Our simulation models legacy Taleo behavior only.",
		KeywordScore:    round(kwScore),
		StructureScore:  78,
		ExperienceScore: round(expScore),
		EducationScore:  round(calcEduScore(cv)),
		MatchedKeywords: matched,
		MissingKeywords: missing,
		Warnings:        warnings,
		Recommendations: recs,
		AutoReject:      autoReject,
		Breakdown: []models.ScoreBreakdown{
			{Category: "Keyword Match (Exact Literal)", Score: kwScore, MaxScore: 100, Weight: 0.60},
			{Category: "Title Alignment", Score: titleScore, MaxScore: 100, Weight: 0.25},
			{Category: "Experience Years", Score: expScore, MaxScore: 100, Weight: 0.15},
		},
	})
}

func calcAbbrevPenalty(cvText, jd string) float64 {
	abbrevPairs := map[string]string{
		"js": "javascript", "ts": "typescript", "py": "python",
		"pg": "postgresql", "k8s": "kubernetes",
	}
	penalty := 0.0
	cvLower := strings.ToLower(cvText)
	jdLower := strings.ToLower(jd)
	for abbrev, full := range abbrevPairs {
		cvHasAbbrev := strings.Contains(cvLower, " "+abbrev+" ")
		jdHasFull := strings.Contains(jdLower, full)
		cvMissesFull := !strings.Contains(cvLower, full)
		if cvHasAbbrev && jdHasFull && cvMissesFull {
			penalty += 4.0
		}
	}
	return penalty
}
