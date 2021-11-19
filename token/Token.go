package token

type TokenType int

const (
	LB = iota
	RB
	KEYWORD
	LITERAL
)

type Token struct {
	TokenType TokenType
	Lexeme    string
	Line      int
	Value     interface{}
}
