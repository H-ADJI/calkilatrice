package lexer

import (
	"fmt"
	"regexp"
)

const DELIM = "--"

const (
	whiteSpace = iota
	addOp
	minusOp
	multOp
	expOp
	divOp
	number
	rightPar
	leftPar
	mathFunc
)

type Rule struct {
	pattern   *regexp.Regexp
	tokenType int
}
type Grammar []Rule

var calculatorGrammar = []Grammar{
	{Rule{pattern: regexp.MustCompile(`^\d+`), tokenType: number}},
	{Rule{pattern: regexp.MustCompile(`^[a-zA-Z]+`), tokenType: mathFunc}},
	{Rule{pattern: regexp.MustCompile(`^\s+`), tokenType: whiteSpace}},
	{Rule{pattern: regexp.MustCompile(`^\+`), tokenType: addOp}},
	{Rule{pattern: regexp.MustCompile(`^\-`), tokenType: minusOp}},
	{Rule{pattern: regexp.MustCompile(`^\*`), tokenType: multOp}},
	{Rule{pattern: regexp.MustCompile(`^/`), tokenType: divOp}},
	{Rule{pattern: regexp.MustCompile(`^\^`), tokenType: expOp}},
	{Rule{pattern: regexp.MustCompile(`^\(`), tokenType: rightPar}},
	{Rule{pattern: regexp.MustCompile(`^\)`), tokenType: leftPar}},
}

type Token struct {
	value     string
	tokenType int
}

func (token Token) String() string {
	return fmt.Sprintf("[%v : %v]", token.tokenType, token.value)
}

type Lexer struct {
	expr   string
	cursor int
	tokens []Token
}

func NewLexer(mathExpression string) Lexer {
	return Lexer{expr: mathExpression}
}

// func (lexer Lexer) String() string {
// 	builder := strings.Builder{}
// 	for _,t := lexer.tokens {
// 		builder.WriteString(t)
// 	}
// 	return fmt.Sprintf("[%v : %v]", token.tokenType, token.value)
// }
