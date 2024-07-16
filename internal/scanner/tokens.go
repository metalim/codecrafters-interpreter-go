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

	STRING
	NUMBER

	IDENTIFIER

	AND
	CLASS
	ELSE
	FALSE
	FOR
	FUN
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
)

type Token struct {
	Type    TokenType
	Line    int
	Lexeme  string
	Literal string
	Error   *Error
}

type Error struct {
	Message string
	Line    int
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) LineNumber() int {
	return e.Line
}
