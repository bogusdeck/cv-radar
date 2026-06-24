// Package ats provides scoring engines simulating real ATS platform behavior.
package ats

import (
	"github.com/bogusdeck/ats-scanner/internal/models"
)

// Engine is the interface all ATS scoring engines must implement
type Engine interface {
	Name() string
	Vendor() string
	Analyze(cv models.ParsedCV, jd string) models.ATSResult
}

// All returns all registered ATS engines
func All() []Engine {
	return []Engine{
		&WorkdayEngine{},
		&TaleoEngine{},
		&ICIMSEngine{},
		&GreenhouseEngine{},
		&LeverEngine{},
		&SuccessFactorsEngine{},
	}
}

// Get returns a single engine by platform name (case-insensitive)
func Get(name string) Engine {
	for _, e := range All() {
		if equalsCI(e.Name(), name) {
			return e
		}
	}
	return nil
}

func equalsCI(a, b string) bool {
	return toLower(a) == toLower(b)
}
