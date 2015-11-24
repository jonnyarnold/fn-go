package parser

import (
	"errors"
	"fmt"
)

// Parse a function call.
// Function calls are of the form `identifier ( params )`
func parseFunctionCall(tokens tokenList) (FunctionCallExpression, tokenList, error) {
	if tokens.Next().Type != "identifier" {
		return FunctionCallExpression{}, tokens, errors.New(
			fmt.Sprintf("Unexpected token type %s in function call", tokens.Next().Type),
		)
	}

	id := IdentifierExpression{Name: tokens.Next().Value}
	tokens = tokens.Pop() // Eat identifier

	args, tokens, err := parseParams(tokens)
	if err != nil {
		return FunctionCallExpression{}, tokens, err
	}

	return FunctionCallExpression{
		Identifier: id,
		Arguments:  args,
	}, tokens, nil
}

// Parse a parameter list roughly of the form `( [value ,]* )`
func parseParams(tokens tokenList) ([]Expression, tokenList, error) {
	if tokens.Next().Type != "bracket_open" {
		return nil, tokens, errors.New(
			fmt.Sprintf("Expected bracket_open, found %s in parameter list", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat bracket_open

	var (
		arg Expression
		err error
	)

	args := []Expression{}
	for tokens.Next().Type != "bracket_close" {
		arg, tokens, err = parseValue(tokens)
		if err != nil {
			return args, tokens, err
		}

		args = append(args, arg)

		switch tokens.Next().Type {
		case "comma":
			tokens = tokens.Pop()
			continue
		case "bracket_close":
			continue // The loop will end!
		}

		// If we get here, we didn't get anything we expected...
		return nil, tokens, errors.New(
			fmt.Sprintf("Unexpected %s in parameter list", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat bracket_close

	return args, tokens, nil
}
