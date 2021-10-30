package expression

import (
	"github.com/murphybytes/analyze/context"
	"github.com/murphybytes/analyze/internal/ast"
)

func Evaluate(expression string, ctx context.Context) (bool, error) {
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
