package ast

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/murphybytes/analyze/context"
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

type Variable string

func (v *Variable) Capture(s []string) error {
	*v = Variable(strings.TrimPrefix(strings.Join(s, ""), "$"))
	return nil
}

func (v *Variable) Eval(ctx context.Context) (*Value, error) {
	keys := strings.Split(string(*v), ".")
	return walkCtx(keys, ctx)
}
// Traverse variable segments left to right using each segment to look up object in context data
// until we get to the get to the last element, then return its value.
func walkCtx(keys []string, ctx context.Context)(*Value, error){
	key := keys[0]; keys = keys[1:]
	inf, err := extractVariableElement(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(keys) > 0 {
		nextCtx, ok := inf.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected context element not correct type")
		}
		return walkCtx(keys, nextCtx)
	}
	// we are at our terminal element convert to appropriate value type
	return convertToValue(inf)
}
// matches foo[ "key" ]
var regexObjectRef = regexp.MustCompile(`^[\w\-]*\[\s*"[\w\-]+"\s*\]$`)
// matches foo[2]
var regexArrayRef =  regexp.MustCompile(`^[\w\-]*\[\s*[0-9]+\s*\]$`)

// Variable names can include index expressions to map into an object, or to reference particular array elements
// i.e. $foo.someObject["field"] or $foo.someArray[3]. This function extracts the value referenced by the key and
// returns it. It also handles the case when the  root element refers to an array, number, string etc.
func extractVariableElement(ctx context.Context, reference string)(interface{}, error){
	reference = strings.Trim(reference, " ")

	switch t := ctx.(type) {
	case []interface{}:
		return resolveArrayElement(t, reference)
	case map[string]interface{}:
		return resolveObjectField(t, reference)
	}

	// pass through scalar types
	return ctx, nil
}

func convertToValue(intf interface{})(*Value,error){
	var val Value
	switch t := intf.(type) {
	case int:
		f := float64(t)
		val.Number = &f
	case *int:
		f := float64(*t)
		val.Number = &f
	case float64:
		val.Number = &t
	case *float64:
		val.Number = t
	case string:
		val.String = &t
	case *string:
		val.String = t
	case map[string]interface{}:
		val.Object = t
	case []interface{}:
		val.Array = t
	case nil:
		val.NilSet = true
	default:
		return nil, UnsupportedTypeError(intf)
	}
	return &val, nil
}

func resolveArrayElement(arr []interface{}, reference string)(interface{}, error){
	// make sure the variable expression looks like an array
	if !regexArrayRef.MatchString(reference) {
		return nil, NewSyntaxError("expected array reference got %q", reference)
	}
	pts := strings.Split(reference, "[")
	indexStr := pts[1]
	index, err  := strconv.Atoi(strings.TrimRight(indexStr, ` ]`))
	if err != nil {
		return nil, NewSyntaxError("error resolving array index %q", err )
	}
	if !(index < len(arr)) {
		return nil, NewSyntaxError("index out of range")
	}
	return arr[index], nil
}

func resolveObjectField(obj map[string]interface{}, reference string)(interface{},error){
	// handle index into object object["field"]
	if regexObjectRef.MatchString(reference) {
		p := strings.Split(reference, "[")
		key, index := p[0], p[1]
		index = strings.Trim(index, ` ]"`)
		// we  expect an object
		var ok bool
		if len(key) > 0 {
			if obj, ok = obj[key].(map[string]interface{}); !ok {
				return nil, MissingKeyError(key)
			}
		}
		var result interface{}
		if result, ok = obj[index]; !ok {
			return nil, IndexOutOfRangeError(reference)
		}

		return result, nil
	}

	// handle index into array array[3]
	if regexArrayRef.MatchString(reference) {
		p := strings.Split(reference, "[")
		key, str := p[0], p[1]
		// we expect an array
		index, err := strconv.Atoi(strings.Trim(str, ` ]`))
		if err != nil {
			return nil, NewSyntaxError(fmt.Sprintf("can't resolve %s into an array element", reference))
		}
		arr, ok := obj[key].([]interface{})
		if !ok {
			return nil, MissingKeyError(key)
		}

		if !(index < len(arr)) {
			return nil, IndexOutOfRangeError(reference)
		}

		return arr[index], nil
	}

	return obj[reference], nil
}