package parser

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const DELIM = "--"

const (
	WhiteSpace = iota
	AddOp
	MinusOp
	MultOp
	ExpOp
	DivOp
	Number
	NegativeNumber
	RightPar
	LeftPar
	MathFunc
)

type Rule struct {
	pattern   *regexp.Regexp
	tokenType int
}
type Grammar []Rule

var calculatorGrammar = Grammar{
	Rule{pattern: regexp.MustCompile(`^\d+(?:\.\d+)?`), tokenType: Number},
	Rule{pattern: regexp.MustCompile(`^-\d+(?:\.\d+)?`), tokenType: NegativeNumber},
	Rule{pattern: regexp.MustCompile(`^\-`), tokenType: MinusOp},
	Rule{pattern: regexp.MustCompile(`^\+`), tokenType: AddOp},
	Rule{pattern: regexp.MustCompile(`^\s+`), tokenType: WhiteSpace},
	Rule{pattern: regexp.MustCompile(`^[a-zA-Z]+`), tokenType: MathFunc},
	Rule{pattern: regexp.MustCompile(`^\*`), tokenType: MultOp},
	Rule{pattern: regexp.MustCompile(`^/`), tokenType: DivOp},
	Rule{pattern: regexp.MustCompile(`^\^`), tokenType: ExpOp},
	Rule{pattern: regexp.MustCompile(`^\(`), tokenType: LeftPar},
	Rule{pattern: regexp.MustCompile(`^\)`), tokenType: RightPar},
}

type Token struct {
	Value     string
	TokenType int
}

func (token Token) String() string {
	return fmt.Sprintf("[%v]", token.Value)
}
func (token Token) IsType(types ...int) bool {
	return slices.Contains(types, token.TokenType)
}

type Lexer struct {
	expr   []byte
	cursor int
	tokens []Token
}

func NewLexer(mathExpression string) Lexer {
	return Lexer{expr: []byte(mathExpression)}
}

func (lexer *Lexer) tokenize() {
	if lexer.cursor == len(lexer.expr) {
		return
	}
	for _, rule := range calculatorGrammar {
		currentExpr := lexer.expr[lexer.cursor:]
		match := rule.pattern.Find(currentExpr)
		if match != nil {
			lexer.cursor += len(match)
			// White spaces are ignored, no token is added to our token list
			if rule.tokenType == WhiteSpace {
				lexer.tokenize()
				return
			} else if rule.tokenType == NegativeNumber {
				// when we have
				if len(lexer.tokens) != 0 {
					lexer.tokens = append(lexer.tokens, Token{Value: "0", TokenType: Number})
					lexer.tokens = append(lexer.tokens, Token{Value: "+", TokenType: AddOp})
				}
				lexer.tokens = append(lexer.tokens, Token{Value: string(match), TokenType: rule.tokenType})
				lexer.tokenize()
				return
			} else {
				lexer.tokens = append(lexer.tokens, Token{Value: string(match), TokenType: rule.tokenType})
				lexer.tokenize()
				return
			}
		}
	}

	errMsg := fmt.Sprintf("Invalid Syntax : unkown token at position %d ==> %v", lexer.cursor+1, string(lexer.expr))
	errCursor := strings.Repeat(" ", len(errMsg)-len(lexer.expr)+lexer.cursor)
	fmt.Printf("%v\n%v^\n", errMsg, errCursor)
	lexer.tokens = nil
}

func (lexer *Lexer) Tokens() []Token {
	if lexer.tokens == nil {
		lexer.tokenize()
	}
	return lexer.tokens
}
