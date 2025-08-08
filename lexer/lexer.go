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
}

func NewLexer() *Lexer {
	lex := new(Lexer)
	lex.state = STATE_NORMSTRINGS
	lex.results = make([]LexerToken, 0)

	return lex
}

func (l *Lexer) DoLex(input string) []LexerToken {
	raw_tokens := NewTokenizer().doTokenize(input, uint64(len(input)))
	l.queue = NewTokenQueue(raw_tokens, uint64(len(raw_tokens)))

	l.buffer = make([]string, 0)

	for !l.queue.IsEmpty() {
		symbol := l.queue.Pop()

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

		case KEYWORD_BRACKET_OPEN:
			if l.state == STATE_STRINGVALUE {
				l.buffer = append(l.buffer, InvertedKeywordMap[symbol.token_type])
				continue
			}
			l.results = append(l.results, NewLexerToken(symbol.token_type, NewData()))
			l.state = STATE_OBJNAME

		case KEYWORD_BRACKET_CLOSE:
			if l.state == STATE_STRINGVALUE {
				l.buffer = append(l.buffer, InvertedKeywordMap[symbol.token_type])
				continue
			}
			l.results = append(l.results, NewLexerToken(symbol.token_type, NewData()))
			l.state = STATE_NORMSTRINGS

		case STRING_QUOTEMARK:
			if l.state != STATE_STRINGVALUE {
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
			if l.state == STATE_STRINGVALUE || l.state == STATE_NORMSTRINGS {
				l.buffer = append(l.buffer, InvertedKeywordMap[symbol.token_type])
				continue
			}
			l.flushBuffer()

		case NORM_STRINGS:
			if l.state == STATE_STRINGVALUE || l.state == STATE_NORMSTRINGS {
				l.buffer = append(l.buffer, symbol.GetData())
				continue
			}

			if l.state == STATE_OBJNAME {
				l.results = append(l.results, NewLexerToken(VALUE, NewObjNameData(symbol.GetData())))
				continue
			}

			l.results = append(l.results, NewLexerToken(VALUE, l.getValues(symbol.GetData())))
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
	}

	return NewData()
}

func (l *Lexer) flushBuffer() {
	buffer_d := strings.Join(l.buffer, "")
	var data LexerTokenData

	if buffer_d == "" {
		return
	}

	switch l.state {
	case STATE_STRINGVALUE:
		l.buffer = l.buffer[:0]
		data = NewStrData(buffer_d)

	case STATE_NORMSTRINGS:
		data = NewNormData(buffer_d)
	}

	l.results = append(l.results, NewLexerToken(VALUE, data))
}
