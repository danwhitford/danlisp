package parser

import (
	"whitford.io/danlisp/expr"
	"whitford.io/danlisp/token"
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
		if tokens[current].TokenType == token.LB {
			e := consumeSeq()
			exprs = append(exprs, e)
		} else {
			e := expr.Atom{Value: tokens[current].Value}
			exprs = append(exprs, e)
			current++
		}
	}

	return exprs, nil
}

func getExpression() expr.Expr {
	if source[current].TokenType == token.LB {
		return consumeSeq()
	} else if source[current].TokenType == token.KEYWORD {
		e := expr.Symbol{Name: source[current].Lexeme}
		current++
		return e
	} else {
		e := expr.Atom{Value: source[current].Value}
		current++
		return e
	}
}

func consumeSeq() expr.Seq {
	seq := []expr.Expr{}
	consume() // Consume the LB
	for current < length && peek().TokenType != token.RB {
		seq = append(seq, getExpression())
	}
	consume() // Consume the RB
	return expr.Seq{Exprs: seq}
}

func consume() token.Token {
	s := source[current]
	current++
	return s
}

func peek() token.Token {
	return source[current]
}