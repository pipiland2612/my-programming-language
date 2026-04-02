package token

type TokenType int

const (
	// Literals
	INT TokenType = iota
	STRING
	IDENT
	TRUE
	FALSE

	// Operators
	PLUS
	MINUS
	STAR
	SLASH
	PERCENT
	EQ
	EQEQ
	NEQ
	LT
	GT
	LEQ
	GEQ
	AND
	OR
	NOT
	ARROW  // ->
	FATARROW // =>
	DCOLON // ::

	// Delimiters
	LPAREN
	RPAREN
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE
	COMMA
	COLON
	SEMICOLON
	PIPE
	DOT
	BACKSLASH

	// Keywords
	LET
	IN
	IF
	THEN
	ELSE
	FN
	FIX
	FST
	SND
	INL
	INR
	CASE
	OF
	IMPORT
	PRINT
	PRINTLN

	// Type keywords
	INT_TYPE
	BOOL_TYPE
	STRING_TYPE
	UNIT_TYPE

	// Special
	EOF
	ILLEGAL
)

var tokenNames = map[TokenType]string{
	INT: "INT", STRING: "STRING", IDENT: "IDENT",
	TRUE: "true", FALSE: "false",
	PLUS: "+", MINUS: "-", STAR: "*", SLASH: "/", PERCENT: "%",
	EQ: "=", EQEQ: "==", NEQ: "!=",
	LT: "<", GT: ">", LEQ: "<=", GEQ: ">=",
	AND: "&&", OR: "||", NOT: "!",
	ARROW: "->", FATARROW: "=>", DCOLON: "::",
	LPAREN: "(", RPAREN: ")", LBRACKET: "[", RBRACKET: "]",
	LBRACE: "{", RBRACE: "}",
	COMMA: ",", COLON: ":", SEMICOLON: ";", PIPE: "|", DOT: ".",
	BACKSLASH: "\\",
	LET: "let", IN: "in", IF: "if", THEN: "then", ELSE: "else",
	FN: "fn", FIX: "fix", FST: "fst", SND: "snd",
	INL: "inl", INR: "inr", CASE: "case", OF: "of",
	IMPORT: "import", PRINT: "print", PRINTLN: "println",
	INT_TYPE: "Int", BOOL_TYPE: "Bool", STRING_TYPE: "String", UNIT_TYPE: "Unit",
	EOF: "EOF", ILLEGAL: "ILLEGAL",
}

func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return "UNKNOWN"
}

var Keywords = map[string]TokenType{
	"let":     LET,
	"in":      IN,
	"if":      IF,
	"then":    THEN,
	"else":    ELSE,
	"fn":      FN,
	"fix":     FIX,
	"fst":     FST,
	"snd":     SND,
	"inl":     INL,
	"inr":     INR,
	"case":    CASE,
	"of":      OF,
	"true":    TRUE,
	"false":   FALSE,
	"import":  IMPORT,
	"print":   PRINT,
	"println": PRINTLN,
	"Int":     INT_TYPE,
	"Bool":    BOOL_TYPE,
	"String":  STRING_TYPE,
	"Unit":    UNIT_TYPE,
}

type Pos struct {
	Line   int
	Column int
}

type Token struct {
	Type    TokenType
	Literal string
	Pos     Pos
}
