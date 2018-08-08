package token

type TokenType string

// Possible types of tokens
const (
	// Special token types
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Function/variable names & values
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// Small, easily categorizable data structures
type Token struct {
	Type    TokenType // Type of token
	Literal string    // Literal value of token
}

// Special identifiers
var specialIdentifiers = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func GetIdentifier(input string) TokenType {
	tokenType, ok := specialIdentifiers[input]
	if ok {
		return tokenType
	} else {
		return IDENT
	}
}
