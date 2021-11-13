package context

import (
	"github.com/murphybytes/analyze/errors"
	"github.com/murphybytes/analyze/internal/ast"
	"regexp"
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

// @array(val1, val2, .... valN) converts a list of values to an array
func _array(args []interface{}) (interface{}, error) {
	return args, nil
}

// @in(arr, val) if value is in array returns true
func _in(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New(errors.SyntaxError, "wrong number of arguments for in, expected 2, got %d", len(args))
	}
	switch t := args[0].(type) {
	case []interface{}:
		for _, elt := range t {
			match, err := compare(elt, args[1])
			if err != nil {
				return nil, err
			}
			if match {
				return true, nil
			}
		}
		return false, nil
	}
	return nil, errors.New(errors.TypeMismatch, "expected array for first argument for in function")
}

// @has(obj, fieldname) returns true if field exists in object
func _has(args []interface{})(interface{},error){
	if len(args) != 2 {
		return nil, errors.New(errors.SyntaxError, "wrong number of arguments for has function expected 2, got %d", len(args))
	}
	switch t := args[0].(type) {
	case map[string]interface{}:
		key, ok := args[1].(string)
		if !ok {
			return nil, errors.New(errors.TypeMismatch, "expect string for second has function argument")
		}
		_, present := t[key]
		return present, nil
	}
	return nil, errors.New(errors.TypeMismatch, "expected object for first argument of has function")
}

// @match(string, regex-string) returns true if string matches regex-string of the form /regular expression/
func _match(args []interface{})(interface{},error){
	if len(args) != 2 {
		return nil, errors.New(errors.SyntaxError, "match expects 2 arguments")
	}
	val, ok := args[0].(string)
	if !ok {
		return nil, errors.New(errors.TypeMismatch, "match expects string argument")
	}
	exp, ok := args[1].(string)
	if !ok {
		return nil, errors.New(errors.TypeMismatch, "match expects string argument 2")
	}
	regex, err := regexp.Compile(exp)
	if err != nil {
		return nil, err
	}
	return regex.MatchString(val), nil
}

func compare(l, r interface{}) (bool, error) {
	switch t := l.(type) {
	case int:
		if rt, ok := r.(int); ok {
			return t == rt, nil
		}
		return false, errors.New(errors.TypeMismatch, "type mismatch")
	case *int:
		if rt, ok := r.(*int); ok {
			return *t == *rt, nil
		}
		return false, errors.New(errors.TypeMismatch, "type mismatch")
	case float64:
		if rt, ok := r.(float64); ok {
			return t == rt, nil
		}
		return false, errors.New(errors.TypeMismatch, "type mismatch")
	case *float64:
		if rt, ok := r.(*float64); ok {
			return *t == *rt, nil
		}
		return false, errors.New(errors.TypeMismatch, "type mismatch")
	case string:
		if rt, ok := r.(string); ok {
			return t == rt, nil
		}
		return false, errors.New(errors.TypeMismatch, "type mismatch")
	case *string:
		if rt, ok := r.(*string); ok {
			return *t == *rt, nil
		}
		return false, errors.New(errors.TypeMismatch, "type mismatch")
	}
	return false, errors.New(errors.UnsupportedType, "unsupported type for in function")
}
