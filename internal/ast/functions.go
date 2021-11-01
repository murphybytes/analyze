package ast

import (
	"github.com/murphybytes/analyze/context"
	"github.com/murphybytes/analyze/errors"
)
//nolint
type Function struct {
	Name string `@Function`
	Args []*Expression `"(" ( @@ ( "," @@ )* )? ")"`
}

func(f *Function) Eval(ctx Context)(*Value,error){
	switch f.Name {
	case "len":
		return Len(ctx, f.Args)
	}
	return nil, errors.New(errors.SyntaxError, "unknown function %q", f.Name)
}

func Len(ctx Context, args []*Expression)(*Value, error){
	if len(args) != 1  {
		return nil, errors.New(errors.SyntaxError, "wrong number of arguments for len, expected 1, got %d", len(args))
	}
	v, err := args[0].Eval(ctx)
	if err != nil {
		return nil, err
	}
	f := float64(len(v.Array))
	return &Value{
		Number: &f,
	}, nil
}

// In -> in( value, array ) returns true if the value is in the array
func In(ctx context.Context, args []*Expression)(*Value,error){
	panic("not implemented")
}

