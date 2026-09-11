package main

import (
	"fmt"
	"io/ioutil"
	"github.com/bogusdeck/ats-scanner/internal/ats"
	"github.com/bogusdeck/ats-scanner/internal/parser"
)

func main() {
	cvText, _ := ioutil.ReadFile("cv_test.txt")
	jdText, _ := ioutil.ReadFile("jd_test.txt")

	cv := parser.ParseText(string(cvText))
	engine := ats.Get("Workday")
	res := engine.Analyze(cv, string(jdText))
	
	fmt.Printf("Parsed CV:\n%+v\n", cv)
	fmt.Printf("Result:\n%+v\n", res)
}
