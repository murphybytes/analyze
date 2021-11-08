package ast

import "strings"

type RegularExpression string

func(r *RegularExpression) Capture(input []string) error {
	*r = RegularExpression(strings.Trim(strings.Join(input, ""), "/"))
	return nil
}

func(r *RegularExpression) Eval(_ Context)(*Value, error){
	s := string(*r)
	return &Value{
		String: &s,
	}, nil
}