package parser

import (
	"fmt"

	"github.com/shaftoe44/danlisp/internal/expr"
	"github.com/shaftoe44/danlisp/internal/token"
)

type Parser struct {
	current int
	length  int
	source  []token.Token
}

func NewParser(tokens []token.Token) Parser {
	return Parser{
		current: 0,
		length:  len(tokens),
		source:  tokens,
	}
}

func (parser *Parser) GetExpressions(tokens []token.Token) ([]expr.Expr, error) {
	exprs := []expr.Expr{}

	for parser.current < parser.length {
		expr, err := parser.getExpression()
		if err != nil {
			return exprs, err
		}
		exprs = append(exprs, expr)
	}

	return exprs, nil
}

func (parser *Parser) getExpression() (expr.Expr, error) {
	switch parser.peek().TokenType {
	case token.LB:
		if parser.next().TokenType == token.DEF {
			return parser.consumeDef()
		} else if parser.next().TokenType == token.IF {
			return parser.consumeIf()
		} else if parser.next().TokenType == token.WHILE {
			return parser.consumeWhile()
		} else {
			return parser.consumeSeq()
		}
	case token.KEYWORD:
		return parser.consumeKeyword()
	default:
		return parser.consumeAtom()
	}
}

func (parser *Parser) consumeWhile() (expr.While, error) {
	parser.consume() // Consume the LB
	parser.consume() // Consume the while

	cond, cerr := parser.getExpression()
	if cerr != nil {
		return expr.While{}, cerr
	}

	body := []expr.Expr{}
	for parser.current < parser.length && parser.peek().TokenType != token.RB {
		e, err := parser.getExpression()
		if err != nil {
			return expr.While{Cond: cond, Body: body}, err
		}
		body = append(body, e)
	}

	parser.consume() // Consume the RB
	return expr.While{Cond: cond, Body: body}, nil
}

func (parser *Parser) consumeKeyword() (expr.Symbol, error) {
	return expr.Symbol{Name: parser.consume().Lexeme}, nil
}

func (parser *Parser) consumeAtom() (expr.Atom, error) {
	e := expr.Atom{Value: parser.consume().Value}
	return e, nil
}

func (parser *Parser) consumeSeq() (expr.Seq, error) {
	seq := []expr.Expr{}

	parser.consume() // Consume the LB
	for parser.current < parser.length && parser.peek().TokenType != token.RB {
		e, err := parser.getExpression()
		if err != nil {
			return expr.Seq{Exprs: seq}, err
		}
		seq = append(seq, e)
	}
	if parser.current == parser.length {
		return expr.Seq{}, fmt.Errorf("parse error. missing ')' to close sequence")
	}
	parser.consume() // Consume the RB
	return expr.Seq{Exprs: seq}, nil
}

func (parser *Parser) consumeDef() (expr.Def, error) {
	parser.consume() // Consume the LB
	parser.consume() // Consume the def
	va := parser.consume()
	if va.TokenType != token.KEYWORD {
		return expr.Def{}, fmt.Errorf("parser error. trying to assign to '%v'", va.Lexeme)
	}
	sy := expr.Symbol{Name: va.Lexeme}
	val, _ := parser.getExpression()
	parser.consume() // Consume the RB
	return expr.Def{Var: sy, Value: val}, nil
}

func (parser *Parser) consumeIf() (expr.If, error) {
	parser.consume() // Consume the LB
	parser.consume() // Consume the if
	cond, err := parser.getExpression()
	if err != nil {
		return expr.If{}, err
	}
	trueBranch, err := parser.getExpression()
	if err != nil {
		return expr.If{}, err
	}
	falseBranch, err := parser.getExpression()
	if err != nil {
		return expr.If{}, err
	}
	parser.consume() // Consume the RB
	return expr.If{Cond: cond, TrueBranch: trueBranch, FalseBranch: falseBranch}, nil
}

func (parser *Parser) consume() token.Token {
	s := parser.source[parser.current]
	parser.current++
	return s
}

func (parser *Parser) peek() token.Token {
	return parser.source[parser.current]
}

func (parser *Parser) next() token.Token {
	return parser.source[parser.current+1]
}
