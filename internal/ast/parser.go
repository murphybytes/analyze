package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"sync"
)

var once sync.Once
var _parser *participle.Parser


func Parser() *participle.Parser {
   once.Do(func() {
	   def := lexer.MustSimple([]lexer.Rule{
		   {"String", `"(\\"|[^"])*"`, nil},
		   {"Number", `[-+]?(\d*\.)?\d+`, nil},
		 //  {"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`, nil},
		   {"whitespace", `[ \t]+`, nil},
		   {`Keyword`, `(?i)\b(true|false)\b`, nil },
		   {"ComparisonOperators", `<`, nil },
		  {"LogicalOperators", `&`, nil },
		  {"UnaryOperators", "!", nil },
	   })
	   _parser = participle.MustBuild(&Expression{},
		   participle.Lexer(def),
		   participle.Unquote("String"),
		   participle.UseLookahead(2),
		   )
   })
   return _parser
}