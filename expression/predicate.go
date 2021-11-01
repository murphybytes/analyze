package expression

import (
	"github.com/murphybytes/analyze/internal/ast"
)

// Evaluate takes
func Evaluate(ctx ast.Context, expression string) (bool, error) {
	parser := ast.Parser()
	var t ast.Expression
	if err := parser.ParseString("", expression, &t); err != nil {
		return false, err
	}
	result, err := t.Eval(ctx)
	if err != nil {
		return false, err
	}
	return bool(*result.Bool), nil
}
