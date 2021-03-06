package parser

import (
	"fmt"

	"github.com/danwhitford/danlisp/internal/expr"
	"github.com/danwhitford/danlisp/internal/token"
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

func (parser *Parser) GetExpressions() ([]expr.Expr, error) {
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

// TODO make this switch
func (parser *Parser) getExpression() (expr.Expr, error) {
	switch parser.peek().TokenType {
	case token.LB:
		if parser.next().TokenType == token.SET {
			return parser.consumeSet()
		} else if parser.next().TokenType == token.IF {
			return parser.consumeIf()
		} else if parser.next().TokenType == token.WHILE {
			return parser.consumeWhile()
		} else if parser.next().TokenType == token.DEFN {
			return parser.consumeDefun()
		} else if parser.next().TokenType == token.FOR {
			return parser.consumeFor()
		} else {
			return parser.consumeSeq()
		}
	case token.KEYWORD:
		return parser.consumeKeyword()
	default:
		return parser.consumeAtom()
	}
}

func (parser *Parser) consumeDefun() (expr.Defn, error) {
	// TODO Make a consumeExpected func
	parser.consume() // Consume the LB
	parser.consume() // Consume the defun

	fnName := parser.consume()
	if fnName.TokenType != token.KEYWORD {
		return expr.Defn{}, fmt.Errorf("expected keyword but got %v", fnName)
	}
	fnSymb := expr.Symbol{Name: fnName.Lexeme}

	parser.consume() // Consume the LB for arglist
	argList := []expr.Symbol{}
	for parser.current < parser.length && parser.peek().TokenType != token.RB {
		a := parser.consume()
		if a.TokenType != token.KEYWORD {
			return expr.Defn{}, fmt.Errorf("arguments must be symbols but got %v", a)
		}
		argList = append(argList, expr.Symbol{Name: a.Lexeme})
	}
	parser.consume() // Consume the RB after arglist

	body := []expr.Expr{}
	for parser.current < parser.length && parser.peek().TokenType != token.RB {
		e, err := parser.getExpression()
		if err != nil {
			return expr.Defn{}, err
		}
		body = append(body, e)
	}
	parser.consume() // Consume the RB after function body

	return expr.Defn{Name: fnSymb, Arglist: argList, Body: body}, nil
}

func (parser *Parser) consumeFor() (expr.For, error) {
	parser.consume() // Consume the LB
	parser.consume() // Consume the for

	init, err := parser.getExpression()
	if err != nil {
		return expr.For{}, err
	}
	cond, err := parser.getExpression()
	if err != nil {
		return expr.For{}, err
	}
	step, err := parser.getExpression()
	if err != nil {
		return expr.For{}, err
	}
	body := []expr.Expr{}
	for parser.current < parser.length && parser.peek().TokenType != token.RB {
		e, err := parser.getExpression()
		if err != nil {
			return expr.For{}, err
		}
		body = append(body, e)
	}
	parser.consume() // Consume the RB

	return expr.For{Initialiser: init, Cond: cond, Step: step, Body: body}, nil
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

func (parser *Parser) consumeSet() (expr.Set, error) {
	parser.consume() // Consume the LB
	parser.consume() // Consume the set
	va := parser.consume()
	if va.TokenType != token.KEYWORD {
		return expr.Set{}, fmt.Errorf("parser error. trying to assign to '%v'", va.Lexeme)
	}
	sy := expr.Symbol{Name: va.Lexeme}
	val, _ := parser.getExpression()
	parser.consume() // Consume the RB
	return expr.Set{Var: sy, Value: val}, nil
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
