package tokeniser

import (
	"strings"
)

var symbolInfixOperators = "+-/*.="
var stringInfixOperators = []string{"eq", "and", "or"}

func trySymbolInfixOperator(code *CodeReader) *Token {
	if !strings.ContainsRune(symbolInfixOperators, code.Next()) {
		return nil
	}

	token := Token{
		Type:  "infix_operator",
		Value: string(code.Next()),
	}

	code.Pop() // Eat infix_operator

	return &token
}

func tryStringInfixOperator(id string) *Token {
	var (
		token      Token
		tokenFound = false
	)

	for _, infixOp := range stringInfixOperators {
		if id == infixOp {
			token = Token{
				Type:  "infix_operator",
				Value: id,
			}

			tokenFound = true
			break
		}
	}

	if tokenFound {
		return &token
	} else {
		return nil
	}
}
