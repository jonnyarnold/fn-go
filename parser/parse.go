package parser

import (
	"errors"
	"fmt"
)

// Converts a list of Tokens into a list of Expressions.
func Parse(tokens tokenList) ([]Expression, error) {
	expressions := []Expression{}

	for tokens.Any() {
		newExpression, remainingTokens, err := parsePrimary(tokens)

		if newExpression != nil {
			expressions = append(expressions, newExpression)
		} else if err != nil {
			// Die on the first error.
			return expressions, err
		} else if len(remainingTokens) == len(tokens) {
			return expressions, errors.New("Parser stalled!")
		}

		tokens = remainingTokens
	}

	return expressions, nil
}

// Parse a top-level expression.
// Primaries are of the form `value | import`
func parsePrimary(tokens tokenList) (Expression, tokenList, error) {
	switch tokens.Next().Type {
	case "identifier", "number", "string", "boolean", "bracket_open", "when", "block_open":
		return parseValue(tokens)
	}

	// TODO: import

	return nil, tokens, errors.New(
		fmt.Sprintf("Unexpected token type %s at start of expression", tokens.Next().Type),
	)
}
