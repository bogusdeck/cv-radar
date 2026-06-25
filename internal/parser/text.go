package parser

import (
	"bufio"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/bogusdeck/ats-scanner/internal/models"
)

// ParseText takes raw CV text (from PDF, LaTeX, or plain input) and returns a structured ParsedCV
func ParseText(raw string) models.ParsedCV {
	cv := models.ParsedCV{RawText: raw}
	lines := splitLines(raw)

	cv.Name = extractName(lines)
	cv.Email = extractEmail(raw)
	cv.Phone = extractPhone(raw)
	cv.Location = extractLocation(raw)
	cv.Links = extractLinks(raw)
	cv.Summary = extractSection(raw, []string{"summary", "objective", "profile"})
	cv.Skills = extractSkills(raw)
	cv.Experience = extractExperience(raw)
	cv.Projects = extractProjects(raw)
	cv.Education = extractEducation(raw)
	cv.YearsExp = calculateYearsExp(cv.Experience)

	return cv
}

func splitLines(text string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func extractName(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	// First non-empty line is usually the name
	for _, line := range lines[:min(5, len(lines))] {
		// Skip lines that look like contact info
		if strings.Contains(line, "@") || strings.Contains(line, "http") ||
			regexp.MustCompile(`\d{10}`).MatchString(line) {
			continue
		}
		if len(strings.Fields(line)) >= 2 && len(strings.Fields(line)) <= 5 {
			return line
		}
	}
	return lines[0]
}

func extractEmail(text string) string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	m := re.FindString(text)
	return m
}

func extractPhone(text string) string {
	re := regexp.MustCompile(`[\+]?[(]?[0-9]{1,4}[)]?[-\s\.]?[0-9]{3,5}[-\s\.]?[0-9]{4,6}`)
	m := re.FindString(text)
	return m
}

func extractLocation(text string) string {
	// Common Indian city patterns + general location
	re := regexp.MustCompile(`(?i)(Delhi|Mumbai|Bangalore|Bengaluru|Hyderabad|Chennai|Pune|Noida|Gurugram|Kolkata|Remote)[,\s]*(?:India)?`)
	m := re.FindString(text)
	return strings.TrimSpace(m)
}

func extractLinks(text string) []string {
	re := regexp.MustCompile(`https?://[^\s]+`)
	return re.FindAllString(text, -1)
}

func extractSection(text string, headers []string) string {
	lower := strings.ToLower(text)
	for _, h := range headers {
		idx := strings.Index(lower, h)
		if idx == -1 {
			continue
		}
		// Grab up to 500 chars after the header
		start := idx + len(h)
		end := start + 500
		if end > len(text) {
			end = len(text)
		}
		chunk := strings.TrimSpace(text[start:end])
		// Cut off at next known section
		sections := []string{"experience", "education", "skills", "projects", "work", "employment"}
		for _, s := range sections {
			if si := strings.Index(strings.ToLower(chunk), s); si > 20 {
				chunk = chunk[:si]
				break
			}
		}
		return strings.TrimSpace(chunk)
	}
	return ""
}

// knownSkills is a broad taxonomy used for skill extraction
var knownSkills = []string{
	"python", "javascript", "typescript", "go", "golang", "c++", "java", "ruby", "rust", "scala",
	"sql", "nosql", "bash", "shell",
	"django", "fastapi", "flask", "node.js", "nodejs", "react", "react.js", "next.js", "nextjs",
	"angularjs", "angular", "vue", "express", "express.js",
	"celery", "elasticsearch", "redis", "postgresql", "mysql", "mongodb", "sqlite",
	"docker", "kubernetes", "k8s", "nginx", "aws", "gcp", "azure", "digitalocean",
	"git", "github", "gitlab", "ci/cd", "jenkins", "terraform",
	"graphql", "rest", "grpc", "api", "microservices",
	"machine learning", "deep learning", "llm", "openai", "langchain", "rag",
	"pandas", "numpy", "matplotlib", "scikit-learn", "tensorflow", "pytorch",
	"kafka", "rabbitmq", "mq", "celery",
	"postman", "linux", "neovim", "vim",
	"react native", "flutter",
	"tailwind", "tailwind css", "css", "html",
}

