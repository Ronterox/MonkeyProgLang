package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	STRING   = "STRING"
	TEMPLATE = "TEMPLATE"
	IDENT    = "IDENT"
	INT      = "INT"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	PERCENT  = "%"
	BANG     = "!"
	COLON    = ":"

	LT  = "<"
	GT  = ">"
	LE  = "<="
	GE  = ">="
	EQ  = "=="
	NE  = "!="
	AND = "&"
	OR  = "|"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"
	MACRO    = "MACRO"
	RETURN   = "RETURN"
	FALSE    = "FALSE"
	TRUE     = "TRUE"
	ELSE     = "ELSE"
	IF       = "IF"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"macro":  MACRO,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
