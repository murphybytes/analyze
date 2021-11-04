// Package context defines data that is passed to the Evaluate function and maps to variables.
package context

import (
	"github.com/murphybytes/analyze/errors"
	"github.com/murphybytes/analyze/internal/ast"
	"regexp"
)


type functionTable map[string]ast.UserDefinedFunc
type Option func(*Context) error

type Context struct {
	data      interface{}
	functions functionTable
}

func(c Context) Data() interface{} {
	return c.data
}

func(c Context) Func(name string)(ast.UserDefinedFunc, bool){
	fn, ok := c.functions[name]
	return fn, ok
}

var functionNameMatcher = regexp.MustCompile(`^@[A-Za-z0-9_]\w*`)

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