func extractSkills(text string) []string {
	lower := strings.ToLower(text)
	found := map[string]bool{}
	for _, skill := range knownSkills {
		if strings.Contains(lower, skill) {
			found[skill] = true
		}
	}

	// Also extract from "Skills" section lines
	skillSection := extractSection(text, []string{"technical skills", "skills", "technologies"})
	if skillSection != "" {
		// Split by common delimiters
		parts := regexp.MustCompile(`[,|•·\n]+`).Split(skillSection, -1)
		for _, p := range parts {
			p = strings.TrimSpace(p)
			p = strings.ToLower(p)
			if len(p) > 1 && len(p) < 40 && !strings.ContainsAny(p, "{}\\") {
				found[p] = true
			}
		}
	}

	skills := make([]string, 0, len(found))
	for k := range found {
		skills = append(skills, k)
	}
	return skills
}

var sectionHeaders = regexp.MustCompile(`(?i)^(experience|work experience|employment|professional experience|projects|education|skills|summary|certifications|achievements)$`)

func extractExperience(text string) []models.Job {
	var jobs []models.Job
	lines := splitLines(text)

	inExp := false
	var current *models.Job

	dateRe := regexp.MustCompile(`(?i)(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|\d{4})\s*[-–—]\s*(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|\d{4}|present|now|current)`)

	for _, line := range lines {
		if sectionHeaders.MatchString(line) {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "experience") || strings.Contains(lower, "employment") {
				inExp = true
				if current != nil {
					jobs = append(jobs, *current)
					current = nil
				}
				continue
			} else if inExp {
				break // hit another section
			}
		}

		if !inExp {
			continue
		}

		if dateRe.MatchString(line) {
			if current != nil {
				jobs = append(jobs, *current)
			}
			current = &models.Job{Duration: dateRe.FindString(line)}
			// Title is often on the same line or previous line
			cleaned := strings.TrimSpace(dateRe.ReplaceAllString(line, ""))
			if cleaned != "" {
				current.Title = cleaned
			}
			continue
		}

		if current != nil {
			if current.Title == "" && isTitle(line) {
				current.Title = line
			} else if current.Company == "" && isCompany(line) {
				current.Company = line
			} else if strings.HasPrefix(line, "•") || strings.HasPrefix(line, "-") || len(line) > 30 {
				desc := strings.TrimLeft(line, "•- ")
				current.Description = append(current.Description, desc)
			}
		}
	}

	if current != nil {
		jobs = append(jobs, *current)
	}
	return jobs
}

func isTitle(line string) bool {
	titleKeywords := []string{"engineer", "developer", "intern", "manager", "analyst", "architect", "lead", "sde", "swe"}
	lower := strings.ToLower(line)
	for _, kw := range titleKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func isCompany(line string) bool {
	return len(line) > 2 && len(line) < 60 && !strings.Contains(line, "@") && unicode.IsUpper(rune(line[0]))
}

func extractProjects(text string) []models.Project {
	var projects []models.Project
	lines := splitLines(text)

	inProj := false
	var current *models.Project

	for _, line := range lines {
		if sectionHeaders.MatchString(line) {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "project") {
				inProj = true
				if current != nil {
					projects = append(projects, *current)
					current = nil
				}
				continue
			} else if inProj {
				break
			}
		}

		if !inProj {
			continue
		}

		// Detect project name (bold line or line with | delimiter)
		if strings.Contains(line, "|") || (len(line) < 60 && !strings.HasPrefix(line, "•") && !strings.HasPrefix(line, "-")) {
			if current != nil {
				projects = append(projects, *current)
			}
			parts := strings.Split(line, "|")
			name := strings.TrimSpace(parts[0])
			current = &models.Project{Name: name}
			if len(parts) > 1 {
				techs := regexp.MustCompile(`[,\s]+`).Split(strings.TrimSpace(parts[1]), -1)
				current.Tech = techs
			}
		} else if current != nil {
			desc := strings.TrimLeft(line, "•- ")
			current.Description = append(current.Description, desc)
		}
	}

	if current != nil {
		projects = append(projects, *current)
	}
	return projects
}

