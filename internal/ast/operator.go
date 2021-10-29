package ast

import (
	"strings"
	"sync"

	"github.com/murphybytes/dsl/context"
)

type ComparisonOperator int

const (
	UnknownComparison ComparisonOperator = iota
	OpLessThan
	OpLessThanEqual
)

var comparisonOps = map[string]ComparisonOperator{
	"<":  OpLessThan,
	"<=": OpLessThanEqual,
}

var comparisonOpMut sync.Mutex

func getComparisonOpID(op string) (ComparisonOperator, error) {
	comparisonOpMut.Lock()
	defer comparisonOpMut.Unlock()
	id, ok := comparisonOps[op]
	if !ok {
		return UnknownComparison, NewUnsupportedOperatorError(op)
	}
	return id, nil
}

func (o *ComparisonOperator) Capture(s []string) (err error) {
	key := strings.Join(s, "")
	*o, err = getComparisonOpID(key)
	return
}

func (o ComparisonOperator) Eval(_ context.Context, l, r *Value) (*Value, error) {
	if l.Number != nil && r.Number != nil {
		return compareNumber(o, *l.Number, *r.Number), nil
	}
	if l.String != nil && r.String != nil {
		return compareString(o, *l.String, *r.String), nil
	}
	return nil, TypeMismatchError(l, r)
}

func compareNumber(o ComparisonOperator, l, r float64) *Value {
	switch o {
	case OpLessThan:
		return BoolVal(l < r)
	case OpLessThanEqual:
		return BoolVal(l <= r)
	}
	panic("unsupported operator")
}

func compareString(o ComparisonOperator, l, r string) *Value {
	switch o {
	case OpLessThan:
		return BoolVal(l < r)
	case OpLessThanEqual:
		return BoolVal(l <= r)
	}
	panic("unsupported operator")
}

type UnaryOperator int

const (
	UnassignedUnary UnaryOperator = iota
	UnaryNot
)

var unaryOps = map[string]UnaryOperator{
	"!": UnaryNot,
}

var unaryOpMut sync.Mutex

func getUnaryOpID(op string) (UnaryOperator, error) {
	unaryOpMut.Lock()
	defer unaryOpMut.Unlock()
	id, ok := unaryOps[op]
	if !ok {
		return UnassignedUnary, NewUnsupportedOperatorError(op)
	}
	return id, nil
}

func (o *UnaryOperator) Capture(s []string) (err error) {
	key := strings.Join(s, "")
	*o, err = getUnaryOpID(key)
	return
}

func (o UnaryOperator) Eval(_ context.Context, v *Value) (*Value, error) {
	switch o {
	case UnaryNot:
		if v.Bool == nil {
			return nil, NewSyntaxError("nil bool not allowed")
		}
		return &Value{
			Bool: v.Bool.Not(),
		}, nil
	case UnassignedUnary:
		return v, nil

	}
	panic("unknown operator type")
}

type LogicalOperator int

const (
	UnassignedLogical LogicalOperator = iota
	LogicalAnd
	LogicalOr
)

var logicalOps = map[string]LogicalOperator{
	"&&": LogicalAnd,
	"||": LogicalOr,
}

var logicalOpMut sync.Mutex

func getLogicalOpID(op string) (LogicalOperator, error) {
	logicalOpMut.Lock()
	defer logicalOpMut.Unlock()
	id, ok := logicalOps[op]
	if !ok {
		return UnassignedLogical, NewUnsupportedOperatorError(op)
	}
	return id, nil
}

func (o *LogicalOperator) Capture(s []string) (err error) {
	key := strings.Join(s, "")
	*o, err = getLogicalOpID(key)
	return
}

func (o LogicalOperator) Eval(_ context.Context, l, r *Value) (*Value, error) {
	switch o {
	case LogicalAnd:
		return and(l, r)
	case LogicalOr:
		return or(l, r)
	}
	return nil, NewSyntaxError("unsupported logical operator")
}

func and(l, r *Value) (*Value, error) {
	if l.Bool != nil && r.Bool != nil {
		lv := bool(*l.Bool)
		rv := bool(*r.Bool)
		return BoolVal(lv && rv), nil
	}
	return nil, NewSyntaxError("nil bool not allowed")
}

func or(l, r *Value) (*Value, error) {
	if l.Bool != nil && r.Bool != nil {
		lv := bool(*l.Bool)
		rv := bool(*r.Bool)
		return BoolVal(lv || rv), nil
	}
	return nil, NewSyntaxError("nil bool not allowed")
}
