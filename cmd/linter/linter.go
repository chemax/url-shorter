package main

import (
	"github.com/chemax/url-shorter/linter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var analyzers []*analysis.Analyzer
	//TODO как получить ВСЕ анализаторы из x/analysis??????
	//analyzers = append(analyzers, linter.ErrCheckAnalyzer)
	analyzers = append(analyzers, linter.ExitInMainAnalyzer)
	analyzers = append(analyzers, linter.X...)
	for _, check := range staticcheck.Analyzers {
		analyzers = append(analyzers, check.Analyzer)
	}
	for _, fix := range quickfix.Analyzers {
		analyzers = append(analyzers, fix.Analyzer)
	}
	//TODO внешние анализаторы

	multichecker.Main(analyzers...)
}
