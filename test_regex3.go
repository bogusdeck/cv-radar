package main

import (
	"fmt"
	"regexp"
)

func main() {
	dateRe := regexp.MustCompile(`(?i)((?:(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+)?\d{2,4}\s*[-–—/to]+\s*(?:(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+)?(?:\d{2,4}|present|now|current))`)
	fmt.Println(dateRe.FindString("Dec 2024 – June 2026"))
}
