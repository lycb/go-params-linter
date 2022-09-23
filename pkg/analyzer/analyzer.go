package analyzer

import (
	"fmt"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goparamslinter",
	Doc:      "Check if multiple params have the same type",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type (
	recvFunc struct {
		Index    int
		FuncName string
	}
)

const NAME = 0
const TYPE = 1

func run(pass *analysis.Pass) (any, error) {
	// a decorated ast package "dst" is used to avoid free floating comment issue
	// see https://github.com/golang/go/issues/20744
	dec := decorator.NewDecorator(pass.Fset) // holds mapping between ast and dst

	for _, file := range pass.Files {
		f, err := dec.DecorateFile(file) // transform ast file to dst file
		if err != nil {
			panic(err)
		}

		funcDecl := f.Decls[0].(*dst.FuncDecl)

		// a map of string to save the previous param type
		var tmpmap = make(map[string]float64)

		// get a list of params
		params := funcDecl.Type.Params.List

		for i := 0; i < len(params); i++ {
			var firstParam []string
			var secondParam []string

			if i == 0 && i+1 < (len(params)) { // first param
				firstParam = getParamInfo(params[i])
				secondParam = getParamInfo(params[i+1])

				if tmpmap[firstParam[TYPE]] > 0 {
					fixParams(pass, funcDecl, firstParam, secondParam)
				}
				// if types for 1st and 2nd param != reset first param counter to 1 and set param counter for 2 param to 1
				tmpmap[firstParam[TYPE]] = 0
				tmpmap[secondParam[TYPE]] = 1

			} else if i > 0 { // if has previous param check
				firstParam = getParamInfo(params[i-1])
				secondParam = getParamInfo(params[i])

				if tmpmap[firstParam[TYPE]] > 0 {
					fixParams(pass, funcDecl, firstParam, secondParam)
				}
				// reset the previous param count to 0
				tmpmap[firstParam[TYPE]] = 0
			}
		}

	}

	return nil, nil
}

// Helper functions
func getParamType(param *dst.Field) string {
	if ident, ok := param.Type.(*dst.Ident); ok {
		return ident.Name
	}
	return ""
}

func getParamName(param *dst.Field) string {
	for i := 0; i < len(param.Names); i++ {
		return param.Names[0].Name
	}
	return ""
}

func getParamInfo(param *dst.Field) []string {
	return []string{getParamName(param), getParamType(param)}
}

func fixParams(pass *analysis.Pass, funcDecl *dst.FuncDecl, firstParam, secondParam []string) {
	pass.Report(analysis.Diagnostic{
		Message: fmt.Sprintf("param '%s' with type '%s' for function '%s' should be combined with param '%s' of type '%s'\n",
			firstParam[NAME],
			firstParam[TYPE],
			funcDecl.Name.Name,
			secondParam[NAME],
			secondParam[TYPE]),
		// SuggestedFixes: []analysis.SuggestedFix{
		// 	{
		// 		Message: fmt.Sprintf("should replace `%s` with `%s`", oldExpr, newExpr),
		// TextEdits: []analysis.TextEdit{ // will only call this with -fix flag
		// 	{
		// 		Pos:     node.Pos(),
		// 		End:     node.End(),
		// 		NewText: []byte(newExpr),
		// 	},
		// },
		// },
		// },
	})
}
