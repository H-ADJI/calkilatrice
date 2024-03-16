package main

import (
	"fmt"

	"github.com/H-ADJI/calkilatrice/lexer"
	"github.com/H-ADJI/calkilatrice/parser"
)

func main() {
	input := "1+sin(tan(22*5-1))"
	lexer := lexer.NewLexer(input)
	fmt.Println(lexer.Tokens())
	parser := parser.Paser{}
	ast := parser.AST(input)
	fmt.Println(ast)

}
