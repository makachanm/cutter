package parser

import "cutter/lexer"

type ParserStack struct {
	items []lexer.LexerToken
}

func NewParserStack() *ParserStack {
	s := new(ParserStack)
	s.items = make([]lexer.LexerToken, 0)

	return s
}

func (s *ParserStack) Push(item lexer.LexerToken) {
	s.items = append(s.items, item)
}

func (s *ParserStack) Pop() (*lexer.LexerToken, bool) {
	if s.IsEmpty() {
		return nil, false
	}
	lastIndex := len(s.items) - 1
	item := s.items[lastIndex]
	s.items = s.items[:lastIndex]
	return &item, true
}

func (s *ParserStack) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *ParserStack) Clear() {
	for !s.IsEmpty() {
		s.Pop()
	}
}

type ParserQueue struct {
	contents []lexer.LexerToken
	size     uint64
	pointer  uint64
}

func NewParserQueue(c []lexer.LexerToken, s uint64) *ParserQueue {
	q := new(ParserQueue)
	q.contents = c
	q.size = s
	q.pointer = 0

	return q
}

func (q *ParserQueue) Pop() lexer.LexerToken {
	if q.pointer+1 >= q.size {
		return lexer.LexerToken{}
	}

	q.pointer++
	return q.contents[q.pointer]
}

func (q *ParserQueue) Pushback() lexer.LexerToken {
	if q.pointer <= 0 {
		return lexer.LexerToken{}
	}

	q.pointer--
	return q.contents[q.pointer]
}

func (q *ParserQueue) Seek() lexer.LexerToken {
	return q.contents[q.pointer]
}

func (q *ParserQueue) IsEmpty() bool {
	return (q.pointer+1 >= q.size)
}
