package analyzer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goparamslinter",
	Doc:      "Check if multiple params have the same type",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (any, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}
	inspector.Preorder(nodeFilter, func(node ast.Node) {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return
		}

		// a map of string to save the previous param type
		var tmpmap = make(map[string]float64)

		// get a list of params
		params := funcDecl.Type.Params.List

		// check every 2 params
		for i := 0; i < len(params); i += 2 {
			firstParamType, ok := params[i].Type.(*ast.Ident)
			if !ok { // if type is not an identifier
				return
			}
			firstParamName := params[i].Names[0].Name

			// if there is a param before the current index
			if i > 0 {
				previousParamType, ok := params[i-1].Type.(*ast.Ident)
				if !ok { // if type is not an identifier
					return
				}
				previousParamName := params[i-1].Names[0].Name

				if tmpmap[firstParamType.Name] > 0 {
					oldExpr, newExpr := formatFixExpr(pass, node, funcDecl, firstParamType.Name)
					fixParams(pass, node, funcDecl.Name.Name, previousParamName, previousParamType.Name, firstParamName, firstParamType.Name, oldExpr, newExpr)
					return
				}
				// reset the previous param count to 0
				tmpmap[previousParamType.Name] = 0
			}

			tmpmap[firstParamType.Name] = 1

			// if there is a second param
			if i+1 < len(params) {
				secondParamType, ok := params[i+1].Type.(*ast.Ident)
				if !ok { // if type is not an identifier
					return
				}
				secondParamName := params[i+1].Names[0].Name

				// fail if type already exist in map with count 1
				if tmpmap[secondParamType.Name] > 0 {
					oldExpr, newExpr := formatFixExpr(pass, node, funcDecl, firstParamType.Name)
					fixParams(pass, node, funcDecl.Name.Name, firstParamName, firstParamType.Name, secondParamName, secondParamType.Name, oldExpr, newExpr)
					return
				}
				// if types for 1st and 2nd param != reset first param counter to 1 and set param counter for 2 param to 1
				tmpmap[firstParamType.Name] = 0
				tmpmap[secondParamType.Name] = 1
			}
		}
	})
	return nil, nil
}

// Helper functions

func render(fset *token.FileSet, x any) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, x); err != nil {
		panic(err)
	}
	return buf.String()
}

func fixParams(pass *analysis.Pass, node ast.Node, funcName, firstName, firstType, secondName, secondType, oldExpr, newExpr string) {
	pass.Report(analysis.Diagnostic{
		Pos:     node.Pos(),
		Message: fmt.Sprintf("param '%s' with type '%s' for function '%s' should be combined with param '%s' of type '%s'\n", secondName, secondType, funcName, firstName, firstType),
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fmt.Sprintf("should replace `%s` with `%s`", oldExpr, newExpr),
				TextEdits: []analysis.TextEdit{ // will only call this with -fix flag
					{
						Pos:     node.Pos(),
						End:     node.End(),
						NewText: []byte(newExpr),
					},
				},
			},
		},
	})
}

func formatFixExpr(pass *analysis.Pass, node ast.Node, funcDecl *ast.FuncDecl, paramType string) (string, string) {
	// get a list of old params and remove keyword func
	oldParam := render(pass.Fset, funcDecl.Type)
	oldParam = strings.Trim(oldParam, "func")

	// remove return type if there is one
	position := strings.LastIndex(oldParam, ") (")
	oldParam = oldParam[:position+1]

	// get the whole old function
	oldExpr := render(pass.Fset, node)

	// handle for when function has preceding comments
	paramPosition := strings.LastIndex(oldExpr, "func")
	if strings.HasPrefix(oldExpr, "//") {
		oldExpr = oldExpr[paramPosition:]
	}

	// replace old params with new params
	newParam := strings.Replace(oldParam, " "+paramType, "", 1)
	newExpr := strings.Replace(oldExpr, oldParam, newParam, 1)

	return oldExpr, newExpr
}
