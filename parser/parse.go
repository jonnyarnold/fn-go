package parser

import (
	"errors"
	"fmt"
)

// Converts a list of Tokens into a list of Expressions.
func Parse(tokens tokenList) ([]Expression, error) {
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

func parsePrimary(tokens tokenList) (Expression, tokenList, error) {
	switch tokens.Next().Type {
	case "identifier", "number", "string", "boolean":
		return parseValue(tokens)
	}

	return nil, tokens, errors.New(
		fmt.Sprintf("Unexpected token type %s at start of expression", tokens.Next().Type),
	)
}

func parseValue(tokens tokenList) (Expression, tokenList, error) {
	// We need to parse the left and right hand side for the value.
	var lhs Expression

	switch tokens.Next().Type {

	case "identifier":
		lhs, tokens = IdentifierExpression{Name: tokens.Next().Value}, tokens.Pop()

	// Basic literals
	case "number":
		lhs, tokens = NumberExpression{Value: tokens.Next().Value}, tokens.Pop()
	case "string":
		lhs, tokens = StringExpression{Value: tokens.Next().Value}, tokens.Pop()
	case "boolean":
		lhs, tokens = BooleanExpression{Value: tokens.Next().Value == "true"}, tokens.Pop()

	// More complex structures
	case "block_open":
		lhs, tokens, _ = parseBlock(tokens)
	}

	// Check we parsed the LHS
	if lhs == nil {
		return nil, tokens, errors.New(
			fmt.Sprintf("Unexpected token type %s when parsing value", tokens.Next().Type),
		)
	}

	return parseInfixRhs(tokens, -1, lhs)
}

// parse the Right-Hand side of an expression.
func parseInfixRhs(tokens tokenList, precedence float32, lhs Expression) (Expression, tokenList, error) {
	for true {
		beforeParsePrecedence := precedenceOf(tokens.Next())

		if beforeParsePrecedence < precedence {
			break
		}

		operation := tokens.Next().Value
		tokens = tokens.Pop()

		rhs, tokens, _ := parseValue(tokens)

		// If, after parsing, the current token has a higher precedence,
		// we need to use everything we have so far at the LHS of the higher expression.
		if beforeParsePrecedence < precedenceOf(tokens.Next()) {
			rhs, tokens, _ = parseInfixRhs(tokens, precedence+0.01, rhs)
		}

		lhs = FunctionCallExpression{
			Identifier: IdentifierExpression{Name: operation},
			Arguments:  []Expression{lhs, rhs},
		}
	}

	return lhs, tokens, nil
}

func parseBlock(tokens tokenList) (Expression, tokenList, error) {
	if tokens.Next().Type != "block_open" {
		return nil, tokens, errors.New(
			fmt.Sprintf("Unexpected token type %s at start of block", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat block_open

	// Eat the body
	body := []Expression{}
	for tokens != nil && tokens.Next().Type != "block_close" {
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
			fmt.Sprintf("End of file reached before closing block", tokens.Next().Type),
		)
	}

	return BlockExpression{Body: body}, tokens, nil
}
