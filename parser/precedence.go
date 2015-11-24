package parser

import (
	. "github.com/jonnyarnold/fn-go/tokeniser"
)

// The precedence of infix operations.
var InfixPrecedence = []string{".", "=", "eq", "and", "or", "*", "/", "+", "-"}

// Get the precedence of a token.
func precedenceOf(token Token) float32 {
	for precedence, tokenType := range InfixPrecedence {
		if token.Value == tokenType {
			return float32(precedence)
		}
	}

	return -2
}
