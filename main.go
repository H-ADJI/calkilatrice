package main

import (
	"fmt"

	"github.com/H-ADJI/calkilatrice/lexer"
	"github.com/H-ADJI/calkilatrice/parser"
)

func main() {
	input := "sin(90)+ (120)"
	lexer := lexer.NewLexer(input)
	fmt.Println(lexer.Tokens())
	p := parser.Paser{}
	ast := p.AST(input)
	value := parser.TreeWalk(&ast.Root)
	fmt.Println(ast)
	fmt.Println(value)

}