func extractEducation(text string) []models.Edu {
	var edu []models.Edu
	lines := splitLines(text)

	inEdu := false
	var current *models.Edu

	gpaRe := regexp.MustCompile(`(?i)(cgpa|gpa|grade)[:\s]*([0-9.]+)`)
	yearRe := regexp.MustCompile(`\b(19|20)\d{2}\b`)

	for _, line := range lines {
		if sectionHeaders.MatchString(line) {
			if strings.ToLower(line) == "education" {
				inEdu = true
				if current != nil {
					edu = append(edu, *current)
					current = nil
				}
				continue
			} else if inEdu {
				break
			}
		}

		if !inEdu {
			continue
		}

		lower := strings.ToLower(line)
		if strings.Contains(lower, "bachelor") || strings.Contains(lower, "master") ||
			strings.Contains(lower, "b.tech") || strings.Contains(lower, "b.e.") ||
			strings.Contains(lower, "xii") || strings.Contains(lower, "10th") {
			if current != nil {
				edu = append(edu, *current)
			}
			current = &models.Edu{Degree: line}
			continue
		}

		if current != nil {
			if gpaRe.MatchString(line) {
				m := gpaRe.FindStringSubmatch(line)
				if len(m) > 2 {
					current.GPA = m[2]
				}
			}
			if yearRe.MatchString(line) {
				current.Year = yearRe.FindString(line)
			}
			if current.Institution == "" && len(line) > 3 && len(line) < 80 {
				current.Institution = line
			}
		}
	}

	if current != nil {
		edu = append(edu, *current)
	}
	return edu
}

func calculateYearsExp(jobs []models.Job) float64 {
	if len(jobs) == 0 {
		return 0
	}

	// Match full date patterns like "Dec 2024", "Aug 2024", "2023", "Present"
	fullDateRe := regexp.MustCompile(`(?i)(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[\s\-]+(\d{4})|(\d{4})|(present|now|current)`)

	total := 0.0
	for _, job := range jobs {
		dur := job.Duration
		if dur == "" {
			continue
		}

		matches := fullDateRe.FindAllStringSubmatch(dur, -1)
		if len(matches) < 1 {
			continue
		}

		// Parse start date from first match
		start := parseDateMatch(matches[0])
		if start.IsZero() {
			continue
		}

		// Parse end date
		var end time.Time
		if len(matches) >= 2 {
			end = parseDateMatch(matches[len(matches)-1])
		}
		if end.IsZero() {
			end = time.Now()
		}

		diff := end.Sub(start).Hours() / 8760.0
		if diff > 0 && diff < 20 { // sanity cap
			total += diff
		}
	}

	// Cap total at 30 years
	if total > 30 {
		total = 30
	}
	return math.Round(total*10) / 10
}

func parseDateMatch(m []string) time.Time {
	// m[0]=full, m[1]=month abbr, m[2]=year-after-month, m[3]=bare-year, m[4]=present
	monthMap := map[string]int{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
		"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
	}

	// present/now/current
	if len(m) > 4 && strings.ToLower(m[4]) != "" {
		return time.Now()
	}

	// "Dec 2024" style
	if len(m) > 2 && m[1] != "" && m[2] != "" {
		mo := monthMap[strings.ToLower(m[1])]
		yr, err := strconv.Atoi(m[2])
		if err == nil && yr > 1990 && yr <= time.Now().Year()+1 {
			return time.Date(yr, time.Month(mo), 1, 0, 0, 0, 0, time.UTC)
		}
	}

	// bare year "2023"
	if len(m) > 3 && m[3] != "" {
		yr, err := strconv.Atoi(m[3])
		if err == nil && yr > 1990 && yr <= time.Now().Year()+1 {
			return time.Date(yr, 6, 1, 0, 0, 0, 0, time.UTC) // mid-year estimate
		}
	}

	return time.Time{} // zero = unknown
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
