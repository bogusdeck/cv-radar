import re

with open("internal/parser/text.go", "r") as f:
    text = f.read()

# Replace dateRe with a robust one that captures "Month Year - Month Year"
target = r"dateRe := regexp.MustCompile(`(?i)((?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|\d{1,2}|\d{4})\s*[-–—/to]+\s*(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|\d{1,2}|\d{4}|present|now|current))`)"
new_code = r"dateRe := regexp.MustCompile(`(?i)((?:(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+)?\d{2,4}\s*[-–—/to]+\s*(?:(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+)?(?:\d{2,4}|present|now|current))`)"

text = text.replace(target, new_code)

with open("internal/parser/text.go", "w") as f:
    f.write(text)

