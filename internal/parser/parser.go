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
		if tokens[current].TokenType == token.LB {
			if next().TokenType == token.DEF {
				e, err := consumeDef()
				if err != nil {
					return exprs, err
				}
				exprs = append(exprs, e)
			} else {
				e, err := consumeSeq()
				if err != nil {
					return exprs, err
				}
				exprs = append(exprs, e)
			}
		} else {
			e := expr.Atom{Value: tokens[current].Value}
			exprs = append(exprs, e)
			current++
		}
	}

	return exprs, nil
}

func getExpression() (expr.Expr, error) {
	if source[current].TokenType == token.LB {
		return consumeSeq()
	} else if source[current].TokenType == token.KEYWORD {
		e := expr.Symbol{Name: source[current].Lexeme}
		current++
		return e, nil
	} else {
		e := expr.Atom{Value: source[current].Value}
		current++
		return e, nil
	}
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
