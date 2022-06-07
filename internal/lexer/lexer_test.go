package lexer

import (
	"testing"

	"github.com/danwhitford/danlisp/internal/token"
)

func assertType(t *testing.T, expected, actual token.TokenType) {
	if expected != actual {
		t.Fatalf("Assertion failed. Expected '%v' but got '%v'", expected, actual)
	}
}

func assertString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("Assertion failed. Expected '%v' but got '%v'", expected, actual)
	}
}

func assertNumber(t *testing.T, expected, actual float64) {
	if expected != actual {
		t.Fatalf("Assertion failed. Expected '%v' but got '%v'", expected, actual)
	}
}

func TestSingleLeftBrack(t *testing.T) {
	input := "("
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.LB, tokens[0].TokenType)
}

func TestSingleRightBrack(t *testing.T) {
	input := ")"
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.RB, tokens[0].TokenType)
}

func TestBracketPair(t *testing.T) {
	input := "()"
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.LB, tokens[0].TokenType)
	assertType(t, token.RB, tokens[1].TokenType)
}

func TestKeyword(t *testing.T) {
	input := "let"
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.KEYWORD, tokens[0].TokenType)
	assertString(t, "let", tokens[0].Lexeme)
}

func TestNumber(t *testing.T) {
	input := "123.7"
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.LITERAL, tokens[0].TokenType)
	assertString(t, "123.7", tokens[0].Lexeme)
	assertNumber(t, 123.7, tokens[0].Value.(float64))
}

func TestErrorInNumber(t *testing.T) {
	input := "123notanumber"
	lex := NewLexer(input)
	_, err := lex.GetTokens()
	assertString(t, err.Error(), "error while lexing on line 1. '123notanumber' is not a number")
}

func TestString(t *testing.T) {
	input := "\"i am the fly\""
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.LITERAL, tokens[0].TokenType)
	assertString(t, "\"i am the fly\"", tokens[0].Lexeme)
	assertString(t, "i am the fly", tokens[0].Value.(string))
}

func TestEOLInString(t *testing.T) {
	input := "\"i am the fly\nfly in the fly in the\""
	lex := NewLexer(input)
	_, err := lex.GetTokens()
	assertString(t, err.Error(), "error while lexing on line 1. reached end of line in string '\"i am the fly'")
}

func TestEOFInString(t *testing.T) {
	input := "\"i am the fly"
	lex := NewLexer(input)
	_, err := lex.GetTokens()
	assertString(t, err.Error(), "error while lexing on line 1. reached end of input in string '\"i am the fly'")
}

func TestSeq(t *testing.T) {
	input := "(1 2 3)"
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.LB, tokens[0].TokenType)
	assertType(t, token.LITERAL, tokens[1].TokenType)
	assertType(t, token.LITERAL, tokens[2].TokenType)
	assertType(t, token.LITERAL, tokens[3].TokenType)
	assertType(t, token.RB, tokens[4].TokenType)
}

func TestDefinition(t *testing.T) {
	input := "(set x 5)"
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.LB, tokens[0].TokenType)
	assertType(t, token.SET, tokens[1].TokenType)
	assertType(t, token.KEYWORD, tokens[2].TokenType)
	assertType(t, token.LITERAL, tokens[3].TokenType)
	assertType(t, token.RB, tokens[4].TokenType)
}

func TestIf(t *testing.T) {
	input := "(if (= 5 x))"
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.LB, tokens[0].TokenType)
	assertType(t, token.IF, tokens[1].TokenType)
}

func TestWhile(t *testing.T) {
	input := "(while (> x 0) (set x (- x 1)))"
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.LB, tokens[0].TokenType)
	assertType(t, token.WHILE, tokens[1].TokenType)
}

func TestWorksWithWhitespace(t *testing.T) {
	input := `l
	
	`
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.KEYWORD, tokens[0].TokenType)
	assertString(t, "l", tokens[0].Lexeme)
}

func TestFor(t *testing.T) {
	input := "(for (set i 0) (lt i 10) (set i (+ i 1)) (prn i))"
	lex := NewLexer(input)
	tokens, _ := lex.GetTokens()
	assertType(t, token.LB, tokens[0].TokenType)
	assertType(t, token.FOR, tokens[1].TokenType)
}
