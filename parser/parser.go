package parser

import "cutter/lexer"

type Parser struct {
	tokenqueue  *ParserQueue
	bufferstack *ParserStack

	head HeadNode
}

func NewParser() *Parser {
	p := new(Parser)
	p.bufferstack = NewParserStack()
	p.head = HeadNode{}
	p.head.Bodys = make([]BodyObject, 0)

	return p
}

func (p *Parser) DoParse(tokens []lexer.LexerToken) {
	p.tokenqueue = NewParserQueue(tokens, uint64(len(tokens)))

	for !p.tokenqueue.IsEmpty() {
		c_token := p.tokenqueue.Pop()
		switch c_token.Type {

		}

	}
}
