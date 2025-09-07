package lexer

import (
	"sort"
)

type Tokenizer struct {
	pointer uint64
	maxsize uint64

	targets   []rune
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

	tk.targets = []rune(input)
	tk.maxsize = uint64(len(tk.targets))

	buffer := make([]rune, 0)

	for {
		if tk.pointer >= tk.maxsize {
			if len(buffer) > 0 {
				tk.tokenized = append(tk.tokenized, NewDataToken(NORM_STRINGS, string(buffer)))
			}

			tk.tokenized = append(tk.tokenized, NewToken(TERMINATOR))
			break
		}

		token_type, token_size := tk.matchToken()

		if token_type != NORM_STRINGS {
			if len(buffer) > 0 {
				tk.tokenized = append(tk.tokenized, NewDataToken(NORM_STRINGS, string(buffer)))
			}

			buffer = buffer[:0]

			tk.tokenized = append(tk.tokenized, NewToken(token_type))
			tk.pointer += uint64(token_size)
		} else {
			buffer = append(buffer, tk.targets[tk.pointer])
			tk.pointer++
		}
	}

	return tk.tokenized
}

type tokenHead struct {
	token TokenType
	len   int
}

func (tk *Tokenizer) matchToken() (TokenType, int) {
	keyword_tokens := make([]tokenHead, 0)

	for key, value := range KeywordMap {
		keywordRunes := []rune(key)
		if int(tk.pointer)+len(keywordRunes) > len(tk.targets) {
			continue
		}

		match := true
		for i, r := range keywordRunes {
			if tk.targets[tk.pointer+uint64(i)] != r {
				match = false
				break
			}
		}

		if match {
			keyword_tokens = append(keyword_tokens, tokenHead{token: value, len: len(keywordRunes)})
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
