package lexer

import (
	"fmt"
	"my-programming-language/token"
	"unicode"
)

type Lexer struct {
	input   []rune
	pos     int
	line    int
	col     int
	tokens  []token.Token
}

func New(input string) *Lexer {
	return &Lexer{
		input: []rune(input),
		pos:   0,
		line:  1,
		col:   1,
	}
}

func (l *Lexer) Tokenize() ([]token.Token, error) {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		// Skip whitespace
		if unicode.IsSpace(ch) {
			l.advance()
			continue
		}

		// Single-line comment
		if ch == '/' && l.peek() == '/' {
			l.skipLineComment()
			continue
		}

		// Multi-line comment
		if ch == '/' && l.peek() == '*' {
			if err := l.skipBlockComment(); err != nil {
				return nil, err
			}
			continue
		}

		// Numbers
		if unicode.IsDigit(ch) {
			l.readNumber()
			continue
		}

		// Strings
		if ch == '"' {
			if err := l.readString(); err != nil {
				return nil, err
			}
			continue
		}

		// Identifiers and keywords
		if unicode.IsLetter(ch) || ch == '_' {
			l.readIdent()
			continue
		}

		// Operators and punctuation
		if err := l.readOperator(); err != nil {
			return nil, err
		}
	}

	l.tokens = append(l.tokens, token.Token{Type: token.EOF, Literal: "", Pos: l.curPos()})
	return l.tokens, nil
}

func (l *Lexer) curPos() token.Pos {
	return token.Pos{Line: l.line, Column: l.col}
}

func (l *Lexer) advance() rune {
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) peek() rune {
	if l.pos+1 < len(l.input) {
		return l.input[l.pos+1]
	}
	return 0
}

func (l *Lexer) emit(typ token.TokenType, lit string, pos token.Pos) {
	l.tokens = append(l.tokens, token.Token{Type: typ, Literal: lit, Pos: pos})
}

func (l *Lexer) skipLineComment() {
	for l.pos < len(l.input) && l.input[l.pos] != '\n' {
		l.advance()
	}
}

func (l *Lexer) skipBlockComment() error {
	pos := l.curPos()
	l.advance() // /
	l.advance() // *
	for l.pos < len(l.input) {
		if l.input[l.pos] == '*' && l.peek() == '/' {
			l.advance()
			l.advance()
			return nil
		}
		l.advance()
	}
	return fmt.Errorf("%d:%d: unterminated block comment", pos.Line, pos.Column)
}

func (l *Lexer) readNumber() {
	pos := l.curPos()
	start := l.pos
	for l.pos < len(l.input) && unicode.IsDigit(l.input[l.pos]) {
		l.advance()
	}
	l.emit(token.INT, string(l.input[start:l.pos]), pos)
}

func (l *Lexer) readString() error {
	pos := l.curPos()
	l.advance() // opening "
	start := l.pos
	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		if l.input[l.pos] == '\\' {
			l.advance() // skip escape char
		}
		if l.pos < len(l.input) {
			l.advance()
		}
	}
	if l.pos >= len(l.input) {
		return fmt.Errorf("%d:%d: unterminated string", pos.Line, pos.Column)
	}
	lit := string(l.input[start:l.pos])
	l.advance() // closing "
	l.emit(token.STRING, lit, pos)
	return nil
}

func (l *Lexer) readIdent() {
	pos := l.curPos()
	start := l.pos
	for l.pos < len(l.input) && (unicode.IsLetter(l.input[l.pos]) || unicode.IsDigit(l.input[l.pos]) || l.input[l.pos] == '_') {
		l.advance()
	}
	lit := string(l.input[start:l.pos])
	if typ, ok := token.Keywords[lit]; ok {
		l.emit(typ, lit, pos)
	} else {
		l.emit(token.IDENT, lit, pos)
	}
}

func (l *Lexer) readOperator() error {
	pos := l.curPos()
	ch := l.advance()

	switch ch {
	case '+':
		l.emit(token.PLUS, "+", pos)
	case '*':
		l.emit(token.STAR, "*", pos)
	case '/':
		l.emit(token.SLASH, "/", pos)
	case '%':
		l.emit(token.PERCENT, "%", pos)
	case '(':
		l.emit(token.LPAREN, "(", pos)
	case ')':
		l.emit(token.RPAREN, ")", pos)
	case '[':
		l.emit(token.LBRACKET, "[", pos)
	case ']':
		l.emit(token.RBRACKET, "]", pos)
	case '{':
		l.emit(token.LBRACE, "{", pos)
	case '}':
		l.emit(token.RBRACE, "}", pos)
	case ',':
		l.emit(token.COMMA, ",", pos)
	case ';':
		l.emit(token.SEMICOLON, ";", pos)
	case '\\':
		l.emit(token.BACKSLASH, "\\", pos)
	case '.':
		l.emit(token.DOT, ".", pos)
	case '-':
		if l.pos < len(l.input) && l.input[l.pos] == '>' {
			l.advance()
			l.emit(token.ARROW, "->", pos)
		} else {
			l.emit(token.MINUS, "-", pos)
		}
	case '=':
		if l.pos < len(l.input) && l.input[l.pos] == '>' {
			l.advance()
			l.emit(token.FATARROW, "=>", pos)
		} else if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.advance()
			l.emit(token.EQEQ, "==", pos)
		} else {
			l.emit(token.EQ, "=", pos)
		}
	case '!':
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.advance()
			l.emit(token.NEQ, "!=", pos)
		} else {
			l.emit(token.NOT, "!", pos)
		}
	case '<':
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.advance()
			l.emit(token.LEQ, "<=", pos)
		} else {
			l.emit(token.LT, "<", pos)
		}
	case '>':
		if l.pos < len(l.input) && l.input[l.pos] == '=' {
			l.advance()
			l.emit(token.GEQ, ">=", pos)
		} else {
			l.emit(token.GT, ">", pos)
		}
	case '&':
		if l.pos < len(l.input) && l.input[l.pos] == '&' {
			l.advance()
			l.emit(token.AND, "&&", pos)
		} else {
			return fmt.Errorf("%d:%d: unexpected character '&'", pos.Line, pos.Column)
		}
	case '|':
		if l.pos < len(l.input) && l.input[l.pos] == '|' {
			l.advance()
			l.emit(token.OR, "||", pos)
		} else {
			l.emit(token.PIPE, "|", pos)
		}
	case ':':
		if l.pos < len(l.input) && l.input[l.pos] == ':' {
			l.advance()
			l.emit(token.DCOLON, "::", pos)
		} else {
			l.emit(token.COLON, ":", pos)
		}
	default:
		return fmt.Errorf("%d:%d: unexpected character '%c'", pos.Line, pos.Column, ch)
	}
	return nil
}
