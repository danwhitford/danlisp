package interpreter

import (
	"testing"

	"whitford.io/danlisp/expr"
	"whitford.io/danlisp/lexer"
	"whitford.io/danlisp/parser"
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

func getExpressions(input string) []expr.Expr {
	tokens, _ := lexer.GetTokens(input)
	exprs, _ := parser.GetExpressions(tokens)
	return exprs
}

func TestJustNumber(t *testing.T) {
	exprs := getExpressions("101")
	ret, _ := Interpret(exprs)
	assertNumber(t, 101, ret.(float64))
}

func TestJustString(t *testing.T) {
	exprs := getExpressions("\"testing testing\"")
	ret, _ := Interpret(exprs)
	assertString(t, "testing testing", ret.(string))
}

func TestAddTwo(t *testing.T) {
	exprs := getExpressions("(+ 2 7)")
	ret, _ := Interpret(exprs)
	assertNumber(t, 9, ret.(float64))
}

func TestSubtract(t *testing.T) {
	exprs := getExpressions("(- 2 7)")
	ret, _ := Interpret(exprs)
	assertNumber(t, -5, ret.(float64))
}

func TestErrorFuncNotFound(t *testing.T) {
	exprs := getExpressions("(nonsuch 2 7)")
	_, err := Interpret(exprs)
	assertString(t, "Runtime error. Could not find function 'nonsuch'.", err.Error())
}

func TestMoreBasicOperators(t *testing.T) {
	exprs := getExpressions("(* 2 7)")
	ret, _ := Interpret(exprs)
	assertNumber(t, 14, ret.(float64))

	exprs = getExpressions("(/ 10 4)")
	ret, _ = Interpret(exprs)
	assertNumber(t, 2.5, ret.(float64))

	exprs = getExpressions("(mod 10 4)")
	ret, _ = Interpret(exprs)
	assertNumber(t, 2, float64(ret.(int)))
}
