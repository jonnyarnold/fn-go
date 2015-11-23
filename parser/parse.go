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
		}

		tokens = remainingTokens
	}

	return expressions, nil
}

func parsePrimary(tokens tokenList) (Expression, tokenList, error) {
	switch tokens.Next().Type {
	case "identifier", "number", "string", "boolean", "bracket_open", "when":
		return parseValue(tokens)
	}

	// TODO: use, import

	return nil, tokens, errors.New(
		fmt.Sprintf("Unexpected token type %s at start of expression", tokens.Next().Type),
	)
}

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

// parse the Right-Hand side of an expression.
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

func parseBlock(tokens tokenList) (BlockExpression, tokenList, error) {
	if tokens.Next().Type != "block_open" {
		return BlockExpression{}, tokens, errors.New(
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
			return BlockExpression{}, tokens, err
		}

		tokens = remainingTokens
	}

	// Check we're at a closing block.
	if tokens == nil {
		return BlockExpression{}, tokens, errors.New(
			fmt.Sprintf("End of file reached before closing block", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat block_close

	return BlockExpression{Body: body}, tokens, nil
}

func parseFunctionCall(tokens tokenList) (FunctionCallExpression, tokenList, error) {
	if tokens.Next().Type != "identifier" {
		return FunctionCallExpression{}, tokens, errors.New(
			fmt.Sprintf("Unexpected token type %s in function call", tokens.Next().Type),
		)
	}

	id := IdentifierExpression{Name: tokens.Next().Value}
	tokens = tokens.Pop() // Eat identifier

	// TODO: Error checking!
	args, tokens, _ := parseParams(tokens)

	return FunctionCallExpression{
		Identifier: id,
		Arguments:  args,
	}, tokens, nil
}

func parseParams(tokens tokenList) ([]Expression, tokenList, error) {
	if tokens.Next().Type != "bracket_open" {
		return nil, tokens, errors.New(
			fmt.Sprintf("Expected bracket_open, found %s in parameter list", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat bracket_open

	var arg Expression
	args := []Expression{}
	for tokens.Next().Type != "bracket_close" {
		arg, tokens, _ = parseParam(tokens)
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

func parseParam(tokens tokenList) (Expression, tokenList, error) {
	if tokens.Next().Type == "bracket_open" {
		return parseFunctionDefinition(tokens)
	} else {
		return parseValue(tokens)
	}
}

func parseFunctionDefinition(tokens tokenList) (FunctionPrototypeExpression, tokenList, error) {
	// TODO: Error checking!
	args, tokens, _ := parseArgs(tokens)
	body, tokens, _ := parseBlock(tokens)

	return FunctionPrototypeExpression{
		Arguments: args,
		Body:      body,
	}, tokens, nil
}

func parseArgs(tokens tokenList) ([]IdentifierExpression, tokenList, error) {
	fmt.Println("parseArgs")

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

func parseWhen(tokens tokenList) (ConditionalExpression, tokenList, error) {
	if tokens.Next().Type != "when" {
		return ConditionalExpression{}, tokens, errors.New(
			fmt.Sprintf("Expected when, found %s in when expression", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat when

	if tokens.Next().Type != "block_open" {
		return ConditionalExpression{}, tokens, errors.New(
			fmt.Sprintf("Expected block_open, found %s in when expression", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat block_open

	var (
		condition Expression
		block     BlockExpression
		err       error
	)
	branches := []ConditionalBranchExpression{}

	for tokens.Next().Type != "block_close" {
		condition, tokens, err = parseValue(tokens)
		if err != nil {
			return ConditionalExpression{}, tokens, err
		}

		block, tokens, err = parseBlock(tokens)
		if err != nil {
			return ConditionalExpression{}, tokens, err
		}

		branches = append(branches, ConditionalBranchExpression{
			Condition: condition,
			Body:      block,
		})
	}

	tokens = tokens.Pop() // Eat block_close

	return ConditionalExpression{Branches: branches}, tokens, nil
}
