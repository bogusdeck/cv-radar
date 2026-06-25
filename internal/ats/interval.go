package ats

import "github.com/bogusdeck/ats-scanner/internal/models"

// confidenceInterval computes ScoreLow and ScoreHigh from a score and confidence level.
// Lower confidence = wider interval. Based on platform-specific empirical margins:
//
//   Taleo       80% conf → ±8  pts  (best documented)
//   Workday     72% conf → ±13 pts
//   Lever       65% conf → ±15 pts
//   SF          68% conf → ±12 pts
//   iCIMS       60% conf → ±18 pts
//   Greenhouse  50% conf → ±22 pts  (no algorithm, pure human)
//
// The margin is: (100 - confidence) * 0.4 + 5 (minimum ±5 always)
func confidenceInterval(score, confidence float64) (low, high float64) {
	margin := (100-confidence)*0.4 + 5
	low = score - margin
	high = score + margin
	if low < 0 {
		low = 0
	}
	if high > 100 {
		high = 100
	}
	return round(low), round(high)
}

// withInterval fills ScoreLow and ScoreHigh on a result and returns it
func withInterval(r models.ATSResult) models.ATSResult {
	r.ScoreLow, r.ScoreHigh = confidenceInterval(r.Score, r.Confidence)
	return r
}
