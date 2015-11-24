package parser

import (
	"errors"
	"fmt"
)

// Parse brackets of the form `(primary)`
func parseBrackets(tokens tokenList) (Expression, tokenList, error) {
	if tokens.Next().Type != "bracket_open" {
		return nil, tokens, errors.New(
			fmt.Sprintf("Expected bracket_open, found %s in bracketed expression", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat bracket_open

	expr, tokens, err := parsePrimary(tokens)
	if err != nil {
		return nil, tokens, err
	}

	if tokens.Next().Type != "bracket_close" {
		return nil, tokens, errors.New(
			fmt.Sprintf("Expected bracket_close, found %s in bracketed expression", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat bracket_close

	return expr, tokens, nil
}
