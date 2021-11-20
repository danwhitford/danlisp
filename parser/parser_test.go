package parser

import (
	"testing"
	"whitford.io/danlisp/expr"
	"whitford.io/danlisp/lexer"
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

func TestNumber(t *testing.T) {
	input := "123.7"
	tokens, _ := lexer.GetTokens(input)
	expressions, _ := GetExpressions(tokens)
	assertNumber(t, 123.7, expressions[0].(expr.Atom).Value.(float64))
}

func TestString(t *testing.T) {
	input := "\"egghead\""
	tokens, _ := lexer.GetTokens(input)
	expressions, _ := GetExpressions(tokens)
	assertString(t, "egghead", expressions[0].(expr.Atom).Value.(string))
}

func TestSeq(t *testing.T) {
	input := "(0 1 2 3 4 5)"
	tokens, _ := lexer.GetTokens(input)
	expressions, _ := GetExpressions(tokens)
	seq := expressions[0].(expr.Seq)
	for i := 0; i < 6; i++ {
		val := seq.Exprs[i].(expr.Atom).Value.(float64)
		assertNumber(t, float64(i), val)
	}
}
