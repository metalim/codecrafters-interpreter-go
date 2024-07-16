package scanner

import "unicode/utf8"

type Scanner struct {
	source           string
	emitPos, scanPos int
	line             int
	Next             chan *Token
}

func New(source string) *Scanner {
	return &Scanner{source: source, line: 1, Next: make(chan *Token, 10)}
}

func (s *Scanner) ScanTokens() {
	for r := s.advance(); r != 0; r = s.advance() {
		switch r {
		case '(':
			s.emit(LEFT_PAREN)
		case ')':
			s.emit(RIGHT_PAREN)
		case '{':
			s.emit(LEFT_BRACE)
		case '}':
			s.emit(RIGHT_BRACE)
		case ',':
			s.emit(COMMA)
		case '.':
			s.emit(DOT)
		case '-':
			s.emit(MINUS)
		case '+':
			s.emit(PLUS)
		case ';':
			s.emit(SEMICOLON)
		case '*':
			s.emit(STAR)
		case '/':
			if s.match('/') {
				for r = s.advance(); r != '\n' && r != 0; r = s.advance() {
				}
				if r == '\n' {
					s.line++
				}
				s.eat()
			} else {
				s.emit(SLASH)
			}
		case '<':
			if s.match('=') {
				s.emit(LESS_EQUAL)
			} else {
				s.emit(LESS)
			}
		case '>':
			if s.match('=') {
				s.emit(GREATER_EQUAL)
			} else {
				s.emit(GREATER)
			}
		case '=':
			if s.match('=') {
				s.emit(EQUAL_EQUAL)
			} else {
				s.emit(EQUAL)
			}
		case '!':
			if s.match('=') {
				s.emit(BANG_EQUAL)
			} else {
				s.emit(BANG)
			}
		case '\n':
			s.line++
			s.eat()
		case ' ', '\r', '\t':
			s.eat()
		case '"':
			for r = s.advance(); r != '"' && r != 0; r = s.advance() {
				if r == '\n' {
					s.line++
				}
			}
			if r == 0 {
				s.emitError("Unterminated string.")
			} else {
				s.emitLiteral(STRING, s.source[s.emitPos+1:s.scanPos-1])
			}
		default:
			switch {
			case isDigit(r):
				s.advanceDigits()
				if s.peek() == '.' && isDigit(s.peekNext()) {
					s.advance()
					s.advanceDigits()
					// truncate excess zeroes
					lastNonZeroPos := s.scanPos - 1
					for s.source[lastNonZeroPos] == '0' {
						lastNonZeroPos--
					}
					if s.source[lastNonZeroPos] == '.' {
						lastNonZeroPos++
					}
					s.emitLiteral(NUMBER, s.source[s.emitPos:lastNonZeroPos+1])
				} else {
					// doh... why? Java nerds
					s.emitLiteral(NUMBER, s.source[s.emitPos:s.scanPos]+".0")
				}
			case isAlpha(r):
				for r = s.peek(); isAlphaNumeric(r); r = s.peek() {
					s.advance()
				}
				if keyword, ok := keywords[s.source[s.emitPos:s.scanPos]]; ok {
					s.emit(keyword)
				} else {
					s.emit(IDENTIFIER)
				}
			default:
				s.emitError("Unexpected character: " + string(r))
			}
		}
	}
	s.emit(EOF)
	close(s.Next)
}

func (s *Scanner) peekSize() (r rune, size int) {
	if len(s.source[s.scanPos:]) == 0 {
		return 0, 0
	}
	return utf8.DecodeRuneInString(s.source[s.scanPos:])
}

func (s *Scanner) peek() rune {
	if len(s.source[s.scanPos:]) == 0 {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(s.source[s.scanPos:])
	return r
}

func (s *Scanner) peekNext() rune {
	if len(s.source[s.scanPos:]) == 0 {
		return 0
	}
	_, size := utf8.DecodeRuneInString(s.source[s.scanPos:])
	r, _ := utf8.DecodeRuneInString(s.source[s.scanPos+size:])
	return r
}

func (s *Scanner) advance() rune {
	r, size := s.peekSize()
	s.scanPos += size
	return r
}

func (s *Scanner) advanceDigits() {
	for r := s.peek(); r >= '0' && r <= '9'; r = s.peek() {
		s.advance()
	}
}

func (s *Scanner) match(expected rune) bool {
	peek, size := s.peekSize()
	if peek != expected {
		return false
	}
	s.scanPos += size
	return true
}

func (s *Scanner) eat() {
	s.emitPos = s.scanPos
}

func (s *Scanner) emit(t TokenType) {
	s.Next <- &Token{Type: t, Lexeme: s.source[s.emitPos:s.scanPos], Line: s.line}
	s.eat()
}

func (s *Scanner) emitLiteral(t TokenType, literal string) {
	s.Next <- &Token{Type: t, Lexeme: s.source[s.emitPos:s.scanPos], Line: s.line, Literal: literal}
	s.eat()
}

func (s *Scanner) emitError(message string) {
	s.Next <- &Token{Line: s.line, Error: &Error{Message: message, Line: s.line}}
	s.eat()
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_'
}

func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || isDigit(r)
}
