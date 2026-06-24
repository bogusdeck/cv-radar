package models

// ParsedCV holds structured data extracted from a CV
type ParsedCV struct {
	RawText    string   `json:"raw_text"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Location   string   `json:"location"`
	Summary    string   `json:"summary"`
	Skills     []string `json:"skills"`
	Experience []Job    `json:"experience"`
	Projects   []Project `json:"projects"`
	Education  []Edu    `json:"education"`
	Links      []string `json:"links"`
	YearsExp   float64  `json:"years_exp"`
}

type Job struct {
	Title       string   `json:"title"`
	Company     string   `json:"company"`
	Duration    string   `json:"duration"`
	Description []string `json:"description"`
}

type Project struct {
	Name        string   `json:"name"`
	Tech        []string `json:"tech"`
	Description []string `json:"description"`
}

type Edu struct {
	Degree      string `json:"degree"`
	Institution string `json:"institution"`
	Year        string `json:"year"`
	GPA         string `json:"gpa"`
}

// AnalysisRequest is the API request payload
type AnalysisRequest struct {
	CVText   string `json:"cv_text"`
	JDText   string `json:"jd_text"`
	Platform string `json:"platform,omitempty"` // empty = all platforms
}

// AnalysisResponse is the full API response
type AnalysisResponse struct {
	ParsedCV ParsedCV              `json:"parsed_cv"`
	Results  []ATSResult           `json:"results"`
}

// ATSResult holds the score and breakdown for one ATS platform
type ATSResult struct {
	Platform        string          `json:"platform"`
	Vendor          string          `json:"vendor"`
	Score           float64         `json:"score"`           // 0–100
	Grade           string          `json:"grade"`           // A, B, C, D, F
	KeywordScore    float64         `json:"keyword_score"`
	StructureScore  float64         `json:"structure_score"`
	ExperienceScore float64         `json:"experience_score"`
	EducationScore  float64         `json:"education_score"`
	MatchedKeywords []string        `json:"matched_keywords"`
	MissingKeywords []string        `json:"missing_keywords"`
	Warnings        []string        `json:"warnings"`
	Recommendations []string        `json:"recommendations"`
	AutoReject      bool            `json:"auto_reject"`
	Breakdown       []ScoreBreakdown `json:"breakdown"`
}

type ScoreBreakdown struct {
	Category string  `json:"category"`
	Score    float64 `json:"score"`
	MaxScore float64 `json:"max_score"`
	Weight   float64 `json:"weight"`
}

func Grade(score float64) string {
	switch {
	case score >= 85:
		return "A"
	case score >= 70:
		return "B"
	case score >= 55:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}
