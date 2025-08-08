package lexer

type TokenQueue struct {
	contents []Token
	size     uint64
	pointer  uint64
}

func NewTokenQueue(c []Token, s uint64) *TokenQueue {
	q := new(TokenQueue)
	q.contents = c
	q.size = s
	q.pointer = 0

	return q
}

func (q *TokenQueue) Pop() Token {
	if q.pointer+1 >= q.size {
		return Token{}
	}

	q.pointer++
	return q.contents[q.pointer]
}

func (q *TokenQueue) Pushback() Token {
	if q.pointer <= 0 {
		return Token{}
	}

	q.pointer--
	return q.contents[q.pointer]
}

func (q *TokenQueue) Seek() Token {
	return q.contents[q.pointer]
}

func (q *TokenQueue) IsEmpty() bool {
	return (q.pointer+1 >= q.size)
}
