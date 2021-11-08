// Package ast contains the components of an abstract syntax tree that represents a predicate Expression.
package ast

type UserDefinedFunc func(a []interface{}) (interface{}, error)

type Context interface {
	Data() interface{}
	Func(string) (UserDefinedFunc, bool)
}

// Value represents data types supported by the predicate expression.
//nolint
type Value struct {
	// Number is the represents floats and integer types in expressions.
	Number *float64 ` @Number`
	// String string literals represented by characters surrounded by double quotes.
	String *string `| @String`
	// Bool true or false keywords
	Bool   *Boolean `| @("true" | "false")`
	NilSet NilFlag  ` | @("nil")`

	// Subexpressions surrounded by parenthesis, innermost subexpression are higher precedence.
	Subexpression *Expression `| "(" @@ ")"`
	// Variables are represented by a leading $ with subelements delimited by dots $foo.bar that are associated
	// with map keys in passed in contexts that are used to pass in data.
	Variable *Variable `| @Variable`
	RegularExpression *RegularExpression `| @RegularExpression`
	// Function
	Function *Function `| @@`
	// These are not set directly in expressions and are used to represent data passed by context.
	Object map[string]interface{}
	Array  []interface{}
}

func (v Value) IsNil() bool {
	if v.Number != nil {
		return false
	}
	if v.String != nil {
		return false
	}
	if v.Bool != nil {
		return false
	}
	if v.Object != nil {
		return false
	}
	if v.Array != nil {
		return false
	}
	return true
}

func (v *Value) Eval(ctx Context) (*Value, error) {
	// TODO: Fix this so it only gets evaluated once
	if v.Subexpression != nil {
		return v.Subexpression.Eval(ctx)
	}
	if v.Variable != nil {
		return v.Variable.Eval(ctx)
	}
	if v.RegularExpression != nil {
		return v.RegularExpression.Eval(ctx)
	}
	if v.Function != nil {
		return v.Function.Eval(ctx)
	}
	return v, nil
}

//nolint
type UnaryOpValue struct {
	Operator *Operator `@("!")?`
	Value    *Value    `@@`
}

func (un *UnaryOpValue) Eval(ctx Context) (*Value, error) {
	v, err := un.Value.Eval(ctx)
	if err != nil {
		return nil, err
	}
	if un.Operator != nil {
		return un.Operator.Eval(ctx, v)
	}
	return v, nil
}

//nolint
type ComparisonOpValue struct {
	Operator Operator      `@("<" | "<=" | "==" | "!=" | ">" | ">=" )?`
	Value    *UnaryOpValue `@@`
}

func (c *ComparisonOpValue) Eval(ctx Context) (*Value, error) {
	return c.Value.Eval(ctx)
}

//nolint
type ComparisonOpTerm struct {
	Left  *UnaryOpValue        `@@`
	Right []*ComparisonOpValue `@@*`
}

func (c *ComparisonOpTerm) Eval(ctx Context) (*Value, error) {
	lv, err := c.Left.Eval(ctx)
	if err != nil {
		return nil, err
	}
	for _, exp := range c.Right {
		rv, err := exp.Value.Eval(ctx)
		if err != nil {
			return nil, err
		}
		lv, err = exp.Operator.Eval(ctx, lv, rv)
		if err != nil {
			return nil, err
		}
	}
	return lv, nil
}

//nolint
type LogicalOpValue struct {
	Operator Operator          `@("&&" | "||")`
	Value    *ComparisonOpTerm `@@`
}

//nolint
type Expression struct {
	Left  *ComparisonOpTerm `@@`
	Right []*LogicalOpValue `@@*`
}

func (t *Expression) Eval(ctx Context) (*Value, error) {
	lv, err := t.Left.Eval(ctx)
	if err != nil {
		return nil, err
	}

	for _, expr := range t.Right {
		rv, err := expr.Value.Eval(ctx)
		if err != nil {
			return nil, err
		}
		if lv, err = expr.Operator.Eval(ctx, lv, rv); err != nil {
			return nil, err
		}
	}
	return lv, nil
}
