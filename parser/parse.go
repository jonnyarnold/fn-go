package parser

import (
	"errors"
	"fmt"
	. "github.com/jonnyarnold/fn-go/tokeniser"
)

// Converts a list of Tokens into a list of Expressions.
func Parse(tokens []Token) ([]Expression, error) {
	expressions := []Expression{}

	for tokens != nil {
		newExpression, remainingTokens, err := parsePrimary(tokens)

		if newExpression != nil {
			expressions = append(expressions, newExpression)
		} else if err != nil {
			// Die on the first error.
			return expressions, err
		}

		tokens = remainingTokens
	}

	return expressions, nil
}

func parsePrimary(tokens []Token) (Expression, []Token, error) {
	switch tokens[0].Type {
	case "comment", "space":
		return nil, tokens[1:], nil
	case "identifier", "number", "string", "boolean":
		return parseValue(tokens)
	}

	return nil, tokens, errors.New(
		fmt.Sprintf("Unexpected token type %s at start of expression", tokens[0].Type),
	)
}

func parseValue(tokens []Token) (Expression, []Token, error) {
	switch tokens[0].Type {

	case "identifier":
		return IdentifierExpression{Name: tokens[0].Value}, tokens[1:], nil

	// Basic literals
	case "number":
		return NumberExpression{Value: tokens[0].Value}, tokens[1:], nil
	case "string":
		return StringExpression{Value: tokens[0].Value}, tokens[1:], nil
	case "boolean":
		return BooleanExpression{Value: tokens[0].Value == "true"}, tokens[1:], nil

	// More complex structures
	case "block_open":
		return parseBlock(tokens)
	}

	return nil, tokens, errors.New(
		fmt.Sprintf("Unexpected token type %s when parsing value", tokens[0].Type),
	)
}

func parseBlock(tokens []Token) (Expression, []Token, error) {
	if tokens[0].Type != "block_open" {
		return nil, tokens, errors.New(
			fmt.Sprintf("Unexpected token type %s at start of block", tokens[0].Type),
		)
	}

	tokens = tokens[1:] // Eat block_open

	// Eat the body
	body := []Expression{}
	for tokens != nil && tokens[0].Type != "block_close" {
		newExpression, remainingTokens, err := parsePrimary(tokens)

		if newExpression != nil {
			body = append(body, newExpression)
		} else if err != nil {
			// Die on the first error.
			return body, tokens, err
		}

		tokens = remainingTokens
	}

	// Check we're at a closing block.
	if tokens == nil {
		return nil, tokens, errors.New(
			fmt.Sprintf("End of file reached before closing block", tokens[0].Type),
		)
	}

	return BlockExpression{Body: body}, tokens, nil
}
