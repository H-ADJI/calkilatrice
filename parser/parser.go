package parser

import (
	"errors"
	"fmt"
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
func debugHelper(tokens []Token, until int) int {
	length := 0
	for i, t := range tokens {
		if until == i {
			break
		}
		length += len(t.Value)
	}
	return length
}
func (parser *Paser) Consume(tokenType int) (Token, error) {
	if parser.lookahead.TokenType != tokenType {
		errMsg := fmt.Sprintf("unable to parse Expression, invalid syntax at position %d ==> %s", parser.cursor-1, string(parser.mathExpression))
		errCursor := strings.Repeat(" ", len(errMsg)-len(parser.mathExpression)+debugHelper(parser.tokens, parser.cursor)+1)
		return Token{}, fmt.Errorf("%s\n%s^", errMsg, errCursor)
	}
	defer parser.Next()
	return parser.lookahead, nil
}
func (parser *Paser) AST(mathExpression string) (*AST, error) {
	tokenizer := NewLexer(mathExpression)
	tokens := tokenizer.Tokens()
	parser.tokens = tokens
	parser.mathExpression = mathExpression
	if len(tokens) > 0 {
		parser.lookahead = tokens[0]
		parser.cursor = 0
		root, err := parser.expression()
		if err != nil {
			return nil, err
		}
		if len(tokens) > parser.cursor {
			errMsg := fmt.Sprintf("unable to parse Expression, invalid syntax at position %d ==> %s", parser.cursor-1, string(parser.mathExpression))
			errCursor := strings.Repeat(" ", len(errMsg)-len(parser.mathExpression)+debugHelper(parser.tokens, parser.cursor)+1)
			return nil, errors.New(errMsg + "\n" + errCursor)
		}
		return &AST{Root: *root}, nil
	}
	return &AST{}, nil
}

func (parser *Paser) expression() (*astNode, error) {
	root, err := parser.addition()
	if err != nil {
		return nil, err
	}
	return root, nil
}

func (parser *Paser) addition() (*astNode, error) {
	leftNode, err := parser.mathFunc()
	if err != nil {
		return nil, err
	}
	for parser.lookahead.IsType(AddOp, MinusOp) {
		token, err := parser.Consume(parser.lookahead.TokenType)
		if err != nil {
			return nil, err
		}
		right, err := parser.mathFunc()
		if err != nil {
			return nil, err
		}
		leftNode = &astNode{token: token, left: leftNode, right: right}
	}
	return leftNode, nil

}
func (parser *Paser) mathFunc() (*astNode, error) {
	leftNode, err := parser.multiplication()
	if err != nil {
		return nil, err
	}
	if parser.lookahead.IsType(LeftPar) {
		_, err := parser.Consume(parser.lookahead.TokenType)
		if err != nil {
			return nil, err
		}
		arg, err := parser.expression()
		if err != nil {
			return nil, err
		}
		_, err = parser.Consume(RightPar)
		if err != nil {
			return nil, err
		}
		if leftNode.token.TokenType == MathFunc {
			return &astNode{token: leftNode.token, left: arg}, nil
		}
		leftNode.right.left = arg
		return leftNode, nil
	}
	return leftNode, nil
}
func (parser *Paser) multiplication() (*astNode, error) {
	leftNode, err := parser.exponentiation()
	if err != nil {
		return nil, err
	}
	for parser.lookahead.IsType(MultOp, DivOp) {
		token, err := parser.Consume(parser.lookahead.TokenType)
		if err != nil {
			return nil, err
		}
		right, err := parser.exponentiation()
		if err != nil {
			return nil, err
		}
		leftNode = &astNode{token: token, left: leftNode, right: right}
	}
	return leftNode, nil
}

func (parser *Paser) exponentiation() (*astNode, error) {
	leftNode, err := parser.terminals()
	if err != nil {
		return nil, err
	}
	for parser.lookahead.IsType(ExpOp) {
		token, err := parser.Consume(parser.lookahead.TokenType)
		if err != nil {
			return nil, err
		}
		right, err := parser.terminals()
		if err != nil {
			return nil, err
		}
		leftNode = &astNode{token: token, left: leftNode, right: right}
	}
	return leftNode, nil
}
func (parser *Paser) terminals() (*astNode, error) {
	if parser.lookahead.IsType(LeftPar) {
		parser.Consume(LeftPar)
		exp, err := parser.expression()
		if err != nil {
			return nil, err
		}
		parser.Consume(RightPar)
		return exp, nil
	}
	if parser.lookahead.IsType(Number, NegativeNumber) {
		token, err := parser.Consume(parser.lookahead.TokenType)
		if err != nil {
			return nil, err
		}
		return &astNode{token: token}, nil
	}
	if parser.lookahead.IsType(MathFunc) {
		token, err := parser.Consume(parser.lookahead.TokenType)
		if err != nil {
			return nil, err
		}
		return &astNode{token: token}, nil
	}
	errMsg := fmt.Sprintf("unable to parse Expression, invalid syntax at position %d ==> %s", parser.cursor-1, string(parser.mathExpression))
	errCursor := strings.Repeat(" ", len(errMsg)-len(parser.mathExpression)+debugHelper(parser.tokens, parser.cursor))
	return nil, errors.New(errMsg + "\n" + errCursor + "^")
	// return nil, errors.New("Invalid Syntax while parsing token : " + parser.lookahead.Value)
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
			if root.token.Value == "cos" || root.token.Value == "tan" || root.token.Value == "sin" {
				result = mathFunc(TreeWalk(root.left, useDegrees) * math.Pi / 180)
			} else {
				result = mathFunc(TreeWalk(root.left, useDegrees))

			}
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
