package expression

import (
	"testing"

	"github.com/murphybytes/analyze/context"
	"github.com/stretchr/testify/require"
)

func TestEval(t *testing.T) {
	tt := []struct {
		name       string
		expression string
		expected   bool
		wantErr    bool
		context    context.Context
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
			context: map[string]interface{}{
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
			context: map[string]interface{}{
				"foo": map[string]interface{}{
					"value": 5,
				},
			},
			expected: true,
		},
		{
			name: "string equality",
			expression: `"one" == "one"`,
			expected: true,
		},
		{
			name: "string equality false ",
			expression: `"one" == "two"`,
			expected: false,
		},
		{
			name: "string not equals ",
			expression: `"one" != "two"`,
			expected: true,
		},
		{
			name: "binary or",
			expression: `( 2 < 1 ) || ( 5 < 6)`,
			expected: true,
		},
		{
			name: "greater than",
			expression: `2 > 1`,
			expected: true,
		},
		{
			name: "greater than or equal to",
			expression: `3 >= 4 && 3 >= 3`,
			expected: true,
		},
		{
			name:       "index into array",
			expression: `1 < $foo[1].bar`,
			context: map[string]interface{}{
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
			expression: `len( $foo ) < 10`,
			context: map[string]interface{}{
				"foo": []interface{}{1, 2, 3},
			},
			expected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Evaluate(tc.expression, tc.context)
			if tc.wantErr {
				require.NotNil(t, err)
				return
			}

			require.Nil(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}

}
