package lexer

import (
	"fmt"
	"strconv"
	"strings"
	"whitford.io/danlisp/internal/token"
)

var length int
var current int
var source string
var line int

func GetTokens(input string) ([]token.Token, error) {
	current = 0
	length = len(input)
	source = input
	line = 1

	var tokens []token.Token
	for current < length {
		c := peek()
		if c == "(" {
			c = consume()
			r := token.Token{TokenType: token.LB, Lexeme: c, Line: 1}
			tokens = append(tokens, r)
		} else if c == ")" {
			c = consume()
			r := token.Token{TokenType: token.RB, Lexeme: c, Line: 1}
			tokens = append(tokens, r)
		} else if isDigit(c) {
			t, err := consumeNumber()
			if err != nil {
				return tokens, err
			}
			tokens = append(tokens, t)
		} else if c == "\"" {
			t, err := consumeString()
			if err != nil {
				return tokens, err
			}
			tokens = append(tokens, t)
		} else if isWhitespace(c) {
			current++
		} else {
			lexeme := consumeLexeme()
			if lexeme == "def" {
				tokens = append(tokens, token.Token{TokenType: token.DEF, Lexeme: lexeme, Line: line})
			} else {
				tokens = append(tokens, token.Token{TokenType: token.KEYWORD, Lexeme: lexeme, Line: line})
			}
		}
	}

	return tokens, nil
}

func peek() string {
	return source[current : current+1]
}

func consume() string {
	s := source[current : current+1]
	current++
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

func consumeLexeme() string {
	var b strings.Builder
	var c string
	for current < length && !endsToken(peek()) {
		c = consume()
		b.WriteString(c)
	}
	return b.String()
}

// func consumeKeyword() token.Token {
// 	var b strings.Builder
// 	var c string
// 	for current < length && !endsToken(peek()) {
// 		c = consume()
// 		b.WriteString(c)
// 	}
// 	return token.Token{TokenType: token.KEYWORD, Lexeme: b.String(), Line: 1}
// }

func consumeNumber() (token.Token, error) {
	var b strings.Builder
	var c string
	for current < length && !endsToken(peek()) {
		c = consume()
		b.WriteString(c)
	}
	val, ok := strconv.ParseFloat(b.String(), 64)
	if ok != nil {
		return token.Token{}, fmt.Errorf("error while lexing on line %d. '%v' is not a number", line, b.String())
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

func consumeString() (token.Token, error) {
	var b strings.Builder
	var c string
	b.WriteString(consume()) // Consume the first quote
	for current < length && peek() != "\"" {
		if peek() == "\n" {
			return token.Token{}, fmt.Errorf("error while lexing on line %d. reached end of line in string '%v'", line, b.String())
		}
		c = consume()
		b.WriteString(c)
	}
	if current == length {
		return token.Token{}, fmt.Errorf("error while lexing on line %d. reached end of input in string '%v'", line, b.String())
	}
	b.WriteString(consume()) // Consume the final quote
	lexeme := b.String()
	val, _ := strconv.Unquote(lexeme)

	return token.Token{TokenType: token.LITERAL, Lexeme: lexeme, Value: val, Line: 1}, nil
}