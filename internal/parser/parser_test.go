package parser

import (
	"testing"

	"whitford.io/danlisp/internal/expr"
	"whitford.io/danlisp/internal/lexer"
)

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

func TestNestedSeq(t *testing.T) {
	input := "(concat (0 1 2) (3 4 5))"
	tokens, _ := lexer.GetTokens(input)
	expressions, _ := GetExpressions(tokens)
	sym, ok := expressions[0].(expr.Seq).Exprs[0].(expr.Symbol)
	if !ok {
		t.Fatalf("Expected symbol to be symbol")
	}
	assertString(t, "concat", sym.Name)

	firstNested := expressions[0].(expr.Seq).Exprs[1].(expr.Seq)
	for i := 0; i < 3; i++ {
		val := firstNested.Exprs[i].(expr.Atom).Value.(float64)
		assertNumber(t, float64(i), val)
	}

	secondNested := expressions[0].(expr.Seq).Exprs[2].(expr.Seq)
	for i := 0; i < 3; i++ {
		val := secondNested.Exprs[i].(expr.Atom).Value.(float64)
		assertNumber(t, float64(i)+3, val)
	}
}

func TestErrorWhenSeqNotClosed(t *testing.T) {
	input := "(+ 1 2 3 4"
	tokens, _ := lexer.GetTokens(input)
	_, err := GetExpressions(tokens)
	assertString(t, "parse error. missing ')' to close sequence.", err.Error())
}
