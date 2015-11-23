package parser

import (
	"fmt"
	. "github.com/jonnyarnold/fn-go/tokeniser"
)

// A tokenList is a wrapper around a token array
// and provides handy functions.
type tokenList []Token

// Returns the first element.
func (tokens tokenList) Next() Token {
	return tokens[0]
}

// Returns every element except the first.
func (tokens tokenList) Pop() tokenList {
	fmt.Printf("Popped %s\n", tokens[0])
	return tokens[1:]
}
