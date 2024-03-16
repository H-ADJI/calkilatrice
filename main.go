package main

import (
	"fmt"

	"github.com/H-ADJI/calkilatrice/lexer"
	"github.com/H-ADJI/calkilatrice/parser"
)

func main() {
	input := "53333+2 222-2 2222"
	lexer := lexer.NewLexer(input)
	fmt.Println(lexer.Tokens())
	parser := parser.Paser{}
	ast := parser.AST(input)
	fmt.Println(ast)

}
