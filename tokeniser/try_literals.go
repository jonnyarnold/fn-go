package tokeniser

import (
	"strings"
)

var numerics = "0123456789"

func tryString(code *CodeReader) *Token {
	if code.Next() != '"' {
		return nil
	}

	code.Pop() // Eat opening "
	str := code.EatUntil("\"")

	token := Token{
		Type:  "string",
		Value: str,
	}

	code.Pop() // Eat closing "

	return &token
}

func tryNumber(code *CodeReader) *Token {
	if !strings.ContainsRune(numerics, code.Next()) {
		return nil
	}

	// Numbers
	num := code.EatWhile(numerics + ".")

	token := Token{
		Type:  "number",
		Value: num,
	}

	return &token
}

func tryBoolean(id string) *Token {
	if id != "true" && id != "false" {
		return nil
	}

	token := Token{
		Type:  "boolean",
		Value: id,
	}

	return &token
}
