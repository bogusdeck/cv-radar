import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

target = r"""	// Remove common problematic unicode characters"""

new_code = r"""	// Sanitize common LLM hallucinations
	output = regexp.MustCompile(`(?s)\\resumeSummary\{(.*?)\}`).ReplaceAllString(output, `$1`)
	output = regexp.MustCompile(`(?s)\\resumeProjectDescription\{(.*?)\}`).ReplaceAllString(output, `\\resumeItem{$1}`)
	output = regexp.MustCompile(`(?s)\\resumeSkillsBullet\{(.*?)\}`).ReplaceAllString(output, `\\resumeItem{$1}`)
	output = strings.ReplaceAll(output, "\\resumeProjectStart", "\\resumeSubHeadingListStart")
	output = strings.ReplaceAll(output, "\\resumeProjectEnd", "\\resumeSubHeadingListEnd")
	
	// Remove common problematic unicode characters"""

text = text.replace(target, new_code)

# Add regexp import if missing
if '"regexp"' not in text:
    text = text.replace('"path/filepath"', '"path/filepath"\n\t"regexp"')

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

