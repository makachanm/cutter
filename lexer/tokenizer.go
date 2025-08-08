package lexer

import (
	"sort"
	"strings"
)

type Tokenizer struct {
	pointer uint64
	maxsize uint64

	targets   string
	tokenized []Token
}

func NewTokenizer() *Tokenizer {
	tokenizer := new(Tokenizer)

	tokenizer.pointer = 0
	tokenizer.maxsize = 0

	return tokenizer
}

func (tk *Tokenizer) doTokenize(input string, size uint64) []Token {
	tk.tokenized = make([]Token, 0)

	tk.maxsize = size
	tk.targets = input

	token_type, token_size := tk.matchToken()
	buffer := make([]string, 0)

	for {
		if tk.pointer >= tk.maxsize {
			if d := strings.Join(buffer, ""); d != "" {
				tk.tokenized = append(tk.tokenized, NewDataToken(NORM_STRINGS, strings.Join(buffer, "")))
			}

			break
		}

		token_type, token_size = tk.matchToken()

		if token_type != NORM_STRINGS {
			if d := strings.Join(buffer, ""); d != "" {
				tk.tokenized = append(tk.tokenized, NewDataToken(NORM_STRINGS, d))
			}

			buffer = buffer[:0]

			tk.tokenized = append(tk.tokenized, NewToken(token_type))
			tk.pointer += uint64(token_size) - 1
		} else {
			buffer = append(buffer, string(tk.targets[tk.pointer]))
		}

		tk.pointer++
	}

	return tk.tokenized
}

func (tk *Tokenizer) currCharVariableSize(size int) string {
	return string(tk.targets[tk.pointer : tk.pointer+uint64(size)])
}

type tokenHead struct {
	token TokenType
	len   int
}

func (tk *Tokenizer) matchToken() (TokenType, int) {
	keyword_tokens := make([]tokenHead, 0)

	for key := range KeywordMap {
		if tk.pointer+uint64(len(key))-1 >= tk.maxsize {
			continue
		}

		value, exist := KeywordMap[tk.currCharVariableSize(len(key))]
		if exist {
			keyword_tokens = append(keyword_tokens, tokenHead{token: value, len: len(key)})
		}
	}

	if len(keyword_tokens) <= 0 {
		return NORM_STRINGS, 1
	} else {
		sort.Slice(keyword_tokens, func(i, j int) bool {
			return keyword_tokens[i].len > keyword_tokens[j].len
		})

		return keyword_tokens[0].token, keyword_tokens[0].len
	}

}
