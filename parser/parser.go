package parser

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/H-ADJI/calkilatrice/lexer"
)

type astNode struct {
	token lexer.Token
	left  *astNode
	right *astNode
}

type AST struct {
	Root astNode
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
	return tree.Root.String()
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
		return &AST{Root: *parser.expression()}
	}
	return &AST{}
}

// func (parser *Paser) characterPosition() int {
// 	var sum int
// 	if parser.cursor == len(parser.tokens)-1 {
// 		sum = -1
// 	}
// 	for _, token := range parser.tokens[:parser.cursor] {
// 		sum += len(token.Value)
// 	}
// 	return sum
// }

func (parser *Paser) expression() *astNode {
	root := parser.addition()
	return root
}

func (parser *Paser) addition() *astNode {
	leftNode := parser.mathFunc()
	for parser.lookahead.IsType(lexer.AddOp, lexer.MinusOp) {
		leftNode = &astNode{token: parser.Consume(parser.lookahead.TokenType), left: leftNode, right: parser.mathFunc()}
	}
	return leftNode

}
func (parser *Paser) mathFunc() *astNode {
	leftNode := parser.multiplication()
	if parser.lookahead.IsType(lexer.LeftPar) {
		parser.Consume(parser.lookahead.TokenType)
		arg := parser.expression()
		parser.Consume(lexer.RightPar)
		return &astNode{token: leftNode.token, left: arg}
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
	if parser.lookahead.IsType(lexer.MathFunc) {
		return &astNode{token: parser.Consume(parser.lookahead.TokenType)}
	}
	panic("Unkown Terminal")
}

func TreeWalk(root *astNode) float64 {
	var result float64
	if root.token.TokenType == lexer.AddOp {
		result = result + TreeWalk(root.left) + TreeWalk(root.right)
		return result
	} else if root.token.TokenType == lexer.MinusOp {
		result = result + TreeWalk(root.left) - TreeWalk(root.right)
		return result
	} else if root.token.TokenType == lexer.MultOp {
		result = result + TreeWalk(root.left)*TreeWalk(root.right)
		return result
	} else if root.token.TokenType == lexer.DivOp {
		result = result + TreeWalk(root.left)/TreeWalk(root.right)
		return result
	} else if root.token.TokenType == lexer.ExpOp {
		result = result + math.Pow(TreeWalk(root.left), TreeWalk(root.right))
		return result
	} else if root.token.TokenType == lexer.MathFunc {
		mathFunc, err := mathFuncEval(root.token.Value)
		if err != nil {
			panic(err)
		}
		return mathFunc(TreeWalk(root.left))
	} else if root.token.IsType(lexer.Number, lexer.NegativeNumber) {
		value, _ := strconv.ParseFloat(root.token.Value, 64)
		return value
	}
	return result
}

func mathFuncEval(funcName string) (func(float64) float64, error) {
	switch funcName {
	case "sin":
		return math.Sin, nil
	case "asin":
		return math.Asin, nil
	case "tan":
		return math.Tan, nil
	case "atan":
		return math.Atan, nil
	case "cos":
		return math.Atan, nil
	case "acos":
		return math.Acos, nil
	case "sqrt":
		return math.Sqrt, nil
	default:
		return nil, errors.New("Unsupported function " + funcName)
	}
}
