package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"", nil},
		{"a", "a"},
		{"foo", "foo"},
		{"()", []any{}},
		{"( )", []any{}},
		{"( ( ) )", []any{[]any{}}},
		{"(()(()))", []any{[]any{}, []any{[]any{}}}},
		{`(+ 1 2)`, []any{"+", 1, 2}},
		{`(* (+ 1 2) (- 3 5))`, []any{"*", []any{"+", 1, 2}, []any{"-", 3, 5}}},
		{
			`
		(define (fact n)
		  (if (= n 0)
			1
			(* n (fact (- n 1)))))`,
			[]any{
				"define",
				[]any{"fact", "n"},
				[]any{
					"if",
					[]any{"=", "n", 0},
					1,
					[]any{"*", "n", []any{"fact", []any{"-", "n", 1}}},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := Parse(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
