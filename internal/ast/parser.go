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
		  //{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`, nil},
		   {"whitespace", `[ \t]+`, nil},
		   {`Keyword`, `(?i)\b(true|false)\b`, nil },
		   {"Operators", `!=|<=|>=|&&|==|\|\||[!()<>]`, nil},
		   {"Variable", `^\$[a-zA-Z\-"\._\[\]1-9 ]+`, nil },
		   {"Function", `(?i)\b(len|in)\b`, nil },
	   })
	   _parser = participle.MustBuild(&Expression{},
		   participle.Lexer(def),
		   participle.Unquote("String"),
		   participle.UseLookahead(2),
		   )
   })
   return _parser
}