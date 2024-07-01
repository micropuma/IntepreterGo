package lexer

import "monkey/token"

type Lexer struct {
	input 	  string
	position  int  // current position in input (points to current char)
	readPosition int // current reading position in input (after current char)
	ch        byte // current char under examination
} 

// that means New is a construction funciton that can construct a Lexer pointer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// that means reachar is a method of Lexer struct
func (l* Lexer)readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

// that means NextToken is a method of Lexer struct
func (l* Lexer)NextToken() token.Token{
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
		// for '=' and '!' branch, need to do some extension
		case '=':
			// 前瞻
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
			} else {
				tok = newToken(token.ASSIGN, l.ch)
			}
		case ';':
			tok = newToken(token.SEMICOLON, l.ch)
		case '(':
			tok = newToken(token.LPAREN, l.ch)
		case ')':
			tok = newToken(token.RPAREN, l.ch)
		case ',':
			tok = newToken(token.COMMA, l.ch)
		case '+':
			tok = newToken(token.PLUS, l.ch)
		case '{':
			tok = newToken(token.LBRACE, l.ch)
		case '}':
			tok = newToken(token.RBRACE, l.ch)
		case '-':
			tok = newToken(token.MINUS, l.ch)
		case '!':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
			} else { 
				tok = newToken(token.BANG, l.ch)
			}
		case '/':
			tok = newToken(token.SLASH, l.ch)
		case '*':
			tok = newToken(token.ASTERISK, l.ch)
		case '<':
			tok = newToken(token.LT, l.ch)
		case '>':
			tok = newToken(token.GT, l.ch)
		case 0:
			tok.Literal = ""
			tok.Type = token.EOF
		default:
			if isLetter(l.ch) {
				// get the content of token and the type of token
				tok.Literal = l.readIdentifier()
				tok.Type = token.LookupIdent(tok.Literal)
				return tok
			} else if isDigit(l.ch) {
				tok.Type = token.INT
				tok.Literal = l.readNumber()
				return tok
			} else {
				tok = newToken(token.ILLEGAL, l.ch)
			}	
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// read identifier
// get the start of identifier and the end of the identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// read digit
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position : l.position]
}

// check if the ch byte is a letter or not 
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_';
}

// check if the ch byte is a digit or not
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// skip the blank space 
func (l *Lexer)skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// specially deal with the "==" | "!="
func (l *Lexer)peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
