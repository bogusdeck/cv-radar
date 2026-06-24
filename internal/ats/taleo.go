package ats

import (
	"github.com/bogusdeck/ats-scanner/internal/models"
)

// TaleoEngine simulates Oracle Taleo behavior:
// - Strictest: literal exact string matching only
// - Auto-reject (Req Rank) if below ~40% threshold
// - Weight: Exact-keywords(60%) Title(25%) Years(15%)
type TaleoEngine struct{}

func (e *TaleoEngine) Name() string   { return "Taleo" }
func (e *TaleoEngine) Vendor() string { return "Oracle" }

func (e *TaleoEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractKeywords(jd)
	matched, missing := exactMatch(cv.RawText, keywords)

	kwRatio := 0.0
	if len(keywords) > 0 {
		kwRatio = float64(len(matched)) / float64(len(keywords))
	}
	kwScore := scoreToPercent(kwRatio)
	titleScore := calcTitleScore(cv, jd)
	expScore := calcExpScore(cv.YearsExp, jd)

	total := kwScore*0.60 + titleScore*0.25 + expScore*0.15
	autoReject := total < 40

	var warnings []string
	if autoReject {
		warnings = append(warnings, "Score is below Taleo Req Rank threshold (~40%) — likely auto-rejected before human review")
	}

	var recs []string
	recs = append(recs, buildKeywordRecs(missing, "exact literal")...)
	recs = append(recs, "Taleo does NOT support synonyms — mirror exact phrasing from the JD")
	recs = append(recs, "Avoid abbreviations: write the full term if the JD uses the full term")

	return models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		KeywordScore:    round(kwScore),
		StructureScore:  80,
		ExperienceScore: round(expScore),
		EducationScore:  0,
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
	}
}
