package lexer

type KeywordMatchingItem map[string]TokenType
type InvertedKeywordMatchingItem map[TokenType]string

var KeywordMap KeywordMatchingItem = KeywordMatchingItem{
	"@":       KEYWORD_CALL,
	"@define": KEYWORD_DEFINE,
	"@include": KEYWORD_INCLUDE,
	"(":       KEYWORD_BRACKET_OPEN,
	")":       KEYWORD_BRACKET_CLOSE,

	"`":  STRING_QUOTEMARK,
	"!t": BOOLEAN_TRUE,
	"!f": BOOLEAN_FALSE,

	" ":  WHITESPACE,
	"\n": NEWLINE,
}

var InvertedKeywordMap = make(InvertedKeywordMatchingItem)

func init() {
	for key, val := range KeywordMap {
		InvertedKeywordMap[val] = key
	}
}
