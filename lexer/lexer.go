package lexer

import (
	"strconv"
	"strings"
)

type Lexer struct {
	state LexerStatus
	queue *TokenQueue

	buffer  []string
	results []LexerToken

	inDefine           bool
	defineBracketLevel int
	inInclude          bool
	ignoreNextNewline  bool
}

func NewLexer() *Lexer {
	lex := new(Lexer)
	lex.state = STATE_NORMSTRINGS
	lex.results = make([]LexerToken, 0)
	lex.inDefine = false
	lex.defineBracketLevel = 0
	lex.inInclude = false
	lex.ignoreNextNewline = false

	return lex
}

func (l *Lexer) DoLex(input string) []LexerToken {
	raw_tokens := NewTokenizer().doTokenize(input, uint64(len(input)))
	l.queue = NewTokenQueue(raw_tokens, int64(len(raw_tokens)))

	l.buffer = make([]string, 0)

	for !l.queue.IsEmpty() {
		symbol := l.queue.Pop()

		if l.ignoreNextNewline && symbol.GetType() != NEWLINE {
			l.ignoreNextNewline = false
		}

		switch symbol.GetType() {
		case KEYWORD_CALL:
			if l.state == STATE_STRINGVALUE {
				l.buffer = append(l.buffer, InvertedKeywordMap[symbol.token_type])
				continue
			}
			l.flushBuffer()
			l.results = append(l.results, NewLexerToken(symbol.token_type, NewData()))
			l.state = STATE_OBJNAME

		case KEYWORD_DEFINE:
			if l.state == STATE_STRINGVALUE {
				l.buffer = append(l.buffer, InvertedKeywordMap[symbol.token_type])
				continue
			}
			l.flushBuffer()
			l.results = append(l.results, NewLexerToken(symbol.token_type, NewData()))
			l.inDefine = true
		
		case KEYWORD_INCLUDE:
			if l.state == STATE_STRINGVALUE {
				l.buffer = append(l.buffer, InvertedKeywordMap[symbol.token_type])
				continue
			}
			l.flushBuffer()
			l.results = append(l.results, NewLexerToken(symbol.token_type, NewData()))
			l.inInclude = true

		case KEYWORD_BRACKET_OPEN:
			if l.state == STATE_STRINGVALUE {
				l.buffer = append(l.buffer, InvertedKeywordMap[symbol.token_type])
				continue
			}
			if l.inDefine {
				l.defineBracketLevel++
			}
			l.results = append(l.results, NewLexerToken(symbol.token_type, NewData()))
			l.state = STATE_OBJNAME

		case KEYWORD_BRACKET_CLOSE:
			if l.state == STATE_STRINGVALUE {
				l.buffer = append(l.buffer, InvertedKeywordMap[symbol.token_type])
				continue
			}
			l.results = append(l.results, NewLexerToken(symbol.token_type, NewData()))
			next := l.queue.Pop()
			l.queue.Pushback()

			if next.GetType() == WHITESPACE {
				l.state = STATE_OBJNAME
			} else {
				l.state = STATE_NORMSTRINGS
			}
			if l.inDefine {
				l.defineBracketLevel--
				if l.defineBracketLevel == 0 {
					l.inDefine = false // End of define block
					l.ignoreNextNewline = true
				}
			}
			if l.inInclude {
				l.inInclude = false
				l.ignoreNextNewline = true
			}

		case STRING_QUOTEMARK:
			if l.state != STATE_STRINGVALUE {
				l.flushBuffer()
				l.state = STATE_STRINGVALUE
			} else {
				l.flushBuffer()
				l.state = STATE_OBJNAME
			}

		case BOOLEAN_TRUE, BOOLEAN_FALSE:
			if l.state == STATE_STRINGVALUE {
				l.buffer = append(l.buffer, InvertedKeywordMap[symbol.token_type])
				continue
			}
			var data bool
			switch symbol.token_type {
			case BOOLEAN_TRUE:
				data = true

			case BOOLEAN_FALSE:
				data = false
			}

			l.results = append(l.results, NewLexerToken(VALUE, NewBoolData(data)))

		case WHITESPACE, NEWLINE:
			if l.ignoreNextNewline && symbol.GetType() == NEWLINE {
				l.ignoreNextNewline = false
				continue
			}
			if l.state == STATE_STRINGVALUE || l.state == STATE_NORMSTRINGS {
				l.buffer = append(l.buffer, InvertedKeywordMap[symbol.token_type])
				continue
			}

			//l.results = append(l.results, NewLexerToken(symbol.token_type, NewData()))

			l.flushBuffer()

		case NORM_STRINGS:
			if l.state == STATE_STRINGVALUE || l.state == STATE_NORMSTRINGS {
				l.buffer = append(l.buffer, symbol.GetData())
				continue
			}

			//if l.state == STATE_OBJNAME {
			//	l.results = append(l.results, NewLexerToken(VALUE, NewObjNameData(symbol.GetData())))
			//	continue
			//}

			l.results = append(l.results, NewLexerToken(VALUE, l.getValues(symbol.GetData())))

		case TERMINATOR:
			l.flushBuffer()
			l.results = append(l.results, NewLexerToken(TERMINATOR, NewData()))
		}
	}

	return l.results
}

func (l *Lexer) getValues(data string) LexerTokenData {
	var dtype LexerTokenDataType = DATA_INT
	for _, chars := range strings.Split(data, "") {
		if dtype != PossibleValueMap[chars] {
			dtype = PossibleValueMap[chars]
			break
		}
	}

	switch dtype {
	case DATA_INT:
		d, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			panic(err)
		}
		return NewIntData(d)

	case DATA_REAL:
		d, err := strconv.ParseFloat(data, 64)
		if err != nil {
			panic(err)
		}
		return NewRealData(d)

	default:
		return NewObjNameData(data)
	}
}

func (l *Lexer) flushBuffer() {
	buffer_d := strings.Join(l.buffer, "")
	var data LexerTokenData

	if buffer_d == "" {
		return
	}

	var tokentype TokenType

	switch l.state {
	case STATE_STRINGVALUE:
		l.buffer = l.buffer[:0]
		data = NewStrData(buffer_d)
		tokentype = VALUE

	case STATE_NORMSTRINGS:
		l.buffer = l.buffer[:0]
		data = NewNormData(buffer_d)
		tokentype = NORM_STRINGS
	}

	l.results = append(l.results, NewLexerToken(tokentype, data))
}
