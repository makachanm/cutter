package lexer

type TokenType int

const (
	KEYWORD_CALL TokenType = iota + 0
	KEYWORD_DEFINE
	KEYWORD_BRACKET_OPEN
	KEYWORD_BRACKET_CLOSE

	STRING_QUOTEMARK
	BOOLEAN_TRUE
	BOOLEAN_FALSE

	WHITESPACE
	NEWLINE

	NORM_STRINGS
	VALUE
)

type Token struct {
	token_type TokenType
	token_data string
}

func NewToken(t_type TokenType) Token {
	return Token{token_type: t_type}
}

func NewDataToken(t_type TokenType, t_data string) Token {
	return Token{token_type: t_type, token_data: t_data}
}

func (t *Token) GetType() TokenType {
	return t.token_type
}

func (t *Token) GetData() string {
	return t.token_data
}

const (
	DATA_INT LexerTokenDataType = iota + 1
	DATA_REAL
	DATA_STR
	DATA_BOOL
	DATA_NORMSTRING
	DATA_OBJNAME

	NODEF
)

type LexerTokenDataType int

type LexerToken struct {
	Type TokenType
	Data LexerTokenData
}

func NewLexerToken(t_type TokenType, t_data LexerTokenData) LexerToken {
	return LexerToken{
		Type: t_type,
		Data: t_data,
	}
}

type LexerTokenData struct {
	Type LexerTokenDataType

	IntData     int64
	RealData    float64
	StrData     string
	BoolData    bool
	ObjNameData string
	NormData    string
}

func NewData() LexerTokenData {
	return LexerTokenData{Type: NODEF}
}

func NewIntData(d_data int64) LexerTokenData {
	return LexerTokenData{Type: DATA_INT, IntData: d_data}
}

func NewRealData(d_data float64) LexerTokenData {
	return LexerTokenData{Type: DATA_REAL, RealData: d_data}
}

func NewStrData(d_data string) LexerTokenData {
	return LexerTokenData{Type: DATA_STR, StrData: d_data}
}

func NewObjNameData(d_data string) LexerTokenData {
	return LexerTokenData{Type: DATA_OBJNAME, StrData: d_data}
}

func NewNormData(d_data string) LexerTokenData {
	return LexerTokenData{Type: DATA_NORMSTRING, NormData: d_data}
}

func NewBoolData(d_data bool) LexerTokenData {
	return LexerTokenData{Type: DATA_INT, BoolData: d_data}
}
