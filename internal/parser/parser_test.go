package parser

import (
	"testing"

	"github.com/shaftoe44/danlisp/internal/expr"
	"github.com/shaftoe44/danlisp/internal/lexer"
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
	lex := lexer.NewLexer(input)
	tokens, _ := lex.GetTokens()
	parser := NewParser(tokens)
	expressions, _ := parser.GetExpressions(tokens)
	assertNumber(t, 123.7, expressions[0].(expr.Atom).Value.(float64))
}

func TestString(t *testing.T) {
	input := "\"egghead\""
	lex := lexer.NewLexer(input)
	tokens, _ := lex.GetTokens()
	parser := NewParser(tokens)
	expressions, _ := parser.GetExpressions(tokens)
	assertString(t, "egghead", expressions[0].(expr.Atom).Value.(string))
}

func TestSeq(t *testing.T) {
	input := "(0 1 2 3 4 5)"
	lex := lexer.NewLexer(input)
	tokens, _ := lex.GetTokens()
	parser := NewParser(tokens)
	expressions, _ := parser.GetExpressions(tokens)
	seq := expressions[0].(expr.Seq)
	for i := 0; i < 6; i++ {
		val := seq.Exprs[i].(expr.Atom).Value.(float64)
		assertNumber(t, float64(i), val)
	}
}

func TestNestedSeq(t *testing.T) {
	input := "(concat (0 1 2) (3 4 5))"
	lex := lexer.NewLexer(input)
	tokens, _ := lex.GetTokens()
	parser := NewParser(tokens)
	expressions, _ := parser.GetExpressions(tokens)
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
	lex := lexer.NewLexer(input)
	tokens, _ := lex.GetTokens()
	parser := NewParser(tokens)
	_, err := parser.GetExpressions(tokens)
	assertString(t, "parse error. missing ')' to close sequence", err.Error())
}

func TestDefinition(t *testing.T) {
	input := "(def x 5)"
	lex := lexer.NewLexer(input)
	tokens, _ := lex.GetTokens()
	parser := NewParser(tokens)
	exprs, _ := parser.GetExpressions(tokens)
	defe, ok := exprs[0].(expr.Def)
	if !ok {
		t.Fatalf("Conversion to Def expression failed")
	}
	assertString(t, defe.Var.Name, "x")
	assertNumber(t, 5, defe.Value.(expr.Atom).Value.(float64))
}

func TestIf(t *testing.T) {
	input := `(if (= 2 2) "yes" "no")`
	lex := lexer.NewLexer(input)
	tokens, _ := lex.GetTokens()
	parser := NewParser(tokens)
	exprs, _ := parser.GetExpressions(tokens)
	ife, ok := exprs[0].(expr.If)
	if !ok {
		t.Fatalf("Conversion to If expression failed")
	}
	if ife.Cond.(expr.Seq).Exprs[0].(expr.Symbol).Name != "=" {
		t.Fatal("Condition wasn't right")
	}
	if ife.TrueBranch.(expr.Atom).Value != "yes" {
		t.Fatal("True branch wasn't right")
	}
	if ife.FalseBranch.(expr.Atom).Value != "no" {
		t.Fatal("False branch wasn't right")
	}
}

func TestWhile(t *testing.T) {
	input := `(while (> count 0) (def count (- count 1)))`
	lex := lexer.NewLexer(input)
	tokens, _ := lex.GetTokens()
	parser := NewParser(tokens)
	exprs, _ := parser.GetExpressions(tokens)
	ife, ok := exprs[0].(expr.While)
	if !ok {
		t.Fatalf("Conversion to While expression failed")
	}
	if ife.Cond.(expr.Seq).Exprs[0].(expr.Symbol).Name != ">" {
		t.Fatal("Condition wasn't right")
	}
}
