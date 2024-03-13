package main

import (
	"fmt"

	"github.com/H-ADJI/calkilatrice/lexer"
)

func main() {
	input := "5*sin(3)"
	lexer := lexer.NewLexer(input)
	fmt.Println(lexer.Tokens())

}
