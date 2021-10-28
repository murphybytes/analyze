// Package ast contains the components of an abstract syntax tree that represents a predicate Expression.
package ast

type Boolean bool

func(b *Boolean) Capture(values []string) error {
	if len(values) == 0 {
		panic("no values in capture")
	}
	*b = values[0] == "true"
	return nil
}

func(b Boolean) Not() *Boolean {
	v := Boolean(!bool(b))
	return &v
}

func BoolVal(v bool) *Value {
	b := Boolean(v)
	return &Value{
		Bool: &b,
	}
}


type Value struct {
	Number *float64 ` @Number`
	String *string `| @String`
	Bool *Boolean  `| @("true" | "false")`
}


func(v *Value) Eval()(*Value, error){
	return v, nil
}

type UnaryOpValue struct {
	Operator UnaryOperator `@("!")?`
	Value *Value `@@`
}

func (un *UnaryOpValue) Eval()(*Value,error) {
	return un.Operator.Eval(un.Value)
}

type ComparisonOpValue struct {
	Operator ComparisonOperator `@("<")?`
	Value    *UnaryOpValue             `@@`
}

func(c *ComparisonOpValue) Eval()(*Value,error) {
	return c.Value.Eval()
}

type ComparisonOpTerm struct {
	Left *UnaryOpValue `@@`
	Right []*ComparisonOpValue `@@*`
}

func(c *ComparisonOpTerm) Eval()(*Value, error){
	lv, err := c.Left.Eval()
	if err  != nil {
		return nil, err
	}
	for _, exp := range c.Right {
		rv, err := exp.Value.Eval()
		if err != nil {
			return nil, err
		}
		lv, err = exp.Operator.Eval(lv, rv)
		if err != nil {
			return nil, err
		}
	}
	return lv, nil
}


type LogicalOpValue struct {
	Operator LogicalOperator `@("&" "&")`
	Value *ComparisonOpValue `@@`
}


type Expression struct {
	Left *ComparisonOpTerm         `@@`
	Right []*LogicalOpValue `@@*`
}

func (t *Expression) Eval()(*Value, error) {
	lv, err := t.Left.Eval()
	if err != nil {
		return nil, err
	}

	for _, expr := range t.Right {
		rv, err := expr.Value.Eval()
		if err != nil {
			return nil, err
		}
		if lv, err = expr.Operator.Eval(lv, rv); err != nil {
			return nil, err
		}
	}
	return lv, nil
}


