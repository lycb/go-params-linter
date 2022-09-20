package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goparamslinter",
	Doc:      "Check if params have the same type",
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

		var tmpmap = make(map[string]float64)

		params := funcDecl.Type.Params.List

		for i := 0; i < len(params); i += 2 {
			firstParamName := params[i].Names[0].Name
			firstParamType, ok := params[i].Type.(*ast.Ident)
			if !ok {
				return
			}

			if i > 0 {
				previousParamName := params[i-1].Names[0].Name
				previousParamType, ok := params[i-1].Type.(*ast.Ident)
				if !ok {
					return
				}

				if tmpmap[firstParamType.Name] > 0 {
					pass.Reportf(node.Pos(), "param '%s' with type '%s' for function '%s' should be combined with param '%s' of type '%s'\n",
						firstParamName, firstParamType.Name, funcDecl.Name.Name, previousParamName, previousParamType.Name)
					return
				}
				// reset the previous second param count to 0
				tmpmap[previousParamType.Name] = 0
			}

			tmpmap[firstParamType.Name] = 1

			// if there is a second param
			if i+1 < len(params) {
				secondParamName := params[i+1].Names[0].Name
				secondParamType, ok := params[i+1].Type.(*ast.Ident)
				if !ok {
					return
				}

				// fail if type already exist in map with count 1
				if tmpmap[secondParamType.Name] > 0 {
					pass.Reportf(node.Pos(), "param '%s' with type '%s' for function '%s' should be combined with param '%s' of type '%s'\n",
						secondParamName, secondParamType.Name, funcDecl.Name.Name, firstParamName, firstParamType.Name)
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
