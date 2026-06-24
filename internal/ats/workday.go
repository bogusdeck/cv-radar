package ats

import (
	"fmt"
	"strings"

	"github.com/bogusdeck/ats-scanner/internal/models"
)

// WorkdayEngine simulates Workday + HiredScore AI behavior:
// - Exact keyword matching + title matching
// - Penalizes creative/complex formatting
// - Strict section parsing (ignores headers/footers)
// - Weight: Title(30%) Skills-exact(40%) Experience-years(20%) Education(10%)
type WorkdayEngine struct{}

func (e *WorkdayEngine) Name() string   { return "Workday" }
func (e *WorkdayEngine) Vendor() string { return "Workday" }

func (e *WorkdayEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractKeywords(jd)
	matched, missing := exactMatch(cv.RawText, keywords)

	// Keyword score (40% weight)
	kwRatio := 0.0
	if len(keywords) > 0 {
		kwRatio = float64(len(matched)) / float64(len(keywords))
	}
	kwScore := scoreToPercent(kwRatio)

	// Title match (30% weight) — check if CV job titles match JD title keywords
	titleScore := calcTitleScore(cv, jd)

	// Experience years (20% weight)
	expScore := calcExpScore(cv.YearsExp, jd)

	// Education (10% weight)
	eduScore := calcEduScore(cv)

	// Format penalties
	var warnings []string
	formatPenalty := 0.0
	issues := detectFormatIssues(cv.RawText)
	for _, iss := range issues {
		warnings = append(warnings, iss)
		formatPenalty += 5.0
	}

	total := (kwScore*0.40 + titleScore*0.30 + expScore*0.20 + eduScore*0.10) - formatPenalty
	if total < 0 {
		total = 0
	}

	var recs []string
	recs = append(recs, buildKeywordRecs(missing, "exact")...)
	if titleScore < 50 {
		recs = append(recs, "Align your job title more closely with the target role title")
	}
	if formatPenalty > 0 {
		recs = append(recs, "Use a simple single-column ATS-friendly format; avoid tables and special characters")
	}

	return models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		KeywordScore:    round(kwScore),
		StructureScore:  round(100 - formatPenalty),
		ExperienceScore: round(expScore),
		EducationScore:  round(eduScore),
		MatchedKeywords: matched,
		MissingKeywords: missing,
		Warnings:        warnings,
		Recommendations: recs,
		AutoReject:      total < 30,
		Breakdown: []models.ScoreBreakdown{
			{Category: "Keyword Match (Exact)", Score: kwScore, MaxScore: 100, Weight: 0.40},
			{Category: "Title Alignment", Score: titleScore, MaxScore: 100, Weight: 0.30},
			{Category: "Experience Years", Score: expScore, MaxScore: 100, Weight: 0.20},
			{Category: "Education", Score: eduScore, MaxScore: 100, Weight: 0.10},
		},
	}
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func calcTitleScore(cv models.ParsedCV, jd string) float64 {
	jdLower := strings.ToLower(jd)
	titleKeywords := []string{"engineer", "developer", "backend", "frontend", "full stack", "sde", "senior", "lead", "architect", "analyst", "scientist"}
	hits := 0
	for _, kw := range titleKeywords {
		if strings.Contains(jdLower, kw) {
			// Check if it appears in any of the CV job titles
			for _, job := range cv.Experience {
				if strings.Contains(strings.ToLower(job.Title), kw) {
					hits++
					break
				}
			}
		}
	}
	if len(titleKeywords) == 0 {
		return 50
	}
	return scoreToPercent(float64(hits) / float64(len(titleKeywords)))
}

func calcExpScore(yearsExp float64, jd string) float64 {
	// Try to detect required years in JD
	required := 2.0
	jdLower := strings.ToLower(jd)
	if strings.Contains(jdLower, "5+ years") || strings.Contains(jdLower, "5 years") {
		required = 5
	} else if strings.Contains(jdLower, "3+ years") || strings.Contains(jdLower, "3 years") {
		required = 3
	} else if strings.Contains(jdLower, "1+ year") || strings.Contains(jdLower, "1 year") {
		required = 1
	}
	if yearsExp >= required {
		return 100
	}
	return scoreToPercent(yearsExp / required)
}

func calcEduScore(cv models.ParsedCV) float64 {
	if len(cv.Education) == 0 {
		return 50
	}
	for _, edu := range cv.Education {
		lower := strings.ToLower(edu.Degree)
		if strings.Contains(lower, "bachelor") || strings.Contains(lower, "b.tech") || strings.Contains(lower, "b.e") {
			return 85
		}
		if strings.Contains(lower, "master") || strings.Contains(lower, "m.tech") {
			return 100
		}
	}
	return 60
}

func buildKeywordRecs(missing []string, matchType string) []string {
	if len(missing) == 0 {
		return nil
	}
	top := missing
	if len(top) > 8 {
		top = top[:8]
	}
	return []string{fmt.Sprintf("Add these missing keywords (%s): %s", matchType, strings.Join(top, ", "))}
}

func round(f float64) float64 {
	return float64(int(f*10+0.5)) / 10
}
