package ast

import (
	"strings"
)

type NilFlag bool

func(n *NilFlag) Capture(values []string) error {
	*n = true
	return nil
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	if len(values) == 0 {
		panic("no values in capture")
	}
	*b = strings.Join(values, "") == "true"
	return nil
}

func (b Boolean) Not() *Boolean {
	v := Boolean(!bool(b))
	return &v
}

func BoolVal(v bool) *Value {
	b := Boolean(v)
	return &Value{
		Bool: &b,
	}
}
