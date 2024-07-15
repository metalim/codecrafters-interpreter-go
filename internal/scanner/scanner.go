package scanner

type Scanner struct {
	source string
	line   int
	queue  queue
	Next   chan Token
}

func New(source string) *Scanner {
	return &Scanner{source: source, line: 1, Next: make(chan Token, 10)}
}

type queue int

const (
	queueEmpty queue = iota
	queueEqual
	queueBang
	queueLess
	queueGreater
	queueSlash
	queueComment
)

func (s *Scanner) ScanTokens() {
	for _, b := range s.source {
		if s.handleQueue(b) {
			continue
		}
		switch b {
		case '(':
			s.Next <- Token{Type: LEFT_PAREN, Lexeme: "(", Line: s.line}
		case ')':
			s.Next <- Token{Type: RIGHT_PAREN, Lexeme: ")", Line: s.line}
		case '{':
			s.Next <- Token{Type: LEFT_BRACE, Lexeme: "{", Line: s.line}
		case '}':
			s.Next <- Token{Type: RIGHT_BRACE, Lexeme: "}", Line: s.line}
		case ',':
			s.Next <- Token{Type: COMMA, Lexeme: ",", Line: s.line}
		case '.':
			s.Next <- Token{Type: DOT, Lexeme: ".", Line: s.line}
		case '-':
			s.Next <- Token{Type: MINUS, Lexeme: "-", Line: s.line}
		case '+':
			s.Next <- Token{Type: PLUS, Lexeme: "+", Line: s.line}
		case ';':
			s.Next <- Token{Type: SEMICOLON, Lexeme: ";", Line: s.line}
		case '*':
			s.Next <- Token{Type: STAR, Lexeme: "*", Line: s.line}
		case '/':
			s.queue = queueSlash
		case '<':
			s.queue = queueLess
		case '>':
			s.queue = queueGreater
		case '=':
			s.queue = queueEqual
		case '!':
			s.queue = queueBang
		case '\n':
			s.line++
		default:
			s.Next <- Token{Type: INVALID, Lexeme: string(b), Line: s.line}
		}
	}
	s.handleQueue(0)
	s.Next <- Token{Type: EOF, Lexeme: "", Line: s.line}
	close(s.Next)
}

func (s *Scanner) handleQueue(next rune) bool {
	switch s.queue {
	case queueEqual:
		if next == '=' {
			s.Next <- Token{Type: EQUAL_EQUAL, Lexeme: "==", Line: s.line}
			s.queue = queueEmpty
			return true
		}
		s.Next <- Token{Type: EQUAL, Lexeme: "=", Line: s.line}
		s.queue = queueEmpty
		return false

	case queueBang:
		if next == '=' {
			s.Next <- Token{Type: BANG_EQUAL, Lexeme: "!=", Line: s.line}
			s.queue = queueEmpty
			return true
		}
		s.Next <- Token{Type: BANG, Lexeme: "!", Line: s.line}
		s.queue = queueEmpty
		return false

	case queueLess:
		if next == '=' {
			s.Next <- Token{Type: LESS_EQUAL, Lexeme: "<=", Line: s.line}
			s.queue = queueEmpty
			return true
		}
		s.Next <- Token{Type: LESS, Lexeme: "<", Line: s.line}
		s.queue = queueEmpty
		return false

	case queueGreater:
		if next == '=' {
			s.Next <- Token{Type: GREATER_EQUAL, Lexeme: ">=", Line: s.line}
			s.queue = queueEmpty
			return true
		}
		s.Next <- Token{Type: GREATER, Lexeme: ">", Line: s.line}
		s.queue = queueEmpty
		return false

	case queueSlash:
		if next == '/' {
			s.queue = queueComment
			return true
		}
		s.Next <- Token{Type: SLASH, Lexeme: "/", Line: s.line}
		s.queue = queueEmpty
		return false

	case queueComment:
		if next == '\n' || next == 0 {
			s.queue = queueEmpty
		}
		return true
	}
	return false
}
