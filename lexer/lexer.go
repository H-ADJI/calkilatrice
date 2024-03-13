package lexer

import (
	"fmt"
	"regexp"
	"strings"
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

var calculatorGrammar = Grammar{
	Rule{pattern: regexp.MustCompile(`^-?\d+(?:\.\d+)?`), tokenType: number},
	Rule{pattern: regexp.MustCompile(`^\+`), tokenType: addOp},
	Rule{pattern: regexp.MustCompile(`^\s+`), tokenType: whiteSpace},
	Rule{pattern: regexp.MustCompile(`^[a-zA-Z]+`), tokenType: mathFunc},
	Rule{pattern: regexp.MustCompile(`^\-`), tokenType: minusOp},
	Rule{pattern: regexp.MustCompile(`^\*`), tokenType: multOp},
	Rule{pattern: regexp.MustCompile(`^/`), tokenType: divOp},
	Rule{pattern: regexp.MustCompile(`^\^`), tokenType: expOp},
	Rule{pattern: regexp.MustCompile(`^\(`), tokenType: rightPar},
	Rule{pattern: regexp.MustCompile(`^\)`), tokenType: leftPar},
}

type Token struct {
	value     string
	tokenType int
}

func (token Token) String() string {
	return fmt.Sprintf("Token[%v]", token.value)
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
	fmt.Println(lexer.cursor)
	fmt.Println(len(lexer.expr))
	if lexer.cursor == len(lexer.expr) {
		return
	}
	for _, rule := range calculatorGrammar {
		currentExpr := lexer.expr[lexer.cursor:]
		token := rule.pattern.Find(currentExpr)
		if token != nil {
			lexer.cursor += len(token)
			// White spaces are ignored, no token is added to our token list
			if rule.tokenType == whiteSpace {
				lexer.tokenize()
				return
			} else {
				lexer.tokens = append(lexer.tokens, Token{value: string(token), tokenType: rule.tokenType})
				lexer.tokenize()
				return
			}
		}
	}

	errMsg := fmt.Sprintf("unkown token at position %d in: %v", lexer.cursor+1, string(lexer.expr))
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
