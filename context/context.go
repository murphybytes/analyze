// Package context defines data that is passed to the Evaluate function and maps to variables.
package context

import (
	"github.com/murphybytes/analyze/errors"
)

type UserDefinedFunc func(a ...interface{}) (interface{}, error)
type functionTable map[string]UserDefinedFunc
type Option func(*Context) error

type Context struct {
	data      interface{}
	functions functionTable
}

func(c Context) Data() interface{} {
	return c.data
}

func(c Context) Func(name string)(UserDefinedFunc, bool){
	fn, ok := c.functions[name]
	return fn, ok
}

func Func(name string, definedFunc UserDefinedFunc) Option {
	return func(ctx *Context) error {
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
		functions: make(functionTable),
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
