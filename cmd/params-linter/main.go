package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	v := Visitor{fset: token.NewFileSet()}
	for _, filePath := range os.Args[1:] {
		if filePath == "--" {
			continue
		}

		f, err := parser.ParseFile(v.fset, filePath, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %s", filePath, err)
		}

		ast.Walk(&v, f)
	}
}

type Visitor struct {
	fset *token.FileSet
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)

	if !ok {
		return v
	}

	var tmpmap = make(map[string]float64)

	params := funcDecl.Type.Params.List

	for i := 0; i < len(params); i += 2 {
		firstParamName := params[i].Names[0].Name
		firstParamType, ok := params[i].Type.(*ast.Ident)
		if !ok {
			return v
		}

		if tmpmap[firstParamType.Name] > 0 {
			return v
		} else {
			if i > 0 {
				previousParam, ok := params[i-1].Type.(*ast.Ident)
				if !ok {
					fmt.Printf("first param !ok")
					return v
				}
				// reset the previous second param count to 0
				tmpmap[previousParam.Name] = 0
			}
			tmpmap[firstParamType.Name] = 1
		}

		// if there is a second param
		if i+1 < len(params) {
			secondParamName := params[i+1].Names[0].Name
			secondParamType, ok := params[i+1].Type.(*ast.Ident)
			if !ok {
				return v
			}

			// fail if type already exist in map with count 1
			if tmpmap[secondParamType.Name] > 0 {
				fmt.Printf("%s: param '%s' with type '%s' for function '%s' should be combined with param '%s' of type '%s'\n",
					v.fset.Position(node.Pos()), secondParamName, secondParamType.Name, funcDecl.Name.Name, firstParamName, firstParamType.Name)
				return v
			} else {
				// if types for 1st and 2nd param != reset first param counter to 1 and set param counter for 2 param to 1
				tmpmap[firstParamType.Name] = 0
				tmpmap[secondParamType.Name] = 1
			}
		}
	}

	return v
}
