package main

import (
	"flag"
	"fmt"

	"github.com/H-ADJI/calkilatrice/parser"
)

func main() {
	mathExpr := flag.String("expr", "2*sin(30)", "Mathematical expression to evaluate. Example : 1/2 - (4 - cos(10^2))")
	supportedFunctions := flag.Bool("S", false, "Shows the list of mathematical operations supported by this calculator")
	showTokens := flag.Bool("include-tokens", false, "Add looging the parsed tokens to the output")
	showAst := flag.Bool("include-ast", false, "Add looging the evaluated ast to the output")
	radianToDegree := flag.Bool("use-degrees", false, "Enable Trigonometric functions evaluation using Degrees as a unit instead of Radians")
	flag.Parse()
	lexer := parser.NewLexer(*mathExpr)
	if *supportedFunctions {
		fmt.Println("The list of supported operations is the following:")
		fmt.Println("\t=> cos()")
		fmt.Println("\t=> acos()")
		fmt.Println("\t=> sin()")
		fmt.Println("\t=> asin()")
		fmt.Println("\t=> tan()")
		fmt.Println("\t=> atan()")
		fmt.Println("\t=> sqrt()")
		fmt.Println("\t=> * : for multiplication")
		fmt.Println("\t=> ^ : for exponentiation")
		fmt.Println("\t=> + : for addition")
		fmt.Println("\t=> - : for substraction")
		return
	}
	if *showTokens {
		fmt.Printf("Tokens : %v \n\n", lexer.Tokens())
	}
	p := parser.Paser{}
	ast, err := p.AST(*mathExpr)
	if err != nil {
		fmt.Println(err)
		return
	}
	if *showAst {
		fmt.Println("The Abstract Syntax Tree")
		fmt.Println(ast)
		fmt.Println()

	}
	value := parser.TreeWalk(&ast.Root, *radianToDegree)
	if *mathExpr == "" {
		fmt.Println("No expression provided, Please use --help")
		return
	}
	fmt.Printf("%s = %.2f\n", *mathExpr, value)
}
