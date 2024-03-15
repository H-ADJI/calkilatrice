package parser

import (
	"github.com/H-ADJI/calkilatrice/lexer"
)

type astNode struct {
	token lexer.Token
	left  *astNode
	right *astNode
}

type AST struct {
	root astNode
}

type Paser struct {
	mathExpression string
	tokens         []lexer.Token
	cursor         int
	lookahead      lexer.Token
}

func NewParser(mathExpression string) Paser {
	tokenizer := lexer.NewLexer(mathExpression)
	return Paser{mathExpression: mathExpression, tokens: tokenizer.Tokens()}
}
func (parser *Paser) Next() {
	parser.cursor += 1
	parser.lookahead = parser.tokens[parser.cursor]
}

func (parser *Paser) Consume(tokenType int) lexer.Token {
	if parser.cursor == len(parser.tokens) {
		panic("No more tokens available ")
	}
	if parser.lookahead.TokenType != tokenType {
		panic("Wrong token type")

	}
	defer parser.Next()
	return parser.lookahead
}
func (parser *Paser) Expression() *astNode {
	return parser.addition()
}

func (parser *Paser) addition() *astNode {
	leftNode := parser.terminals()
	for parser.lookahead.IsType(lexer.AddOp, lexer.MinusOp) {
		leftNode = &astNode{token: parser.Consume(parser.lookahead.TokenType), left: leftNode, right: parser.terminals()}
	}
	return leftNode

}
func (parser *Paser) terminals() *astNode {
	if parser.lookahead.IsType(lexer.LeftPar) {
		parser.Consume(lexer.LeftPar)
		exp := parser.Expression()
		parser.Consume(lexer.RightPar)
		return exp

	}
	if parser.lookahead.IsType(lexer.Number) {
		return &astNode{token: parser.Consume(lexer.Number)}
	}
	panic("Unkown Terminal")
}
