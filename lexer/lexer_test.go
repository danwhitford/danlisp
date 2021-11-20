package lexer

import (
	"testing"
	"whitford.io/danlisp/token"
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
	tokens, _ := GetTokens(input)
	assertType(t, token.LB, tokens[0].TokenType)
}

func TestSingleRightBrack(t *testing.T) {
	input := ")"
	tokens, _ := GetTokens(input)
	assertType(t, token.RB, tokens[0].TokenType)
}

func TestBracketPair(t *testing.T) {
	input := "()"
	tokens, _ := GetTokens(input)
	assertType(t, token.LB, tokens[0].TokenType)
	assertType(t, token.RB, tokens[1].TokenType)
}

func TestKeyword(t *testing.T) {
	input := "let"
	tokens, _ := GetTokens(input)
	assertType(t, token.KEYWORD, tokens[0].TokenType)
	assertString(t, "let", tokens[0].Lexeme)
}

func TestNumber(t *testing.T) {
	input := "123.7"
	tokens, _ := GetTokens(input)
	assertType(t, token.LITERAL, tokens[0].TokenType)
	assertString(t, "123.7", tokens[0].Lexeme)
	assertNumber(t, 123.7, tokens[0].Value.(float64))
}

func TestErrorInNumber(t *testing.T) {
	input := "123notanumber"
	_, err := GetTokens(input)
	assertString(t, err.Error(), "error while lexing on line 1. '123notanumber' is not a number")
}

func TestString(t *testing.T) {
	input := "\"i am the fly\""
	tokens, _ := GetTokens(input)
	assertType(t, token.LITERAL, tokens[0].TokenType)
	assertString(t, "\"i am the fly\"", tokens[0].Lexeme)
	assertString(t, "i am the fly", tokens[0].Value.(string))
}

func TestEOLInString(t *testing.T) {
	input := "\"i am the fly\nfly in the fly in the\""
	_, err := GetTokens(input)
	assertString(t, err.Error(), "error while lexing on line 1. reached end of line in string '\"i am the fly'")
}

func TestEOFInString(t *testing.T) {
	input := "\"i am the fly"
	_, err := GetTokens(input)
	assertString(t, err.Error(), "error while lexing on line 1. reached end of input in string '\"i am the fly'")
}
