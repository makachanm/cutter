package parser

import (
	"cutter/lexer"
	"fmt"
)

type Parser struct {
	targets *ParserQueue
}

func NewParser() *Parser {
	p := new(Parser)

	return p
}

func (p *Parser) makeTokenError(expected lexer.TokenType, err lexer.LexerToken) {
	d := fmt.Sprintf("invalid token: %#v, expected was: %d", err, expected)
	panic("unexpected syntax error - " + d)
}

func (p *Parser) makeDataError(expected lexer.LexerTokenDataType) {
	d := fmt.Sprintf("invalid data: %d, expected was: %d", p.targets.Seek().Data.Type, expected)
	panic("unexpected syntax error - " + d)
}

func (p *Parser) validCheckPop(target_token lexer.TokenType) lexer.LexerToken {
	val, _ := p.targets.Pop()
	if val.Type != target_token {
		p.makeTokenError(target_token, val)
		return lexer.LexerToken{}
	} else {
		return val
	}
}

func (p *Parser) DoParse(tokens []lexer.LexerToken) HeadNode {
	head := HeadNode{}
	head.Bodys = make([]BodyObject, 0)
	p.targets = NewParserQueue(tokens, int64(len(tokens)))

	for !p.targets.IsEmpty() {
		c_token, _ := p.targets.Pop()

		switch c_token.Type {
		case lexer.KEYWORD_CALL:
			call := p.doCallParse()
			head.Bodys = append(head.Bodys, BodyObject{
				Type: FUNCTION_CALL,
				Call: call,
			})

		case lexer.KEYWORD_DEFINE:
			fun := p.doDefineParse()
			head.Bodys = append(head.Bodys, BodyObject{
				Type: FUCNTION_DEFINITION,
				Func: fun,
			})

		case lexer.NORM_STRINGS:
			head.Bodys = append(head.Bodys, BodyObject{
				Type: NORM_STRINGS,
				Norm: NormStringObject{Data: c_token.Data.NormData},
			})

		default:
			continue
		}

	}

	return head
}

func (p *Parser) doCallParse() CallObject {
	var call CallObject = CallObject{}
	call.Arguments = make([]Argument, 0)

	object := p.validCheckPop(lexer.VALUE)
	if object.Data.Type != lexer.DATA_OBJNAME {
		p.makeDataError(lexer.DATA_OBJNAME)
	}
	call.Name = object.Data.ObjNameData

	object = p.validCheckPop(lexer.KEYWORD_BRACKET_OPEN)
	for object.Type != lexer.KEYWORD_BRACKET_CLOSE {
		object, _ = p.targets.Pop()
		if object.Type == lexer.WHITESPACE {
			continue
		}

		switch object.Data.Type {
		case lexer.DATA_INT:
			call.Arguments = append(call.Arguments, Argument{Type: ARG_LITERAL, Literal: makeIntValueObj(object.Data.IntData)})
		case lexer.DATA_REAL:
			call.Arguments = append(call.Arguments, Argument{Type: ARG_LITERAL, Literal: makeRealValueObj(object.Data.RealData)})
		case lexer.DATA_STR:
			call.Arguments = append(call.Arguments, Argument{Type: ARG_LITERAL, Literal: makeStrValueObj(object.Data.StrData)})
		case lexer.DATA_BOOL:
			call.Arguments = append(call.Arguments, Argument{Type: ARG_LITERAL, Literal: makeBoolValueObj(object.Data.BoolData)})

		case lexer.DATA_OBJNAME:
			next, _ := p.targets.Pop()
			p.targets.Pushback()

			if next.Type != lexer.KEYWORD_BRACKET_OPEN {
				call.Arguments = append(call.Arguments, Argument{Type: ARG_VARIABLE, VarName: object.Data.ObjNameData})
			} else {
				p.targets.Pushback()
				subcall := p.doCallParse()

				call.Arguments = append(call.Arguments, Argument{Type: ARG_CALLABLE, Callable: subcall})
			}

		}
	}

	return call
}

func (p *Parser) doDefineParse() FunctionObject {
	var fun FunctionObject = FunctionObject{}

	object := p.validCheckPop(lexer.KEYWORD_BRACKET_OPEN)
	object = p.validCheckPop(lexer.VALUE)
	if object.Data.Type != lexer.DATA_OBJNAME {
		p.makeDataError(lexer.DATA_OBJNAME)
	}

	fun.Name = object.Data.ObjNameData

	for object.Type != lexer.KEYWORD_BRACKET_CLOSE {
		object, _ = p.targets.Pop()

		switch object.Data.Type {
		case lexer.DATA_OBJNAME:
			next, _ := p.targets.Pop()
			p.targets.Pushback()

			if next.Type != lexer.KEYWORD_BRACKET_OPEN {
				fun.Args = append(fun.Args, CallObject{Name: object.Data.ObjNameData})
			} else {
				p.targets.Pushback()
				fun.Args = append(fun.Args, p.doCallParse())
			}

		case lexer.DATA_INT:
			fun.StaticData = makeIntValueObj(object.Data.IntData)
		case lexer.DATA_REAL:
			fun.StaticData = makeRealValueObj(object.Data.RealData)
		case lexer.DATA_STR:
			fun.StaticData = makeStrValueObj(object.Data.StrData)
		case lexer.DATA_BOOL:
			fun.StaticData = makeBoolValueObj(object.Data.BoolData)
		}
	}

	return fun
}