package parser

import "cutter/lexer"

type ParserQueue struct {
	contents []lexer.LexerToken
	size     int64
	pointer  int64
}

func NewParserQueue(c []lexer.LexerToken, s int64) *ParserQueue {
	q := new(ParserQueue)
	q.contents = c
	q.size = s
	q.pointer = -1

	return q
}

func (q *ParserQueue) Pop() (lexer.LexerToken, bool) {
	if q.pointer+1 >= q.size {
		return lexer.LexerToken{}, false
	}

	q.pointer += 1
	return q.contents[q.pointer], true
}

func (q *ParserQueue) Pushback() (lexer.LexerToken, bool) {
	if q.pointer <= 0 {
		return lexer.LexerToken{}, false
	}

	q.pointer -= 1
	return q.contents[q.pointer], true
}

func (q *ParserQueue) Seek() lexer.LexerToken {
	return q.contents[q.pointer]
}

func (q *ParserQueue) IsEmpty() bool {
	return (q.pointer+1 >= q.size)
}
