package expression

import (
	"fmt"
	"github.com/murphybytes/analyze/internal/ast"
	"testing"

	"github.com/murphybytes/analyze/context"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestEval(t *testing.T) {
	tt := []struct {
		name       string
		expression string
		expected   bool
		wantErr    bool
		data       interface{}
	}{
		{
			name:       "int less than",
			expression: "1 < 2",
			expected:   true,
		},
		{
			name:       "int less than",
			expression: "3 < 2",
			expected:   false,
		},
		{
			name:       "negative int less than",
			expression: "-3 < 2",
			expected:   true,
		},
		{
			name:       "extra whitespace",
			expression: "  -3 <   2",
			expected:   true,
		},
		{
			name:       "string less than",
			expression: `"aardvark" < "arhus"`,
			expected:   true,
		},
		{
			name:       "not",
			expression: "!true",
			expected:   false,
		},
		{
			name:       "and",
			expression: "1 < 2 && 8 < 9",
			expected:   true,
		},
		{
			name:       "subexpression",
			expression: "(1 < 2) && (9 < 8)",
			expected:   false,
		},
		{
			name:       "less than equal to",
			expression: "1 <= 3",
			expected:   true,
		},
		{
			name:       "not subexpression",
			expression: "!(1 < 2)",
			expected:   false,
		},
		{
			name:       "nested subexpression",
			expression: "(1 < 2) && ((3 < 4) && (5 < 4))",
			expected:   false,
		},
		{
			name:       "simple variable",
			expression: "1 < $foo.value",
			data: map[string]interface{}{
				"foo": map[string]interface{}{
					"value":     4,
					"something": "xxx",
				},
			},
			expected: true,
		},
		{
			name:       "index into object",
			expression: `1 < $foo["value"]`,
			data: map[string]interface{}{
				"foo": map[string]interface{}{
					"value": 5,
				},
			},
			expected: true,
		},
		{
			name:       "string equality",
			expression: `"one" == "one"`,
			expected:   true,
		},
		{
			name:       "string equality false ",
			expression: `"one" == "two"`,
			expected:   false,
		},
		{
			name:       "string not equals ",
			expression: `"one" != "two"`,
			expected:   true,
		},
		{
			name:       "binary or",
			expression: `( 2 < 1 ) || ( 5 < 6)`,
			expected:   true,
		},
		{
			name:       "greater than",
			expression: `2 > 1`,
			expected:   true,
		},
		{
			name:       "greater than or equal to",
			expression: `3 >= 4 && 3 >= 3`,
			expected:   true,
		},
		{
			name:       "index into array",
			expression: `1 < $foo[1].bar`,
			data: map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"bar": 0,
					},
					map[string]interface{}{
						"bar": 5,
					},
				},
			},
			expected: true,
		},
		{
			name:       "simple function",
			expression: `@len( $foo ) < 10`,
			data: map[string]interface{}{
				"foo": []interface{}{1, 2, 3},
			},
			expected: true,
		},
		{
			name:       "array root",
			expression: "$[1] == 3",
			data: []interface{}{
				5,
				3,
			},
			expected: true,
		},
		{
			name:       "object root",
			expression: `$["field"] == 3 && $["another-field"] < 5`,
			data: map[string]interface{}{
				"field":         3,
				"another-field": 4,
			},
			expected: true,
		},
		{
			name:       "scalar root",
			expression: "$ < 6",
			data:       5,
			expected:   true,
		},
		{
			name:       "nil type",
			expression: `$ != nil`,
			data:       "xxx",
			expected:   true,
		},
		{
			name:       "nil type",
			expression: `$ != nil`,
			data:       nil,
			expected:   false,
		},
		{
			// avoid type mismatch because $foo.bar == 3 never is evaluated
			name:       "short circuit and",
			expression: `false && $foo.bar == 3`,
			data: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": nil,
				},
			},
			expected: false,
		},
		{
			// avoid type mismatch because $foo.bar == 3 never is evaluated
			name:       "short circuit or",
			expression: `true || $foo.bar == 3`,
			data: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": nil,
				},
			},
			expected: true,
		},
		{
			name:       "type mismatch unary",
			expression: "!3",
			wantErr:    true,
		},
		{
			name:       "select test",
			expression: `@len( @select( $arr, "$elt == 3" ) ) > 2`,
			data:       []interface{}{3, 1, 2, 3, 3},
			expected:   true,
		},
		{
			name:       "in func test",
			expression: `@in( @array(1, 2, 3), 2)`,
			expected:   true,
		},
		{
			name:       "in with string arr var",
			expression: `@in( $arr, "foo")`,
			data: []interface{}{
				"zip",
				"foo",
				"bazz",
			},
			expected: true,
		},
		{
			name:       "has function",
			expression: `@has($bar, "foo") && $bar.foo == 3`,
			data: map[string]interface{}{
				"bar": map[string]interface{}{
					"foo": 3,
				},
			},
			expected: true,
		},
		{
			name:       "root object dotted reference",
			expression: `3 < $foo`,
			data: map[string]interface{}{
				"foo": 4,
			},
			expected: true,
		},
		{
			name: "match",
			expression: `@match("10.10.10.10", /^([0-9]{1,3}\.){3}[0-9]{1,3}$/)`,
			expected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx, err := context.New(tc.data)
			require.Nil(t, err)
			actual, err := Evaluate(ctx, tc.expression)
			if tc.wantErr {
				require.NotNil(t, err)
				return
			}

			require.Nil(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}

}

func TestUserDefinedFunctions(t *testing.T) {
	tt := []struct {
		name       string
		fns        map[string]ast.UserDefinedFunc
		expression string
		expected   bool
		wantErr    bool
		data       interface{}
	}{
		{
			name: "simple function",
			fns: map[string]ast.UserDefinedFunc{
				"@tester": func(a []interface{}) (interface{}, error) {
					return a[0], nil
				},
			},
			expression: `@tester($foo) == $foo`,
			data:       2,
			expected:   true,
		},
		{
			name: "string arg",
			fns: map[string]ast.UserDefinedFunc{
				"@tester": func(a []interface{}) (interface{}, error) {
					return a[0], nil
				},
			},
			expression: `@tester(2 , "$foo == 3") == 2`,
			data:       2,
			expected:   true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var opts []context.Option
			for name, fn := range tc.fns {
				opts = append(opts, context.Func(name, fn))
			}
			ctx, err := context.New(tc.data, opts...)
			require.Nil(t, err)
			actual, err := Evaluate(ctx, tc.expression)
			if tc.wantErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, actual, tc.expected)

		})
	}
}

func TestPreparedExpression(t *testing.T) {
	tt := []struct {
		data     interface{}
		expected bool
	}{
		{

			data:     2,
			expected: true,
		},
		{
			data:     5,
			expected: false,
		},
		{
			data:     7,
			expected: false,
		},
	}

	expression, err := Prepare("$ < 3")
	require.Nil(t, err)
	g := new(errgroup.Group)

	for i := 0; i < len(tt); i++ {

		g.Go(func(d interface{}, expected bool, j int) func() error {
			return func() error {
				ctx, err := context.New(d)
				if err != nil {
					return err
				}
				actual, err := expression.Evaluate(ctx)
				if err != nil {
					return err
				}
				if actual != expected {
					return fmt.Errorf("test %d failed", j)
				}
				return nil
			}
		}(tt[i].data, tt[i].expected, i))
	}

	require.Nil(t, g.Wait())

}
