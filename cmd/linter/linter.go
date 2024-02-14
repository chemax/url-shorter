// Package main запускать go run cmd/linter/linter.go ./...
package main

import (
	"github.com/chemax/url-shorter/linter"
	"github.com/tdakkota/asciicheck"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var analyzers []*analysis.Analyzer
	//analyzers = append(analyzers, linter.ErrCheckAnalyzer)
	analyzers = append(analyzers, linter.ExitInMainAnalyzer)
	//Simple linter to check that your code does not contain non-ASCII identifiers
	analyzers = append(analyzers, asciicheck.NewAnalyzer())
	// bodyclose is a static analysis tool which checks whether res.Body is correctly closed.
	analyzers = append(analyzers, bodyclose.Analyzer)
	analyzers = append(analyzers, linter.X...)
	for _, check := range staticcheck.Analyzers {
		analyzers = append(analyzers, check.Analyzer)
	}
	for _, fix := range quickfix.Analyzers {
		analyzers = append(analyzers, fix.Analyzer)
	}
	multichecker.Main(analyzers...)
}
