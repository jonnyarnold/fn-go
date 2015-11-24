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
	var lhs Expression

	switch tokens.Next().Type {

	case "identifier":
		// Peek forward to see if we have a function call
		if tokens.Peek(1).Type == "bracket_open" {
			// TODO: Error checking!
			lhs, tokens, _ = parseFunctionCall(tokens)
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
		if tokenAfterClosingBracket.Type == "block_open" {
			lhs, tokens, _ = parseFunctionDefinition(tokens)
		} else {
			lhs, tokens, _ = parseBrackets(tokens)
		}

	case "block_open":
		// TODO: Error checking!
		lhs, tokens, _ = parseBlock(tokens)

	case "when":
		// TODO: Error checking!
		lhs, tokens, _ = parseWhen(tokens)
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
	var rhs Expression

	for tokens.Any() {
		beforeParsePrecedence := precedenceOf(tokens.Next())

		if beforeParsePrecedence < precedence {
			break
		}

		operation := tokens.Next().Value
		tokens = tokens.Pop() // Eat infix_operator

		// TODO: Error checking!
		rhs, tokens, _ = parseValue(tokens)

		// If, after parsing, the current token has a higher precedence,
		// we need to use everything we have so far at the LHS of the higher expression.
		if beforeParsePrecedence < precedenceOf(tokens.Next()) {
			// TODO: Error checking!
			rhs, tokens, _ = parseInfixRhs(tokens, precedence+0.01, rhs)
		}

		lhs = FunctionCallExpression{
			Identifier: IdentifierExpression{Name: operation},
			Arguments:  []Expression{lhs, rhs},
		}
	}

	return lhs, tokens, nil
}
