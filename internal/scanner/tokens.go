package scanner

//go:generate stringer -type=TokenType
type TokenType int

const (
	INVALID TokenType = iota
	EOF
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	STAR
	SLASH
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}
