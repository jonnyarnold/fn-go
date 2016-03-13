package parser

import (
	"errors"
	"fmt"
)

// Parse a value.
// Values are of the form
// `identifier | function_call | number | string | boolean |
//  function_definition | brackets | block | when`
func parseValue(tokens tokenList) (Expression, tokenList, error) {
	// We need to parse the left and right hand side for the value.
	var (
		lhs Expression
		err error
	)

	switch tokens.Next().Type {

	case "identifier":
		// Peek forward to see if we have a function call
		if tokens.Length() > 1 && tokens.Peek(1).Type == "bracket_open" {
			// TODO: Error checking!
			lhs, tokens, err = parseFunctionCall(tokens)
		} else {
			lhs = IdentifierExpression{Name: tokens.Next().Value}
			tokens = tokens.Pop()
		}

	// Basic literals
	case "number":
		lhs = NumberExpression{Value: tokens.Next().Value}
		tokens = tokens.Pop()
	case "string":
		lhs = StringExpression{Value: tokens.Next().Value}
		tokens = tokens.Pop()
	case "boolean":
		lhs = BooleanExpression{Value: tokens.Next().Value == "true"}
		tokens = tokens.Pop()

	case "bracket_open":
		// Find the token immediately after the closing bracket.
		tokenAfterClosingBracket := tokens.AfterNext("bracket_close")

		// TODO: Error checking!
		if tokenAfterClosingBracket != nil && tokenAfterClosingBracket.Type == "block_open" {
			lhs, tokens, err = parseFunctionDefinition(tokens)
		} else {
			lhs, tokens, err = parseBrackets(tokens)
		}

	case "block_open":
		lhs, tokens, err = parseBlock(tokens)

	case "when":
		lhs, tokens, err = parseWhen(tokens)
	}

	if err != nil {
		return nil, tokens, err
	}

	// Check we parsed the LHS
	if lhs == nil {
		return nil, tokens, errors.New(
			fmt.Sprintf("Unexpected token type %s when parsing value", tokens.Next().Type),
		)
	}

	return parseInfixRhs(tokens, -1, lhs)
}

// Parse the Right-Hand side of an expression, respecting precedence.
func parseInfixRhs(tokens tokenList, precedence float32, lhs Expression) (Expression, tokenList, error) {
	var (
		rhs Expression
		err error
	)

	for tokens.Any() {
		beforeParsePrecedence := precedenceOf(tokens.Next())

		if beforeParsePrecedence < precedence {
			break
		}

		operation := tokens.Next().Value
		tokens = tokens.Pop() // Eat infix_operator

		rhs, tokens, err = parseValue(tokens)
		if err != nil {
			return rhs, tokens, err
		}

		// If, after parsing, the current token has a higher precedence,
		// we need to use everything we have so far at the LHS of the higher expression.
		if tokens.Any() && beforeParsePrecedence < precedenceOf(tokens.Next()) {
			rhs, tokens, err = parseInfixRhs(tokens, precedence, rhs)
			if err != nil {
				return rhs, tokens, err
			}
		}

		lhs = FunctionCallExpression{
			Identifier: IdentifierExpression{Name: operation},
			Arguments: []Expression{
				lhs,
				rhs,
			},
		}
	}

	return lhs, tokens, nil
}
