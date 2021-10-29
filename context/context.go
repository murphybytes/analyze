// Package context defines data that is passed to the Evaluate function and maps to variables. The data takes  the form
// of a graph where intermediate nodes may be of the type map[string]interface{} or an array of interfaces []interface{}.
// An object is defined by map[string]interface{}, the interface{} value my contain other objects, arrays, float64, int,
// string, or bool. Arrays may contain object, string, int, float64, or bool.  IMPORTANT! Arrays can not contain other
// arrays. All elements of arrays must be the same type. Variables  are of the form $segment.segment.segment where a
// segment is of the form key with an optional index qualifier of the form ["key"] that maps to object fields, or
// [index] that maps to an element in an array. For example the variable $foo.bar in the following object would map
// to the value "bozo".   { "foo": { "bing": 23, "bar": "bozo" } }.  This value could also be resolved by $foo["bar"].
// An this example $foo[1].bar[0] maps to the value 23.
// { "foo": [ { bar: [ 3, 4 ] }, { bar: [ 23, 8 ] } ] }
package context

type Context map[string]interface{}
