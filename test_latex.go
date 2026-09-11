package main

import (
	"fmt"
	"io/ioutil"
	"github.com/bogusdeck/ats-scanner/internal/parser"
)

func main() {
	latex, _ := ioutil.ReadFile("test.tex")
	stripped := parser.StripLatex(string(latex))
	fmt.Println("=== STRIPPED LATEX ===")
	fmt.Println(stripped)
	
	cv := parser.ParseText(stripped)
	fmt.Printf("=== PARSED CV ===\n%+v\n", cv)
}
