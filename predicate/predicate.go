package predicate

import (
	"github.com/murphybytes/dsl/internal/ast"
)



func Evaluate( expression string, _ map[string]interface{})(bool,error) {
	parser := ast.Parser()
	var t ast.Expression
	if err := parser.ParseString("", expression, &t); err != nil {
		return false, err
	}
	result, err := t.Eval()
	if err != nil {
		return false, err
	}
	return bool(*result.Bool), nil
}