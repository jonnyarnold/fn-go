package parser

import (
	"errors"
	"fmt"
)

// Parse a function definition
// A function definition has the form `(args) { primary* }`
func parseFunctionDefinition(tokens tokenList) (FunctionPrototypeExpression, tokenList, error) {
	var (
		args []IdentifierExpression
		body BlockExpression
		err  error
	)

	args, tokens, err = parseArgs(tokens)
	if err != nil {
		return FunctionPrototypeExpression{}, tokens, err
	}

	body, tokens, err = parseBlock(tokens)
	if err != nil {
		return FunctionPrototypeExpression{}, tokens, err
	}

	return FunctionPrototypeExpression{
		Arguments: args,
		Body:      body,
	}, tokens, nil
}

// Parses an argument list of the rough form `( [identifier ,]+ )`
func parseArgs(tokens tokenList) ([]IdentifierExpression, tokenList, error) {
	if tokens.Next().Type != "bracket_open" {
		return nil, tokens, errors.New(
			fmt.Sprintf("Expected bracket_open, found %s in argument list", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat bracket_open

	args := []IdentifierExpression{}
	for tokens.Next().Type != "bracket_close" {
		// Arguments are always identifiers.
		if tokens.Next().Type != "identifier" {
			return nil, tokens, errors.New(
				fmt.Sprintf("Expected identifier, found %s in argument list", tokens.Next().Type),
			)
		}

		args = append(args, IdentifierExpression{Name: tokens.Next().Value})
		tokens = tokens.Pop()

		switch tokens.Next().Type {
		case "comma":
			tokens = tokens.Pop() // Remove comma
			continue
		case "bracket_close":
			continue // The loop will end
		}

		// If we get here, we didn't get anything we expected...
		return nil, tokens, errors.New(
			fmt.Sprintf("Unexpected %s in argument list", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat bracket_close
	return args, tokens, nil
}
