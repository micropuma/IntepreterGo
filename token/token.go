package token

type TokenType string

// TokenType is constant
const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	// Identifiers + literals
	IDENT = "IDENT"
	INT = "INT"

	// Operators
	ASSIGN = "=" 
	PLUS = "+" 
	MINUS = "-" 
	BANG = "!"
	ASTERISK = "*"
	SLASH = "/"

	LT = "<"
	GT = ">"

	EQ = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "(" 
	RPAREN = ")" 
	LBRACE = "{" 
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET = "LET"
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"
)

type Token struct {
	Type TokenType
	Literal string
}

// map of keywords
var keywords = map[string] TokenType {
	"fn" : FUNCTION,
	"let" : LET,
	"true" : TRUE,
	"false" : FALSE,
	"if" : IF,
	"else" : ELSE,
	"return" : RETURN,
} 

// LookupIdent : lookup identifier
func LookupIdent(input string) TokenType {
	if token, ok := keywords[input]; ok {
		return token
	}

	return IDENT 
}