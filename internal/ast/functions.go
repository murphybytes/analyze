package ast

import (
	"fmt"
	"github.com/murphybytes/analyze/errors"
)
//nolint
type Function struct {
	Name string `@Function`
	Args []*Expression `"(" ( @@ ( "," @@ )* )? ")"`
}

func(f *Function) Eval(ctx Context)(*Value,error){
	fn, ok := ctx.Func(f.Name)
	if !ok {
		return nil, errors.New(errors.InvalidFunction, "function %q has not been declared", f.Name)
	}
	var args []interface{}
	for _, expr := range f.Args {
		v, err := expr.Eval(ctx)
		if err != nil {
			return nil, err
		}
		arg, err := valToInterface(v)
		if err != nil {
			fmt.Errorf("%q call is invalid %w", f.Name, err )
		}
		args = append(args, arg )
	}
	result, err :=  fn(args)
	if err != nil {
		return nil, err
	}
	return convertToValue(result)
}

func valToInterface(v *Value)(interface{},error){
	switch {
	case v.String != nil :
		return *v.String, nil
	case v.Bool != nil :
		return *v.Bool, nil
	case v.Number != nil :
		return *v.Number, nil
	case bool(v.NilSet):
		return nil, nil
	case v.Object != nil :
		return v.Object, nil
	case v.Array != nil :
		return v.Array, nil
	}
	return nil, errors.New(errors.InvalidArgumentType, "argument type not supported")
}


