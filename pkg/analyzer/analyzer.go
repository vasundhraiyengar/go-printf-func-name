package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goprintffuncname",
	Doc:      "Checks that printf-like functions are named with `f` at the end.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		ce := node.(*ast.CallExpr)

		isLog := isPkgDot(ce.Fun, "log", "Log")
		isLogWith := isPkgDot(ce.Fun, "log", "With")
		if !(isLog || isLogWith) || len(ce.Args) != 1 {
			return
		}

		params := ce.Args
		if isLog && len(params)%2 == 0 {
			return
		}

		if isLogWith && len(params)%2 != 0 {
			return
		}

		pass.Reportf(node.Pos(), "need right number of args")

	})

	return nil, nil
}

func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}
