package lexer

type LexerStatus int

const (
	STATE_OBJECTPARSE LexerStatus = iota + 0
	STATE_NORMSTRINGS
	STATE_STRINGVALUE
	STATE_OBJNAME
)
