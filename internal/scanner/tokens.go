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
	LESS
	MORE
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
