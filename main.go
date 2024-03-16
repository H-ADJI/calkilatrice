package main

import (
	"fmt"

	"github.com/H-ADJI/calkilatrice/parser"
)

func main() {
	input := "5+2-13"
	parser := parser.NewParser(input)
	ast := parser.AST()
	fmt.Println(ast)

}
