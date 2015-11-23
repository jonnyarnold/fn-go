package parser

import (
	. "github.com/jonnyarnold/fn-go/tokeniser"
)

// A tokenList is a wrapper around a token array
// and provides handy functions.
type tokenList []Token

// Returns true if there are any elements.
func (tokens tokenList) Any() bool {
	return len(tokens) != 0
}

// Returns the first element.
func (tokens tokenList) Next() Token {
	return tokens[0]
}

// Returns every element except the first.
func (tokens tokenList) Pop() tokenList {
	return tokens[1:]
}

func (tokens tokenList) Peek(index int) Token {
	return tokens[index]
}

// Returns a pointer to the first token _after_ the token
// of the specified type.
//
// For example, in the token list `(1 + 2) * 3`
// tokenList.AfterNext() would point to `*`
func (tokens tokenList) AfterNext(tokenType string) *Token {
	for idx, token := range tokens {
		if token.Type == tokenType {
			return &tokens[idx+1]
		}
	}

	return nil
}
