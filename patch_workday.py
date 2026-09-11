import re

with open("internal/ats/workday.go", "r") as f:
    text = f.read()

old_logic = """		if strings.Contains(jdLower, kw) {
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
	return scoreToPercent(float64(hits) / float64(len(titleKeywords)))"""

new_logic = """		if strings.Contains(jdLower, kw) {
			jdKeywords++
			for _, job := range cv.Experience {
				if strings.Contains(strings.ToLower(job.Title), kw) {
					hits++
					break
				}
			}
		}
	}
	if jdKeywords == 0 {
		return 50
	}
	return scoreToPercent(float64(hits) / float64(jdKeywords))"""

text = text.replace("hits := 0\n\tfor", "hits := 0\n\tjdKeywords := 0\n\tfor")
text = text.replace(old_logic, new_logic)

with open("internal/ats/workday.go", "w") as f:
    f.write(text)

