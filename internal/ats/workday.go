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
// Confidence: 72% — HiredScore weighting is proprietary but behavior is well-documented
type WorkdayEngine struct{}

func (e *WorkdayEngine) Name() string   { return "Workday" }
func (e *WorkdayEngine) Vendor() string { return "Workday" }

func (e *WorkdayEngine) Analyze(cv models.ParsedCV, jd string) models.ATSResult {
	keywords := extractWeightedKeywords(jd)
	matched, missing, kwScore := weightedExactMatch(cv.RawText, keywords)

	// Workday strictly weights required keywords more — apply penalty for missing required ones
	requiredMissing := 0
	for _, kw := range keywords {
		if kw.Weight >= 0.9 {
			found := false
			for _, m := range matched {
				if m == kw.Word {
					found = true
					break
				}
			}
			if !found {
				requiredMissing++
			}
		}
	}
	// Each missing required keyword is a heavier penalty for Workday
	requiredPenalty := float64(requiredMissing) * 3.0
	kwScore = kwScore - requiredPenalty
	if kwScore < 0 {
		kwScore = 0
	}

	titleScore := calcTitleScore(cv, jd)
	expScore := calcExpScore(cv.YearsExp, jd)
	eduScore := calcEduScore(cv)

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
	if total > 100 {
		total = 100
	}

	var recs []string
	recs = append(recs, buildKeywordRecs(missing, keywords, "exact")...)
	if titleScore < 50 {
		recs = append(recs, "Align your job title more closely with the target role title")
	}
	if formatPenalty > 0 {
		recs = append(recs, "Use a simple single-column ATS-friendly format; avoid tables and special characters")
	}
	if requiredMissing > 0 {
		recs = append(recs, fmt.Sprintf("%d required keywords are missing — Workday heavily penalizes this", requiredMissing))
	}

	return withInterval(models.ATSResult{
		Platform:        e.Name(),
		Vendor:          e.Vendor(),
		Score:           round(total),
		Grade:           models.Grade(total),
		Confidence:      58,
		SimulationNote:  "HiredScore uses A/B/C/D tier grading (not a %) based on career trajectory, tenure, and contextual skill alignment — NOT keyword density. This simulation uses keyword matching as a proxy, which underestimates HiredScore's semantic intelligence. Real score may vary significantly. No auto-reject threshold exists — HiredScore only surfaces prioritization tiers.",
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
			{Category: "Weighted Keyword Match (Exact)", Score: kwScore, MaxScore: 100, Weight: 0.40},
			{Category: "Title Alignment", Score: titleScore, MaxScore: 100, Weight: 0.30},
			{Category: "Experience Years", Score: expScore, MaxScore: 100, Weight: 0.20},
			{Category: "Education", Score: eduScore, MaxScore: 100, Weight: 0.10},
		},
	})
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func calcTitleScore(cv models.ParsedCV, jd string) float64 {
	jdLower := strings.ToLower(jd)
	titleKeywords := []string{"engineer", "developer", "backend", "frontend", "full stack", "sde", "senior", "lead", "architect", "analyst", "scientist"}
	hits := 0
	for _, kw := range titleKeywords {
		if strings.Contains(jdLower, kw) {
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
	required := 2.0
	jdLower := strings.ToLower(jd)
	if strings.Contains(jdLower, "7+") || strings.Contains(jdLower, "7 years") {
		required = 7
	} else if strings.Contains(jdLower, "5+") || strings.Contains(jdLower, "5 years") {
		required = 5
	} else if strings.Contains(jdLower, "3+") || strings.Contains(jdLower, "3 years") {
		required = 3
	} else if strings.Contains(jdLower, "1+") || strings.Contains(jdLower, "1 year") {
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
		if strings.Contains(lower, "master") || strings.Contains(lower, "m.tech") {
			return 100
		}
		if strings.Contains(lower, "bachelor") || strings.Contains(lower, "b.tech") || strings.Contains(lower, "b.e") {
			return 85
		}
	}
	return 60
}

// buildKeywordRecs picks the top missing REQUIRED keywords to recommend
func buildKeywordRecs(missing []string, keywords []WeightedKeyword, matchType string) []string {
	if len(missing) == 0 {
		return nil
	}

	// prioritize high-weight missing keywords
	weightMap := map[string]float64{}
	for _, kw := range keywords {
		weightMap[kw.Word] = kw.Weight
	}

	type mw struct{ word string; w float64 }
	var ranked []mw
	for _, m := range missing {
		ranked = append(ranked, mw{m, weightMap[m]})
	}
	// sort by weight desc
	for i := 0; i < len(ranked)-1; i++ {
		for j := i + 1; j < len(ranked); j++ {
			if ranked[j].w > ranked[i].w {
				ranked[i], ranked[j] = ranked[j], ranked[i]
			}
		}
	}

	top := ranked
	if len(top) > 8 {
		top = top[:8]
	}
	words := make([]string, len(top))
	for i, r := range top {
		words[i] = r.word
	}
	return []string{fmt.Sprintf("Add these missing keywords (%s, highest priority first): %s", matchType, strings.Join(words, ", "))}
}

func round(f float64) float64 {
	return float64(int(f*10+0.5)) / 10
}
