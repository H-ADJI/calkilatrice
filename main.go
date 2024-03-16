package main

import (
	"fmt"

	"github.com/H-ADJI/calkilatrice/lexer"
	"github.com/H-ADJI/calkilatrice/parser"
)

func main() {
	input := "22-2^3"
	lexer := lexer.NewLexer(input)
	lexer.Tokens()
	parser := parser.Paser{}
	ast := parser.AST(input)
	fmt.Println(ast)

}
