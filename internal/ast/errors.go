package ast

import "fmt"

type ErrType int
const (
	TypeMismatch ErrType = iota
	UnsupportedType
	UnsupportedOperator
	SyntaxError
)

type Error interface {
	error
	Type() ErrType
}

type ErrAst struct {
	typ ErrType
	msg string
}

func(ea ErrAst) Error() string {
	return ea.msg
}

func(ea ErrAst) Type() ErrType {
	return ea.typ
}

// TypeMismatchError indicates types of the arguments don't match when they should.
func TypeMismatchError(l,r *Value) error {
	return &ErrAst{
		msg: fmt.Sprintf("type mismatch lval %#v rval %#v", l, r ),
		typ: TypeMismatch,
	}
}

// NewUnsupportedOperatorError indicates that an unsupported operator is being used.
func NewUnsupportedOperatorError(s string) error {
	return &ErrAst{
		msg: fmt.Sprintf("unsupported operator %q", s),
		typ: UnsupportedOperator,
	}
}

func NewSyntaxError() error {
	return &ErrAst{
		msg: fmt.Sprintf("syntax error"),
		typ: SyntaxError,
	}
}
