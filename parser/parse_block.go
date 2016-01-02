package parser

import (
	"errors"
	"fmt"
)

// Parse a block statement.
// Blocks are of the form `{ primary+ }`
func parseBlock(tokens tokenList) (BlockExpression, tokenList, error) {
	if tokens.Next().Type != "block_open" {
		return BlockExpression{}, tokens, errors.New(
			fmt.Sprintf("Unexpected token type %s at start of block", tokens.Next().Type),
		)
	}

	tokens = tokens.Pop() // Eat block_open

	// Eat the body
	body := []Expression{}
	for tokens.Any() && tokens.Next().Type != "block_close" {
		newExpression, remainingTokens, err := parsePrimary(tokens)

		if err != nil {
			return BlockExpression{}, tokens, err
		}

		body = append(body, newExpression)
		tokens = remainingTokens
	}

	// Check we're at a closing block.
	if !tokens.Any() {
		return BlockExpression{}, tokens, errors.New("End of file reached before closing block")
	}

	tokens = tokens.Pop() // Eat block_close

	return BlockExpression{Body: body}, tokens, nil
}
