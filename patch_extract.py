import re

with open("internal/parser/text.go", "r") as f:
    text = f.read()

old_logic = """		for _, line := range lines {
			if sectionHeaders.MatchString(line) {"""

new_logic = """		prevLine := ""
		for _, line := range lines {
			if sectionHeaders.MatchString(line) {"""

text = text.replace(old_logic, new_logic)

old_date_logic = """			if dateRe.MatchString(line) {
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
			}"""

new_date_logic = """			if dateRe.MatchString(line) {
				if current != nil {
					jobs = append(jobs, *current)
				}
				current = &models.Job{Duration: dateRe.FindString(line)}
				// Title is often on the same line or previous line
				cleaned := strings.TrimSpace(dateRe.ReplaceAllString(line, ""))
				if cleaned != "" {
					current.Title = cleaned
				} else if prevLine != "" && !strings.HasPrefix(prevLine, "•") && !strings.HasPrefix(prevLine, "-") {
					current.Title = prevLine
				}
				prevLine = line
				continue
			}"""

text = text.replace(old_date_logic, new_date_logic)

old_end = """				}
			}
		}

		if current != nil {"""

new_end = """				}
			}
			prevLine = line
		}

		if current != nil {"""

text = text.replace(old_end, new_end)

with open("internal/parser/text.go", "w") as f:
    f.write(text)

