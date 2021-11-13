package expression

import (
	"github.com/murphybytes/analyze/context"
	"github.com/murphybytes/analyze/internal/ast"
	"sync"
)


// Evaluate processes predicate expression with variables populated with data passed in as an
// argument.
func Evaluate(data interface{}, expression string)(bool,error) {
	ctx, err := context.New(data)
	if err != nil {
		return false, err
	}
	var t ast.Expression
	if err := ast.Parser().ParseString("", expression, &t); err != nil {
		return false, err
	}
	result, err := t.Eval(ctx)
	if err != nil {
		return false, err
	}
	return bool(*result.Bool), nil
}

// EvaluateContext takes context with data and optional user defined functions and
// evaluates a predicate expression.
func EvaluateContext(ctx ast.Context, expression string) (bool, error) {

	parser := ast.Parser()
	var t ast.Expression
	if err := parser.ParseString("", expression, &t); err != nil {
		return false, err
	}
	result, err := t.Eval(ctx)
	if err != nil {
		return false, err
	}
	return bool(*result.Bool), nil
}

// PreparedExpression is used to create a thread safe expression that can be used more efficiently because the
// expression tree is parsed only once and can be called repeatedly.
type PreparedExpression struct {
	mut sync.Mutex
	tree ast.Expression
}

// Evaluate evaluate a prepared expression.
func(p *PreparedExpression) Evaluate(ctx ast.Context)(bool, error){
	p.mut.Lock()
	defer p.mut.Unlock()
	result, err := p.tree.Eval(ctx)
	if  err != nil {
		return false, err
	}
	return bool(*result.Bool), nil

}

// Prepare create an expression that you can use repeatedly with different input data.
func Prepare(expression string)(*PreparedExpression,error){
	parser := ast.Parser()
	var result PreparedExpression
	if err := parser.ParseString("", expression, &result.tree); err != nil {
		return nil, err
	}
	return &result, nil
}
