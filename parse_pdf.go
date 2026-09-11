package main

import (
	"fmt"
	"github.com/bogusdeck/ats-scanner/internal/parser"
)

func main() {
	text, err := parser.ExtractPDF("test.pdf")
	if err != nil {
		panic(err)
	}
	fmt.Println("=== RAW TEXT ===")
	fmt.Println(text)
	
	cv := parser.ParseText(text)
	fmt.Printf("=== PARSED CV ===\n%+v\n", cv)
}
