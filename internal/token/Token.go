package token

type TokenType int

const (
	LB = iota
	RB
	KEYWORD
	LITERAL
	SET
	IF
	WHILE
)

type Token struct {
	TokenType TokenType
	Lexeme    string
	Line      int
	Value     interface{}
}
