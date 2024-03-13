package main

import (
	"fmt"

	"github.com/H-ADJI/calkilatrice/lexer"
)

func main() {
	input := "sin(cos(90))*5*11-(5-9)"
	lexer := lexer.NewLexer(input)
	lexer.Tokenize()
	fmt.Println(lexer.Tokens)

}
