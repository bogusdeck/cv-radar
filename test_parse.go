package main

import (
	"fmt"
	"regexp"
)

func main() {
	sectionHeaders := regexp.MustCompile(`(?i)^(experience|work experience|employment|professional experience|projects|education|skills|summary|certifications|achievements)$`)
	fmt.Println(sectionHeaders.MatchString("Experience"))
	fmt.Println(sectionHeaders.MatchString("Experience "))
}
