package parser

import (
	"errors"
	"fmt"
)

// Parses a when statement of the form:
// `when { [value { primary+ }]* }`
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
