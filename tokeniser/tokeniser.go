package tokeniser

import (
	"strings"
	"unicode/utf8"
)

var basicTokens = map[rune]string{
	'(': "bracket_open",
	')': "bracket_close",
	',': "comma",
	';': "end_statement",
	'{': "block_open",
	'}': "block_close",
}

var numerics = "0123456789"

var keywords = []string{"use", "import", "when"}

var infixOperators = []string{"+", "-", "*", "/", ".", "=", "eq", "and", "or"}

// Tokenise() converts a string input into an
// ordered array of Tokens.
func Tokenise(input string) []Token {
	var tokens = []Token{}

	// Keep looping until we have eaten the whole array.
	for input != "" {
		firstRune, firstRuneSize := utf8.DecodeRuneInString(input)

		if firstRune == '#' {
			// Comments
			_, input = eatUntil(input[firstRuneSize:], "\n")
		} else if firstRune == ' ' || firstRune == '\n' {
			// Spaces
			_, input = eatWhile(input[firstRuneSize:], " \n")
		} else if basicTokens[firstRune] != "" {
			// Basic tokens (no value)
			tokens = append(tokens, Token{
				Type: basicTokens[firstRune],
			})

			input = input[1:]
		} else if firstRune == '"' {
			// Strings
			var str string

			// We need a special eatUntil, which checks for \"

			str, input = eatUntil(input[firstRuneSize:], "\"")

			tokens = append(tokens, Token{
				Type:  "string",
				Value: str,
			})

			// The first character in the stream will be the closing ".
			// We don't want that.
			if len(input) > 0 {
				input = input[1:]
			}
		} else if strings.ContainsRune(numerics, firstRune) {
			// Numbers
			var num string
			num, input = eatWhile(input, numerics+".")

			tokens = append(tokens, Token{
				Type:  "number",
				Value: num,
			})
		} else {
			// Identifier/keyword
			var id string
			id, input = eatUntil(input, " \n#\"(){},;")

			if id == "true" || id == "false" {
				tokens = append(tokens, Token{
					Type:  "boolean",
					Value: id,
				})

				continue
			}

			// Check keywords
			keywordFound := false
			for _, keyword := range keywords {
				if id == keyword {
					tokens = append(tokens, Token{
						Type: keyword,
					})

					keywordFound = true
					break
				}
			}

			// Check infix operators
			if !keywordFound {
				infixOpFound := false
				for _, infixOp := range infixOperators {
					if id == infixOp {
						tokens = append(tokens, Token{
							Type:  "infix_operator",
							Value: id,
						})

						infixOpFound = true
						break
					}
				}

				if !infixOpFound {
					tokens = append(tokens, Token{
						Type:  "identifier",
						Value: id,
					})
				}
			}
		}
	}

	return tokens
}

// Eats the input string until one of the tokens is hit.
//
// Returns the eaten string as the first argument,
// and the remainder as the second.
//
// eatUntil("abc", "a") == "", "abc"
// eatUntil("abc", "b") == "a", "bc"
// eatUntil("abc", "c") == "ab", "c"
// eatUntil("abc", "d") == "abc", ""
func eatUntil(input string, tokens string) (string, string) {
	breakCharIdx := strings.IndexAny(input, tokens)

	if breakCharIdx == -1 {
		return input, ""
	}

	if breakCharIdx == 0 {
		return "", input[1:]
	}

	return input[:breakCharIdx], input[breakCharIdx:]
}

// Eats the input string while the tokens in the tokens list are found.
//
// Returns the eaten string as the first argument,
// and the remainder as the second.
//
// eatWhile("abc", "abc") == "abc", ""
// eatWhile("abc", "ab") == "ab", "c"
// eatWhile("abc", "a") == "a", "bc"
// eatWhile("abc", "d") == "abc", ""
func eatWhile(input string, tokens string) (string, string) {
	breakCharIdx := strings.IndexFunc(input, func(r rune) bool {
		return !strings.ContainsRune(tokens, r)
	})

	if breakCharIdx == -1 {
		return input, ""
	}

	return input[:breakCharIdx], input[breakCharIdx:]
}
