# Analyze

A predicate expression evaluator designed to embed logical expressions in YAML or JSON. Think of being able to 
apply a SQL where clause on a set of data that you supply.  For example, let say we have an array of objects.
```json
[
  {
    "firstName": "John", 
    "lastName": "Smith"
  },
  {
    "firstName": "Mary",
    "lastName": "Jones"
  }
]
```
You could then apply the following expression to see if there were someone with the first name John in your collection of 
objects. 
```sql
@len( @select(data, "$elt.firstName == 'John'")) > 0
```
This becomes useful when you want to build tools that perform analysis defined in a configuration file, so you might
collect data on resources in kubernetes and check to see if a particular value is defined in a ConfigMap. 
The Go code to perform this looks like this.  In this case we'd put the results of our collection in a map[string]interface{}
rooted at 'resources' and 'config_maps' would refer to an array of config map resources. 
```go
ctx, err := context.New(data)
if err != nil {
	log.Fatal(err)
}
result, err := expression.Evaluate(ctx, `@len( @select( $resources.config_maps, "$elt.firstName == 'John'" ) ) > 0`)
if result {
	fmt.Println("found John!")
}

```
Expressions are composed of operators needed to compose a predicate expression and a way to declare variables that will 
reference a supplied data set as well as a set of builtin functions. If the existing functionality doesn't support your 
use case it's easy to define your own functions. 
```go
func MyFunction(args  []interface{})(interface{}, error) {
	// you implementation 
}

ctx, _ := context.New(data, context.Func("@myfunc", MyFunction))
result, _ := expression.Evaluate(ctx, `@myfunc($foo) == "hello"`)

```