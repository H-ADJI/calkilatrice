package main

import (
	"flag"
	"fmt"

	"github.com/H-ADJI/calkilatrice/parser"
)

func main() {
	mathExpr := flag.String("expr", "", "Mathematical expression to evaluate. Example : 1/2 - (4 - cos(10+2))")
	showTokens := flag.Bool("include-tokens", false, "Add looging the parsed tokens to the output")
	showAst := flag.Bool("include-ast", false, "Add looging the evaluated ast to the output")
	radianToDegree := flag.Bool("use-degrees", false, "Enable Trigonometric functions evaluation using Degrees as a unit instead of Radians")
	flag.Parse()
	lexer := parser.NewLexer(*mathExpr)
	if *showTokens {
		fmt.Println(lexer.Tokens())
	}
	p := parser.Paser{}
	ast := p.AST(*mathExpr)
	if *showAst {
		fmt.Println(ast)
	}
	value := parser.TreeWalk(&ast.Root, *radianToDegree)
	fmt.Println(value)
}
