package ast

import "github.com/murphybytes/analyze/context"
//nolint
type Function struct {
	Name string `@Function`
	Args []*Expression `"(" ( @@ ( "," @@ )* )? ")"`
}

func(f *Function) Eval(ctx context.Context)(*Value,error){
	switch f.Name {
	case "len":
		return Len(ctx, f.Args)
	}
	return nil, NewSyntaxError("unknown function %q", f.Name)
}

func Len(ctx context.Context, args []*Expression)(*Value, error){
	if len(args) != 1  {
		return nil, NewSyntaxError("wrong number of arguments for len, expected 1, got %d", len(args))
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

