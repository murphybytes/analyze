package predicate

import (
	"github.com/stretchr/testify/require"
	"testing"
)


func TestEval(t *testing.T) {
	tt := []struct {
		name string
		expression string
		expected bool
		wantErr bool
	} {
		{
			name: "int less than",
			expression: "1 < 2",
			expected: true,
		},
		{
			name: "int less than",
			expression: "3 < 2",
			expected: false,
		},
		{
			name: "negative int less than",
			expression: "-3 < 2",
			expected: true,
		},
		{
			name: "extra whitespace",
			expression: "  -3 <   2",
			expected: true,
		},
		{
			name: "string less than",
			expression: `"aardvark" < "arhus"`,
			expected: true,
		},
		{
			name: "not",
			expression: "!true",
			expected: false,
		},
		{
			name: "and",
			expression: "true && true",
			expected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T){
			actual, err := Evaluate(tc.expression, nil)
			if tc.wantErr {
				require.NotNil(t, err)
				return
			}

			require.Nil(t, err)
			require.Equal(t, tc.expected, actual )
		})
	}

}