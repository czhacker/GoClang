package token

type Tokentype string

type Token struct {
	Type Tokentype
	Literal string
}


const(
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT"
	INT = "INT"

	//Operators
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

	//Delimiters
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"


	//Keywords
	FUNCTION = "fn"
	LET = "LET"
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"
)

var keywords = map[string]Tokentype {
	"fn": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"if":IF,
	"else": ELSE,
	"return" : RETURN,
}

func LookupIdent(ident string) Tokentype {
	if tok, ok := keywords[ident]; ok{
		return tok
	}
	return IDENT
}

