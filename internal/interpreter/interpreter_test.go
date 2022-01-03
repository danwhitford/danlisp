package interpreter

import (
	"testing"

	"github.com/shaftoe44/danlisp/internal/expr"
	"github.com/shaftoe44/danlisp/internal/lexer"
	"github.com/shaftoe44/danlisp/internal/parser"
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

func assert(t *testing.T, b bool) {
	if !b {
		t.Fatalf("Expected true but was false")
	}
}

func getExpressions(input string) []expr.Expr {
	lex := lexer.NewLexer(input)
	tokens, _ := lex.GetTokens()
	parser := parser.NewParser(tokens)
	exprs, _ := parser.GetExpressions(tokens)
	return exprs
}

func TestJustNumber(t *testing.T) {
	exprs := getExpressions("101")
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assertNumber(t, 101, ret.(float64))
}

func TestJustString(t *testing.T) {
	exprs := getExpressions("\"testing testing\"")
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assertString(t, "testing testing", ret.(string))
}

func TestAddTwo(t *testing.T) {
	exprs := getExpressions("(+ 2 7)")
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assertNumber(t, 9, ret.(float64))
}

func TestSubtract(t *testing.T) {
	exprs := getExpressions("(- 2 7)")
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assertNumber(t, -5, ret.(float64))
}

func TestErrorFuncNotFound(t *testing.T) {
	exprs := getExpressions("(nonsuch 2 7)")
	intr := NewInterpreter()
	_, err := intr.Interpret(exprs)
	assertString(t, "runtime error. Could not find symbol 'nonsuch'", err.Error())
}

func TestMoreBasicOperators(t *testing.T) {
	exprs := getExpressions("(* 2 7)")
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assertNumber(t, 14, ret.(float64))

	exprs = getExpressions("(/ 10 4)")
	ret, _ = intr.Interpret(exprs)
	assertNumber(t, 2.5, ret.(float64))

	exprs = getExpressions("(mod 10 4)")
	ret, _ = intr.Interpret(exprs)
	assertNumber(t, 2, ret.(float64))
}

func TestBitwiseOps(t *testing.T) {
	exprs := getExpressions("(& 255 101)")
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assertNumber(t, 101, ret.(float64))

	exprs = getExpressions("(| 255 72)")
	ret, _ = intr.Interpret(exprs)
	assertNumber(t, 255, ret.(float64))

	exprs = getExpressions("(^ 0 72)")
	ret, _ = intr.Interpret(exprs)
	assertNumber(t, 72, ret.(float64))

	exprs = getExpressions("(&^ 255 72)")
	ret, _ = intr.Interpret(exprs)
	assertNumber(t, 183, ret.(float64))
}

func TestShifts(t *testing.T) {
	exprs := getExpressions("(>> 255 2)")
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assertNumber(t, 63, ret.(float64))

	exprs = getExpressions("(<< 255 2)")
	ret, _ = intr.Interpret(exprs)
	assertNumber(t, 1020, ret.(float64))
}

func TestDefinition(t *testing.T) {
	exprs := getExpressions("(def foo 10) (* foo 5)")
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assertNumber(t, 50, ret.(float64))
}

func TestEquals(t *testing.T) {
	exprs := getExpressions("(def foo 10) (= foo 10)")
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assert(t, ret.(bool))

	exprs = getExpressions("(def foo 10) (= foo 5)")
	ret, _ = intr.Interpret(exprs)
	assert(t, !ret.(bool))

	exprs = getExpressions("(def bar \"dan\") (= bar \"dan\")")
	ret, _ = intr.Interpret(exprs)
	assert(t, ret.(bool))

	exprs = getExpressions("(def bar \"dan\") (= foo \"egg\")")
	ret, _ = intr.Interpret(exprs)
	assert(t, !ret.(bool))
}

func TestAndOr(t *testing.T) {
	exprs := getExpressions("(and (= 2 2) (= 5 5))")
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assert(t, ret.(bool))

	exprs = getExpressions("(and (= 2 2) (= 1 5))")
	ret, _ = intr.Interpret(exprs)
	assert(t, !ret.(bool))

	exprs = getExpressions("(or (= 2 2) (= 1 5))")
	ret, _ = intr.Interpret(exprs)
	assert(t, ret.(bool))

	exprs = getExpressions("(or (= 2 5) (= 1 5))")
	ret, _ = intr.Interpret(exprs)
	assert(t, !ret.(bool))
}

func TestIfExpr(t *testing.T) {
	exprs := getExpressions(`(if (= 2 2) "yes" "no")`)
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assertString(t, "yes", ret.(string))

	exprs = getExpressions(`(if (= (+2 2) 5) "yes" "no")`)
	ret, _ = intr.Interpret(exprs)
	assertString(t, "no", ret.(string))
}

func TestWhileExpr(t *testing.T) {
	exprs := getExpressions(`(def x 5) (def total 0) (while (gt x 0) (def total (+ total x)) (def x (- x 1))) total`)
	intr := NewInterpreter()
	ret, _ := intr.Interpret(exprs)
	assertNumber(t, 15, ret.(float64))
}
