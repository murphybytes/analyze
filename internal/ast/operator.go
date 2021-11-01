package ast

import (
	"github.com/murphybytes/analyze/errors"
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
	OpGreaterThan
	OpGreaterThanOrEqualTo
	OpAnd
	OpOr
	OpEqualTo
	OpNotEqualTo
)

func (o *Operator) Capture(s []string) error {
	key := strings.Join(s, "")
	idMap := map[string]Operator{
		"<":  OpLessThan,
		"<=": OpLessThanEqual,
		"!":  OpUnaryNot,
		"&&": OpAnd ,
		"||": OpOr,
		"==": OpEqualTo,
		"!=": OpNotEqualTo,
		">": OpGreaterThan,
		">=": OpGreaterThanOrEqualTo,
	}
	var ok bool
	if *o, ok = idMap[key]; !ok {
		return errors.New(errors.UnsupportedOperator, "unsupported operator %q", key )
	}
	return nil
}

func (o *Operator) Eval(ctx Context, values ...*Value) (*Value, error) {
	fnMap := map[Operator]func(Context, ...*Value) (*Value, error){
		OpUnaryNot: func(ctx Context, values ...*Value) (*Value, error) {
			return mapUnary(values, func(l *Value) (*Value, error) {
				if !hasNilBools(l) {
					return BoolVal(!bool(*l.Bool)), nil
				}
				return nil, errors.New(errors.SyntaxError, "type mismatch")
			})
		},
		OpLessThan: func(ctx Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				if !hasNilNumbers(l, r) {
					return BoolVal(*l.Number < *r.Number), nil
				}
				if !hasNilStrings(l, r) {
					return BoolVal(*l.String < *r.String), nil
				}
				return nil, errors.New(errors.SyntaxError, "type mismatch")
			})
		},
		OpLessThanEqual: func(ctx Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				if !hasNilNumbers(l, r) {
					return BoolVal(*l.Number <= *r.Number), nil
				}
				if !hasNilStrings(l, r) {
					return BoolVal(*l.String <= *r.String), nil
				}
				return nil, errors.New(errors.SyntaxError, "type mismatch")
			})
		},
		OpGreaterThan: func(ctx Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				if !hasNilNumbers(l, r) {
					return BoolVal(*l.Number > *r.Number), nil
				}
				if !hasNilStrings(l, r) {
					return BoolVal(*l.String > *r.String), nil
				}
				return nil, errors.New(errors.SyntaxError, "type mismatch")
			})
		},
		OpGreaterThanOrEqualTo: func(ctx Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				if !hasNilNumbers(l, r) {
					return BoolVal(*l.Number <= *r.Number), nil
				}
				if !hasNilStrings(l, r) {
					return BoolVal(*l.String <= *r.String), nil
				}
				return nil, errors.New(errors.SyntaxError, "type mismatch")
			})
		},
		OpEqualTo: func(ctx Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				if !hasNilNumbers(l, r) {
					return BoolVal(*l.Number == *r.Number), nil
				}
				if !hasNilStrings(l, r) {
					return BoolVal(*l.String == *r.String), nil
				}
				if !hasNilBools(l,r) {
					return BoolVal(bool(*l.Bool) == bool(*r.Bool) ), nil
				}
				if l.NilSet || r.NilSet {
					return BoolVal(l.IsNil() == r.IsNil()), nil
				}
				return nil, errors.New(errors.SyntaxError, "type mismatch")
			})
		},
		OpNotEqualTo: func(ctx Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				if !hasNilNumbers(l, r) {
					return BoolVal(*l.Number != *r.Number), nil
				}
				if !hasNilStrings(l, r) {
					return BoolVal(*l.String != *r.String), nil
				}
				if !hasNilBools(l,r) {
					return BoolVal(bool(*l.Bool) != bool(*r.Bool) ), nil
				}
				if l.NilSet || r.NilSet {
					return BoolVal(l.IsNil() != r.IsNil()), nil
				}
				return nil, errors.New(errors.SyntaxError, "type mismatch")
			})
		},
		OpAnd: func(ctx Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				// short circuit eval, if lval is false, ignore rval and return false
				if l.Bool != nil && !bool(*l.Bool) {
					return BoolVal(false), nil
				}
				if !hasNilBools(l, r) {
					return BoolVal(bool(*l.Bool) && bool(*r.Bool)), nil
				}
				return nil, errors.New(errors.SyntaxError, "type mismatch")
			})
		},
		OpOr: func(ctx Context, values ...*Value) (*Value, error) {
			return mapBinary(values, func(l, r *Value) (*Value, error) {
				// short circuit eval, if lval is true, ignore rval and return true
				if l.Bool != nil && bool(*l.Bool) {
					return BoolVal(true), nil
				}
				if !hasNilBools(l, r) {
					return BoolVal(bool(*l.Bool) || bool(*r.Bool)), nil
				}
				return nil, errors.New(errors.SyntaxError, "type mismatch")
			})
		},
	}

	// find the function that maps to a particular operator, evaluate the variables it uses to resolve variables,
	// functions, and subexpressions to types, and execute
	if fn, ok := fnMap[*o]; ok {
		return fn(ctx, values...)
	}

	return nil, errors.New(errors.SyntaxError, "eval called on uninitialzed operator")
}

func mapBinary(vals []*Value, fn func(l, r *Value) (*Value, error)) (*Value, error) {
	if len(vals) != 2 {
		return nil, errors.New(errors.SyntaxError, "expected 2 arguments got %d", len(vals))
	}
	return fn(vals[0], vals[1])
}

func mapUnary(vals []*Value, fn func(v *Value) (*Value, error)) (*Value, error) {
	ln := len(vals)
	if ln != 1 {
		return nil, errors.New(errors.SyntaxError, "expected 1 argument got %d", ln)
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
