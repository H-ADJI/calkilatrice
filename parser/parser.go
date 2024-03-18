package parser

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

type astNode struct {
	token Token
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
	tokens         []Token
	cursor         int
	lookahead      Token
}

func (parser *Paser) Next() {
	parser.cursor += 1
	if parser.cursor >= len(parser.tokens) {
		return
	}
	parser.lookahead = parser.tokens[parser.cursor]
}

func (parser *Paser) Consume(tokenType int) Token {
	if parser.lookahead.TokenType != tokenType {
		panic("Wrong token type")
	}
	defer parser.Next()
	return parser.lookahead
}
func (parser *Paser) AST(mathExpression string) *AST {
	tokenizer := NewLexer(mathExpression)
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

func (parser *Paser) expression() *astNode {
	root := parser.addition()
	return root
}

func (parser *Paser) addition() *astNode {
	leftNode := parser.mathFunc()
	for parser.lookahead.IsType(AddOp, MinusOp) {
		leftNode = &astNode{token: parser.Consume(parser.lookahead.TokenType), left: leftNode, right: parser.mathFunc()}
	}
	return leftNode

}
func (parser *Paser) mathFunc() *astNode {
	leftNode := parser.multiplication()
	if parser.lookahead.IsType(LeftPar) {
		parser.Consume(parser.lookahead.TokenType)
		arg := parser.expression()
		parser.Consume(RightPar)
		return &astNode{token: leftNode.token, left: arg}
	}
	return leftNode
}
func (parser *Paser) multiplication() *astNode {
	leftNode := parser.exponentiation()
	for parser.lookahead.IsType(MultOp, DivOp) {
		leftNode = &astNode{token: parser.Consume(parser.lookahead.TokenType), left: leftNode, right: parser.exponentiation()}
	}
	return leftNode
}

func (parser *Paser) exponentiation() *astNode {
	leftNode := parser.terminals()
	for parser.lookahead.IsType(ExpOp) {
		leftNode = &astNode{token: parser.Consume(parser.lookahead.TokenType), left: leftNode, right: parser.terminals()}
	}
	return leftNode
}
func (parser *Paser) terminals() *astNode {
	if parser.lookahead.IsType(LeftPar) {
		parser.Consume(LeftPar)
		exp := parser.expression()
		parser.Consume(RightPar)
		return exp
	}
	if parser.lookahead.IsType(Number, NegativeNumber) {
		return &astNode{token: parser.Consume(parser.lookahead.TokenType)}
	}
	if parser.lookahead.IsType(MathFunc) {
		return &astNode{token: parser.Consume(parser.lookahead.TokenType)}
	}
	panic("Unkown Terminal")
}

func TreeWalk(root *astNode, useDegrees bool) float64 {
	var result float64
	if root.token.TokenType == AddOp {
		result = result + TreeWalk(root.left, useDegrees) + TreeWalk(root.right, useDegrees)
		return result
	} else if root.token.TokenType == MinusOp {
		result = result + TreeWalk(root.left, useDegrees) - TreeWalk(root.right, useDegrees)
		return result
	} else if root.token.TokenType == MultOp {
		result = result + TreeWalk(root.left, useDegrees)*TreeWalk(root.right, useDegrees)
		return result
	} else if root.token.TokenType == DivOp {
		result = result + TreeWalk(root.left, useDegrees)/TreeWalk(root.right, useDegrees)
		return result
	} else if root.token.TokenType == ExpOp {
		result = result + math.Pow(TreeWalk(root.left, useDegrees), TreeWalk(root.right, useDegrees))
		return result
	} else if root.token.TokenType == MathFunc {
		mathFunc, err := mathFuncEval(root.token.Value)
		if err != nil {
			panic(err)
		}
		if useDegrees {
			result = mathFunc(TreeWalk(root.left, useDegrees) * math.Pi / 180)
		} else {
			result = mathFunc(TreeWalk(root.left, useDegrees))
		}
		return result
	} else if root.token.IsType(Number, NegativeNumber) {
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
		return math.Cos, nil
	case "acos":
		return math.Acos, nil
	case "sqrt":
		return math.Sqrt, nil
	default:
		return nil, errors.New("Unsupported function " + funcName)
	}
}
