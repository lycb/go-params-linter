package analyzer

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "goparamslinter",
	Doc:  "Check if params have the same type",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	inspect := func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)

		if !ok {
			return true
		}

		var tmpmap = make(map[string]float64)

		params := funcDecl.Type.Params.List

		for i := 0; i < len(params); i += 2 {
			firstParamName := params[i].Names[0].Name
			firstParamType, ok := params[i].Type.(*ast.Ident)
			if !ok {

				fmt.Printf("What1")
				return true
			}

			if i > 0 {
				previousParamName := params[i-1].Names[0].Name
				previousParamType, ok := params[i-1].Type.(*ast.Ident)
				if !ok {
					fmt.Printf("What3")
					return true
				}

				if tmpmap[firstParamType.Name] > 0 {
					pass.Reportf(node.Pos(), "param '%s' with type '%s' for function '%s' should be combined with param '%s' of type '%s'\n",
						firstParamName, firstParamType.Name, funcDecl.Name.Name, previousParamName, previousParamType.Name)
					return true
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

					fmt.Printf("What4")
					return true
				}

				// fail if type already exist in map with count 1
				if tmpmap[secondParamType.Name] > 0 {
					pass.Reportf(node.Pos(), "param '%s' with type '%s' for function '%s' should be combined with param '%s' of type '%s'\n",
						secondParamName, secondParamType.Name, funcDecl.Name.Name, firstParamName, firstParamType.Name)
					return true
				}
				// if types for 1st and 2nd param != reset first param counter to 1 and set param counter for 2 param to 1
				tmpmap[firstParamType.Name] = 0
				tmpmap[secondParamType.Name] = 1
			}
		}
		return true
	}
	for _, f := range pass.Files {
		ast.Inspect(f, inspect)
	}
	return nil, nil
}
