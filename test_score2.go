package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"bufio"
)

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

func main() {
	cvText, _ := ioutil.ReadFile("cv_test.txt")
	lines := splitLines(string(cvText))

	sectionHeaders := regexp.MustCompile(`(?i)^(experience|work experience|employment|professional experience|projects|education|skills|summary|certifications|achievements)$`)
	dateRe := regexp.MustCompile(`(?i)((?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|\d{1,2}|\d{4})\s*[-–—/to]+\s*(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec|\d{1,2}|\d{4}|present|now|current))`)

	inExp := false
	for _, line := range lines {
		if sectionHeaders.MatchString(line) {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "experience") || strings.Contains(lower, "employment") {
				inExp = true
				fmt.Println("ENTERED EXPERIENCE")
				continue
			} else if inExp {
				fmt.Println("EXITED EXPERIENCE ON:", line)
				inExp = false
			}
		}

		if !inExp {
			continue
		}

		fmt.Println("IN EXP:", line)
		if dateRe.MatchString(line) {
			fmt.Println("  -> MATCHED DATE:", dateRe.FindString(line))
		}
	}
}
