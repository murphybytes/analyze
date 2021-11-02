package context

import (
	"github.com/murphybytes/analyze/errors"
	"github.com/murphybytes/analyze/internal/ast"
)

// @len(arr) returns the length of an array
func _len(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New(errors.SyntaxError, "wrong number of arguments for len expected 1 got %d", len(args))
	}
	arr, ok := args[0].([]interface{})
	if !ok {
		return nil, errors.New(errors.TypeMismatch, "expected array got %T for len function", args[0])
	}
	return len(arr), nil
}

// @select(arr, expression) returns a subset of arr such that elements of the subset are such that expression is true.
// note that each element can be referenced in the expression by a variable, for example: "$foo == 3" would return each
// element in an array that was equal to three.  You would use $foo.bar == "complete" to reference the field "bar" in an
// array of objects, returning all objects where "bar" == "complete"
func _select(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New(errors.SyntaxError, "wrong number of arguments for select, expected 2 got %d", len(args))
	}
	arr, ok := args[0].([]interface{})
	if !ok {
		return nil, errors.New(errors.TypeMismatch, "expected array got %T for first argument of select", args[0])
	}
	expression, ok := args[1].(string)
	if !ok {
		return nil, errors.New(errors.TypeMismatch, "expected string argument for second argument of select")
	}
	parser := ast.Parser()
	var ast ast.Expression
	if err := parser.ParseString("", expression, &ast); err != nil {
		return nil, err
	}
	var selected []interface{}

	for _, elt := range arr {
		ctx, err := New(elt)
		if err != nil {
			return nil, err
		}
		result, err := ast.Eval(ctx)
		if err != nil {
			return nil, err
		}
		if bool(*result.Bool) {
			selected = append(selected, elt)
		}
	}

	return selected, nil
}
