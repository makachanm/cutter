package lexer

type LexerStatus int

const (
	STATE_OBJECTPARSE LexerStatus = iota + 1
	STATE_NORMSTRINGS
	STATE_STRINGVALUE
	STATE_OBJNAME
)
