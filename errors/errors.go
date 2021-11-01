// Package errors communicates common problems to package consumers.
package errors

import "fmt"

type ErrType int

const (
	TypeMismatch ErrType = iota
	UnsupportedType
	UnsupportedOperator
	SyntaxError
	MissingKey
	IndexOutOfRange
	UnexpectedError
	DuplicateFunction
)

type Error interface {
	error
	Type() ErrType
}

type ErrAst struct {
	typ ErrType
	msg string
}

func (ea ErrAst) Error() string {
	return ea.msg
}

func (ea ErrAst) Type() ErrType {
	return ea.typ
}

func New(t ErrType, format string, v ...interface{}) error {
	return &ErrAst{
		msg: fmt.Sprintf(format, v...),
		typ: t,
	}
}

func NewUnexpectedError(form string, v ...interface{}) error {
	return &ErrAst{
		msg: fmt.Sprintf("unexpected error: %s", fmt.Sprintf(form, v...)),
		typ: UnexpectedError,
	}
}

// UnsupportedTypeError is raised when data passed in the context doesn't map
// to a data type we support
func UnsupportedTypeError(t interface{}) error {
	return &ErrAst{
		msg: fmt.Sprintf("unsupported type %T", t),
		typ: UnsupportedType,
	}
}

// MissingKeyError is raised when the data passed in as context doesn't map to an expression variable.
func MissingKeyError(key string) error {
	return &ErrAst{
		msg: fmt.Sprintf("variable segment %q does not map to context data", key),
		typ: MissingKey,
	}
}

// IndexOutOfRangeError is raised when the data passed in as context doesn't map to an expression variable.
func IndexOutOfRangeError(index string) error {
	return &ErrAst{
		msg: fmt.Sprintf("variable segment %q does not map to context data", index),
		typ: IndexOutOfRange,
	}
}


// NewUnsupportedOperatorError indicates that an unsupported operator is being used.
func NewUnsupportedOperatorError(s string) error {
	return &ErrAst{
		msg: fmt.Sprintf("unsupported operator %q", s),
		typ: UnsupportedOperator,
	}
}

func NewSyntaxError(formt string, v ...interface{}) error {
	return &ErrAst{
		msg: fmt.Sprintf("syntax error: %s", fmt.Sprintf(formt, v...)),
		typ: SyntaxError,
	}
}
