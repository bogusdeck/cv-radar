package main

import (
	"fmt"
	"regexp"
)

func main() {
	dateRe := regexp.MustCompile(`(?i)((?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|\d{1,2}|\d{4})\s*[-–—/to]+\s*(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|\d{1,2}|\d{4}|present|now|current))`)
	
	fmt.Println(dateRe.MatchString("Dec 2024 – June 2026"))
	fmt.Println(dateRe.FindString("Dec 2024 – June 2026"))
}
