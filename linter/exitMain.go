// Package linter функционал линтера.

// asciicheck Simple linter to check that your code does not contain non-ASCII identifiers
// bodyclose is a static analysis tool which checks whether res.Body is correctly closed.
package linter

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// ExitInMainAnalyzer check for direct os.Exit calls in main function
var ExitInMainAnalyzer = &analysis.Analyzer{
	Name: "exitinmain",
	Doc:  "check for direct os.Exit calls in main function",
	Run:  runExitInMainAnalyzer,
}

func runExitInMainAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fd, okFD := decl.(*ast.FuncDecl)
			if !okFD || fd.Name.Name != "main" {
				continue
			}

			filePath := pass.Fset.Position(file.Package).Filename
			if strings.Contains(filePath, "go-build") {
				continue
			}

			ast.Inspect(fd, func(n ast.Node) bool {
				ce, okCE := n.(*ast.CallExpr)
				if !okCE {
					return true
				}

				se, okSE := ce.Fun.(*ast.SelectorExpr)
				if !okSE {
					return true
				}

				if ident, ok := se.X.(*ast.Ident); ok && ident.Name == "os" && se.Sel.Name == "Exit" {
					pass.Reportf(ce.Pos(), "avoid using os.Exit directly in main function")
				}

				return true
			})
		}
	}
	return nil, nil
}
