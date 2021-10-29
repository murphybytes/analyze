package ast

import (
	"github.com/murphybytes/dsl/context"
	"strings"
)

// TODO: combine operator into one thing, the precedence is controlled in the ast tags so there is no reason to
// TODO: segregate them here .

type Operator int

const (
	OpUnknown Operator = iota
	OpUnaryNot
	OpLessThan
	OpLessThanEqual
	OpAnd
	OpOr
)

var 	idMap = map[string]Operator{
	"<":  OpLessThan,
	"<=": OpLessThanEqual,
	"!":  OpUnaryNot,
	"&&": OpAnd ,
	"||": OpOr,
}

func getOpID(op string)(Operator, error){
	id, ok := idMap[op]
	if !ok {
		return OpUnknown, NewUnsupportedOperatorError(op)
	}
	return id, nil
}

func (o *Operator) Capture(s []string)(err error) {
	key := strings.Join(s, "")
	*o, err = getOpID(key)
	return
}

func (o *Operator) Eval(ctx context.Context, values ...*Value) (*Value, error) {
	fnMap := map[Operator]func(context.Context, ...*Value) (*Value, error){
		OpUnaryNot: func(ctx context.Context, values ...*Value) (*Value, error) {
			return mapUnary(values, func(l *Value) (*Value, error) {
				if !hasNilBools(l) {
					return BoolVal(!bool(*l.Bool)), nil
				}
				return nil, NewSyntaxError("type mismatch")
			})
		},
		OpLessThan: func(ctx context.Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				if !hasNilNumbers(l, r) {
					return BoolVal(*l.Number < *r.Number), nil
				}
				if !hasNilStrings(l, r) {
					return BoolVal(*l.String < *r.String), nil
				}
				return nil, NewSyntaxError("type mismatch")
			})
		},
		OpLessThanEqual: func(ctx context.Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				if !hasNilNumbers(l, r) {
					return BoolVal(*l.Number <= *r.Number), nil
				}
				if !hasNilStrings(l, r) {
					return BoolVal(*l.String <= *r.String), nil
				}
				return nil, NewSyntaxError("type mismatch")
			})
		},
		OpAnd: func(ctx context.Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				if !hasNilBools(l, r) {
					return BoolVal(bool(*l.Bool) && bool(*r.Bool)), nil
				}
				return nil, NewSyntaxError("type mismatch")
			})
		},
		OpOr: func(ctx context.Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				if !hasNilBools(l, r) {
					return BoolVal(bool(*l.Bool) || bool(*r.Bool)), nil
				}
				return nil, NewSyntaxError("type mismatch")
			})
		},
	}

	// find the function that maps to a particular operator, evaluate the variables it uses to resolve variables,
	// functions, and subexpressions to types, and execute
	if fn, ok := fnMap[*o]; ok {
		return fn(ctx, values...)
	}

	return nil, NewSyntaxError("eval called on uninitialzed operator")
}

func mapBinary(vals []*Value, fn func(l, r *Value) (*Value, error)) (*Value, error) {
	if len(vals) != 2 {
		return nil, NewSyntaxError("expected 2 arguments got %d", len(vals))
	}
	return fn(vals[0], vals[1])
}

func mapUnary(vals []*Value, fn func(v *Value) (*Value, error)) (*Value, error) {
	ln := len(vals)
	if ln != 1 {
		return nil, NewSyntaxError("expected 1 argument got %d", ln)
	}
	return fn(vals[0])
}

func hasNilBools(vals ...*Value) bool {
	for _, v := range vals {
		if v.Bool == nil {
			return true
		}
	}
	return false
}

func hasNilNumbers(vals ...*Value) bool {
	for _, v := range vals {
		if v.Number == nil {
			return true
		}
	}
	return false
}

func hasNilStrings(vals ...*Value) bool {
	for _, v := range vals {
		if v.String == nil {
			return true
		}
	}
	return false
}
