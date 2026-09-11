import re

with open("internal/parser/text.go", "r") as f:
    text = f.read()

new_func = r"""func extractExperience(text string) []models.Job {
	var jobs []models.Job
	lines := splitLines(text)

	inExp := false
	var current *models.Job
	prevLine := ""

	dateRe := regexp.MustCompile(`(?i)((?:(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+)?\d{2,4}\s*[-–—/to]+\s*(?:(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+)?(?:\d{2,4}|present|now|current))`)

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

		if dateRe.MatchString(line) && len(line) < 35 {
			if current != nil {
				jobs = append(jobs, *current)
			}
			current = &models.Job{Duration: dateRe.FindString(line)}
			
			cleaned := strings.TrimSpace(dateRe.ReplaceAllString(line, ""))
			if cleaned != "" {
				current.Title = cleaned
			} else if prevLine != "" && !strings.HasPrefix(prevLine, "•") && !strings.HasPrefix(prevLine, "-") {
				current.Title = prevLine
			}
			prevLine = line
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
		prevLine = line
	}

	if current != nil {
		jobs = append(jobs, *current)
	}
	return jobs
}"""

# Manually find the start and end of extractExperience
start = text.find("func extractExperience(text string) []models.Job {")
if start != -1:
    end = text.find("\n}\n", start) + 3
    text = text[:start] + new_func + "\n" + text[end:]

with open("internal/parser/text.go", "w") as f:
    f.write(text)
