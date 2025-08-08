package lexer

type KeywordMatchingItem map[string]TokenType
type InvertedKeywordMatchingItem map[TokenType]string

var KeywordMap KeywordMatchingItem = KeywordMatchingItem{
	"@":       KEYWORD_CALL,
	"@define": KEYWORD_DEFINE,
	"(":       KEYWORD_BRACKET_OPEN,
	")":       KEYWORD_BRACKET_CLOSE,

	"`":  STRING_QUOTEMARK,
	"!t": BOOLEAN_TRUE,
	"!f": BOOLEAN_FALSE,

	" ":  WHITESPACE,
	"\n": NEWLINE,
}

var InvertedKeywordMap = make(InvertedKeywordMatchingItem)

type PossibleValueMatchingItem map[string]LexerTokenDataType

var PossibleValueMap PossibleValueMatchingItem = PossibleValueMatchingItem{
	"0": DATA_INT,
	"1": DATA_INT,
	"2": DATA_INT,
	"3": DATA_INT,
	"4": DATA_INT,
	"5": DATA_INT,
	"6": DATA_INT,
	"7": DATA_INT,
	"8": DATA_INT,
	"9": DATA_INT,

	".": DATA_REAL,
}

func init() {
	for key, val := range KeywordMap {
		InvertedKeywordMap[val] = key
	}
}
