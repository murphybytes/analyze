package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"sync"
)

var once sync.Once
var _parser *participle.Parser

//nolint
func Parser() *participle.Parser {
	once.Do(func() {
		def := lexer.MustSimple([]lexer.Rule{
			{"String", `"(\\"|[^"])*"`, nil},
			{"Number", `[-+]?(\d*\.)?\d+`, nil},
			{"whitespace", `[ \t]+`, nil},
			{`Keyword`, `(?i)\b(nil|true|false)\b`, nil},
			{"Operators", `!=|<=|>=|&&|==|\|\||[!()<>,]`, nil},
			{"Variable", `\$[a-zA-Z\-"\._\[\]1-9 ]+`, nil},
			{"Function", `^@[a-zA-Z_]\w*`, nil },
			{"RegularExpression", `/\^?[0-9a-zA-Z\(\)\?\:\[\]\{\}\,\.\-\*\+\\]+\$?/`, nil},
		})
		_parser = participle.MustBuild(&Expression{},
			participle.Lexer(def),
			participle.Unquote("String"),
			participle.UseLookahead(2),
		)
	})
	return _parser
}
