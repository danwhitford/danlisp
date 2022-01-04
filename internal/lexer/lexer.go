package lexer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shaftoe44/danlisp/internal/token"
)

type Lexer struct {
	length  int
	current int
	source  string
	line    int
}

func NewLexer(input string) Lexer {
	return Lexer{
		current: 0,
		length:  len(input),
		source:  input,
		line:    1,
	}
}

func (lexer *Lexer) GetTokens() ([]token.Token, error) {
	var tokens []token.Token
	for lexer.current < lexer.length {
		c := lexer.peek()
		if c == "(" {
			c = lexer.consume()
			r := token.Token{TokenType: token.LB, Lexeme: c, Line: 1}
			tokens = append(tokens, r)
		} else if c == ")" {
			c = lexer.consume()
			r := token.Token{TokenType: token.RB, Lexeme: c, Line: 1}
			tokens = append(tokens, r)
		} else if isDigit(c) {
			t, err := lexer.consumeNumber()
			if err != nil {
				return tokens, err
			}
			tokens = append(tokens, t)
		} else if c == "\"" {
			t, err := lexer.consumeString()
			if err != nil {
				return tokens, err
			}
			tokens = append(tokens, t)
		} else if isWhitespace(c) {
			lexer.current++
		} else {
			lexeme := lexer.consumeLexeme()
			if lexeme == "set" {
				tokens = append(tokens, token.Token{TokenType: token.SET, Lexeme: lexeme, Line: lexer.line})
			} else if lexeme == "if" {
				tokens = append(tokens, token.Token{TokenType: token.IF, Lexeme: lexeme, Line: lexer.line})
			} else if lexeme == "while" {
				tokens = append(tokens, token.Token{TokenType: token.WHILE, Lexeme: lexeme, Line: lexer.line})
			} else if lexeme == "defun" {
				tokens = append(tokens, token.Token{TokenType: token.DEFUN, Lexeme: lexeme, Line: lexer.line})
			} else {
				tokens = append(tokens, token.Token{TokenType: token.KEYWORD, Lexeme: lexeme, Line: lexer.line})
			}
		}
	}

	return tokens, nil
}

func (lexer *Lexer) peek() string {
	return lexer.source[lexer.current : lexer.current+1]
}

func (lexer *Lexer) consume() string {
	s := lexer.source[lexer.current : lexer.current+1]
	lexer.current++
	return s
}

func endsToken(c string) bool {
	token_enders := []string{"(", ")", "/n", "/t", " "}
	for _, cc := range token_enders {
		if c == cc {
			return true
		}
	}
	return false
}

func isWhitespace(c string) bool {
	whitespace := []string{" ", "\n", "\t"}
	for _, cc := range whitespace {
		if c == cc {
			return true
		}
	}
	return false
}

func (lexer *Lexer) consumeLexeme() string {
	var b strings.Builder
	var c string
	for lexer.current < lexer.length && !endsToken(lexer.peek()) {
		c = lexer.consume()
		b.WriteString(c)
	}
	return b.String()
}

func (lexer *Lexer) consumeNumber() (token.Token, error) {
	var b strings.Builder
	var c string
	for lexer.current < lexer.length && !endsToken(lexer.peek()) {
		c = lexer.consume()
		b.WriteString(c)
	}
	val, ok := strconv.ParseFloat(b.String(), 64)
	if ok != nil {
		return token.Token{}, fmt.Errorf("error while lexing on line %d. '%v' is not a number", lexer.line, b.String())
	}
	return token.Token{TokenType: token.LITERAL, Lexeme: b.String(), Value: val, Line: 1}, nil
}

func isDigit(c string) bool {
	numbers := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for _, cc := range numbers {
		if c == cc {
			return true
		}
	}
	return false
}

func (lexer *Lexer) consumeString() (token.Token, error) {
	var b strings.Builder
	var c string
	b.WriteString(lexer.consume()) // Consume the first quote
	for lexer.current < lexer.length && lexer.peek() != "\"" {
		if lexer.peek() == "\n" {
			return token.Token{}, fmt.Errorf("error while lexing on line %d. reached end of line in string '%v'", lexer.line, b.String())
		}
		c = lexer.consume()
		b.WriteString(c)
	}
	if lexer.current == lexer.length {
		return token.Token{}, fmt.Errorf("error while lexing on line %d. reached end of input in string '%v'", lexer.line, b.String())
	}
	b.WriteString(lexer.consume()) // Consume the final quote
	lexeme := b.String()
	val, _ := strconv.Unquote(lexeme)

	return token.Token{TokenType: token.LITERAL, Lexeme: lexeme, Value: val, Line: 1}, nil
}
