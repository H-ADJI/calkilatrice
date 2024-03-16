package parser

import (
	"fmt"
	"strings"

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

func (node *astNode) String() string {
	return node.stringWithIndent("", true)
}

func (node *astNode) stringWithIndent(indent string, isTail bool) string {
	if node == nil {
		return ""
	}

	var builder strings.Builder

	if node.right != nil {
		builder.WriteString(node.right.stringWithIndent(indent+"│   ", false))
	}

	builder.WriteString(indent)
	if isTail {
		builder.WriteString("└── ")
	} else {
		builder.WriteString("├── ")
	}
	builder.WriteString(node.token.String() + "\n")

	if node.left != nil {
		builder.WriteString(node.left.stringWithIndent(indent+"    ", true))
	}

	return builder.String()
}

func (tree *AST) String() string {
	return tree.root.String()
}

type Paser struct {
	mathExpression string
	tokens         []lexer.Token
	cursor         int
	lookahead      lexer.Token
}

func (parser *Paser) Next() {
	parser.cursor += 1
	if parser.cursor >= len(parser.tokens) {
		return
	}
	parser.lookahead = parser.tokens[parser.cursor]
}

func (parser *Paser) Consume(tokenType int) lexer.Token {
	if parser.lookahead.TokenType != tokenType {
		panic("Wrong token type")
	}
	defer parser.Next()
	return parser.lookahead
}
func (parser *Paser) AST(mathExpression string) *AST {
	tokenizer := lexer.NewLexer(mathExpression)
	tokens := tokenizer.Tokens()
	parser.tokens = tokens
	parser.mathExpression = mathExpression
	if len(tokens) > 0 {
		parser.lookahead = tokens[0]
		parser.cursor = 0
		return &AST{root: *parser.expression()}
	}
	return &AST{}
}

func (parser *Paser) characterPosition() int {
	var sum int
	if parser.cursor == 1 || parser.cursor == len(parser.tokens)-1 {
		sum = -1
	}
	for _, token := range parser.tokens[:parser.cursor] {
		sum += len(token.Value)
	}
	return sum
}

func (parser *Paser) expression() *astNode {
	root := parser.addition()
	if parser.cursor < len(parser.tokens) {
		errMsg := fmt.Sprintf("Invalid syntax : at position %d ==> %v", parser.characterPosition()+1, string(parser.mathExpression))
		errCursor := strings.Repeat(" ", len(errMsg)-len(parser.mathExpression)+parser.characterPosition())
		fmt.Printf("%v\n%v^\n", errMsg, errCursor)
	}
	return root
}

func (parser *Paser) addition() *astNode {
	leftNode := parser.multiplication()
	for parser.lookahead.IsType(lexer.AddOp, lexer.MinusOp) {
		leftNode = &astNode{token: parser.Consume(parser.lookahead.TokenType), left: leftNode, right: parser.multiplication()}
	}
	return leftNode

}

func (parser *Paser) multiplication() *astNode {
	leftNode := parser.exponentiation()
	for parser.lookahead.IsType(lexer.MultOp, lexer.DivOp) {
		leftNode = &astNode{token: parser.Consume(parser.lookahead.TokenType), left: leftNode, right: parser.exponentiation()}
	}
	return leftNode
}

func (parser *Paser) exponentiation() *astNode {
	leftNode := parser.terminals()
	for parser.lookahead.IsType(lexer.ExpOp) {
		leftNode = &astNode{token: parser.Consume(parser.lookahead.TokenType), left: leftNode, right: parser.terminals()}
	}
	return leftNode
}
func (parser *Paser) terminals() *astNode {
	if parser.lookahead.IsType(lexer.LeftPar) {
		parser.Consume(lexer.LeftPar)
		exp := parser.expression()
		parser.Consume(lexer.RightPar)
		return exp

	}
	if parser.lookahead.IsType(lexer.Number, lexer.NegativeNumber) {
		return &astNode{token: parser.Consume(parser.lookahead.TokenType)}
	}
	panic("Unkown Terminal")
}
