import re

with open("internal/parser/text.go", "r") as f:
    text = f.read()

# Replace double backslashes with single backslashes in dateRe
target = r"\\d{1,2}|\\d{4})\\s*[-–—/to]+\\s*(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|\\d{1,2}|\\d{4}|present|now|current))"
new_code = r"\d{1,2}|\d{4})\s*[-–—/to]+\s*(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|\d{1,2}|\d{4}|present|now|current))"

text = text.replace(target, new_code)

with open("internal/parser/text.go", "w") as f:
    f.write(text)

