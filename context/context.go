// Package context defines data and optional user functions that are used to evaluate expressions.
package context

import (
	"github.com/murphybytes/analyze/errors"
	"github.com/murphybytes/analyze/internal/ast"
	"regexp"
)


type functionTable map[string]ast.UserDefinedFunc
// Option optional function for New.
type Option func(*Context) error

// Context contains data used to evaluate expression.
type Context struct {
	data      interface{}
	functions functionTable
}

// Data returns data that maps to variables defined in expressions.
func(c Context) Data() interface{} {
	return c.data
}

// Func returns a named function.
func(c Context) Func(name string)(ast.UserDefinedFunc, bool){
	fn, ok := c.functions[name]
	return fn, ok
}

var functionNameMatcher = regexp.MustCompile(`^@[A-Za-z0-9_]\w*`)

// Func pass a user defined function to a new context.  The name for the function must be prefaced by '@' for
// example @abs.  You would then use the function in an expression thus: @abs(-10) > 9
func Func(name string, definedFunc ast.UserDefinedFunc) Option {
	return func(ctx *Context) error {
		if !functionNameMatcher.MatchString(name) {
			return errors.New(errors.InvalidFunction, "%q is not a valid function name", name)
		}
		if _, ok := ctx.functions[name]; ok {
			return errors.New(errors.DuplicateFunction, "function name %q already in use", name)
		}
		ctx.functions[name] = definedFunc
		return nil
	}
}

// New creates a new context with data that can be referenced in variables in expressions.  User defined functions
// can optionally be passed as well.
func New(data interface{}, options ...Option) (*Context, error) {
	if err := validate(data); err != nil {
		return nil ,err
	}

	ctx := Context{
		data:      data,
	}

	// builtin functions
	ctx.functions = functionTable{
		"@len": _len,
		"@select": _select,
		"@in": _in ,
		"@array": _array,
		"@has": _has,
		"@match": _match,
	}

	for _, opt := range options {
		if err := opt(&ctx); err != nil {
			return nil, err
		}
	}

	return &ctx, nil
}

func validate(data interface{}) error {
	switch t := data.(type) {
	case int:
		return nil
	case *int:
		return nil
	case float64:
		return nil
	case *float64:
		return nil
	case bool:
		return nil
	case *bool:
		return nil
	case string:
		return nil
	case *string:
		return nil
	case nil:
		return nil
	case []interface{}:
		for _, elt := range t {
			if err := validate(elt); err != nil {
				return err
			}
		}
		return nil
	case map[string]interface{}:
		for _, val := range t {
			if err := validate(val); err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New(errors.UnsupportedType, "input data validation failed because of unsupported type %T", data)
}
