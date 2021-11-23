package parser

import (
	"fmt"

	"whitford.io/danlisp/internal/expr"
	"whitford.io/danlisp/internal/token"
)

var current int
var length int
var source []token.Token

func GetExpressions(tokens []token.Token) ([]expr.Expr, error) {
	current = 0
	length = len(tokens)
	source = tokens

	exprs := []expr.Expr{}

	for current < length {
		expr, err := getExpression()
		if err != nil {
			return exprs, err
		}
		exprs = append(exprs, expr)
	}

	return exprs, nil
}

func getExpression() (expr.Expr, error) {
	switch peek().TokenType {
	case token.LB:
		if next().TokenType == token.DEF {
			return consumeDef()
		} else if next().TokenType == token.IF {
			return consumeIf()
		} else {
			return consumeSeq()
		}
	case token.KEYWORD:
		return consumeKeyword()
	default:
		return consumeAtom()
	}
}

func consumeKeyword() (expr.Symbol, error) {
	return expr.Symbol{Name: consume().Lexeme}, nil
}

func consumeAtom() (expr.Atom, error) {
	e := expr.Atom{Value: consume().Value}
	return e, nil
}

func consumeSeq() (expr.Seq, error) {
	seq := []expr.Expr{}

	consume() // Consume the LB
	for current < length && peek().TokenType != token.RB {
		e, err := getExpression()
		if err != nil {
			return expr.Seq{Exprs: seq}, err
		}
		seq = append(seq, e)
	}
	if current == length {
		return expr.Seq{}, fmt.Errorf("parse error. missing ')' to close sequence.")
	}
	consume() // Consume the RB
	return expr.Seq{Exprs: seq}, nil
}

func consumeDef() (expr.Def, error) {
	consume() // Consume the LB
	consume() // Consume the def
	va := consume()
	if va.TokenType != token.KEYWORD {
		return expr.Def{}, fmt.Errorf("parser error. trying to assign to '%v'", va.Lexeme)
	}
	sy := expr.Symbol{Name: va.Lexeme}
	val, _ := getExpression()
	consume() // Consume the RB
	return expr.Def{Var: sy, Value: val}, nil
}

func consumeIf() (expr.If, error) {
	consume() // Consume the LB
	consume() // Consume the if
	cond, err := getExpression()
	if err != nil {
		return expr.If{}, err
	}
	trueBranch, err := getExpression()
	if err != nil {
		return expr.If{}, err
	}
	falseBranch, err := getExpression()
	if err != nil {
		return expr.If{}, err
	}
	consume() // Consume the RB
	return expr.If{Cond: cond, TrueBranch: trueBranch, FalseBranch: falseBranch}, nil
}

func consume() token.Token {
	s := source[current]
	current++
	return s
}

func peek() token.Token {
	return source[current]
}

func next() token.Token {
	return source[current+1]
}
