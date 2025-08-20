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
	d := fmt.Sprintf("invalid token: %%#v, expected was: %%d", err, expected)
	panic("unexpected syntax error - " + d)
}

func (p *Parser) makeDataError(expected lexer.LexerTokenDataType) {
	d := fmt.Sprintf("invalid data: %%d, expected was: %%d", p.targets.Seek().Data.Type, expected)
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

	p.validCheckPop(lexer.KEYWORD_BRACKET_OPEN)

	for {
		// Peek at the next token to see if it's the end
		next, ok := p.targets.Pop()
		if !ok {
			panic("unexpected end of token stream, expected ')'")
		}
		p.targets.Pushback()

		if next.Type == lexer.KEYWORD_BRACKET_CLOSE {
			p.targets.Pop() // Consume the closing bracket
			break
		}

		object, ok := p.targets.Pop()
		if !ok {
			panic("unexpected end of token stream, expected ')'")
		}

		if object.Type == lexer.WHITESPACE || object.Type == lexer.NEWLINE {
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
		default:
			if object.Type != lexer.KEYWORD_BRACKET_CLOSE {
				panic(fmt.Sprintf("unexpected token in argument list: %#v", object))
			}
		}
	}

	return call
}

func (p *Parser) doDefineParse() FunctionObject {
	fun := FunctionObject{
		Parameters: make([]string, 0),
	}

	p.validCheckPop(lexer.KEYWORD_BRACKET_OPEN)
	object := p.validCheckPop(lexer.VALUE)
	if object.Data.Type != lexer.DATA_OBJNAME {
		p.makeDataError(lexer.DATA_OBJNAME)
	}
	fun.Name = object.Data.ObjNameData

	tempArgs := make([]CallObject, 0)

	for {
		// Peek at the next token to see if it's the end
		next, ok := p.targets.Pop()
		if !ok {
			panic("unexpected end of token stream, expected ')'")
		}
		p.targets.Pushback()

		if next.Type == lexer.KEYWORD_BRACKET_CLOSE {
			p.targets.Pop() // Consume the closing bracket
			break
		}

		// Now we know it's not the end, so we parse an argument/body/staticdata
		object, _ := p.targets.Pop()
		if object.Type == lexer.WHITESPACE || object.Type == lexer.NEWLINE {
			continue
		}

		switch object.Data.Type {
		case lexer.DATA_OBJNAME:
			next, _ := p.targets.Pop()
			p.targets.Pushback()

			if next.Type != lexer.KEYWORD_BRACKET_OPEN {
				tempArgs = append(tempArgs, CallObject{Name: object.Data.ObjNameData})
			} else {
				p.targets.Pushback()
				tempArgs = append(tempArgs, p.doCallParse())
			}

		case lexer.DATA_INT:
			fun.StaticData = makeIntValueObj(object.Data.IntData)
		case lexer.DATA_REAL:
			fun.StaticData = makeRealValueObj(object.Data.RealData)
		case lexer.DATA_STR:
			fun.StaticData = makeStrValueObj(object.Data.StrData)
		case lexer.DATA_BOOL:
			fun.StaticData = makeBoolValueObj(object.Data.BoolData)
		default:
			panic(fmt.Sprintf("unexpected token in function definition: %%#v", object))
		}
	}

	// Now, interpret tempArgs into Parameters and Body
	if len(tempArgs) > 0 {
		// The last item is the body
		fun.Body = tempArgs[len(tempArgs)-1]

		// Everything before the last item is a parameter name
		for i := 0; i < len(tempArgs)-1; i++ {
			// Check that parameters are just names and not calls
			if len(tempArgs[i].Arguments) > 0 {
				panic(fmt.Sprintf("callable object is not allowed as a parameter name: %%#v", tempArgs[i]))
			}
			fun.Parameters = append(fun.Parameters, tempArgs[i].Name)
		}
	}

	return fun
}
