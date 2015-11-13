package parser

import (
	. "github.com/jonnyarnold/fn-go/tokeniser"
)

// Converts a list of Tokens into a list of Expressions.
func Parse(tokens []Token) []Expression {
	expressions := []Expression{}

	for tokens != nil {
		newExpressions, remainingTokens := parsePrimary(tokens)

		expressions = append(newExpressions, expressions)
		tokens = remainingTokens
	}

	return expressions
}

func parsePrimary(tokens []Token) ([]Expression, []Token) {
	switch tokens[0].Type {
	case "comment", "space":
		return nil, tokens[1:]
	}

	panic("Could not parse!")
}
