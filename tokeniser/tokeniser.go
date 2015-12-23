package tokeniser

import (
	"strings"
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

var runeInfixOperators = []rune{'+', '-', '*', '/', '.', '='}
var stringInfixOperators = []string{"eq", "and", "or"}

// Tokenise() converts a string input into an
// ordered array of Tokens.
func Tokenise(input string) []Token {
	tokens := []Token{}
	code, err := NewCodeReader(input)
	if err != nil {
		panic(err)
	}

	// Keep looping until we have eaten the whole array.
	for !code.End() {
		firstRune := code.Next()
		line, col := code.CurrentLine, code.CurrentColumn

		if firstRune == '#' {
			// Comments
			_ = code.EatUntil("\n")
		} else if firstRune == ' ' || firstRune == '\n' {
			// Spaces
			_ = code.EatWhile(" \n")
		} else if basicTokens[firstRune] != "" {
			// Basic tokens (no value)
			tokens = append(tokens, Token{
				Type:   basicTokens[firstRune],
				Line:   line,
				Column: col,
			})

			code.Pop() // Eat token
		} else if firstRune == '"' {
			// Strings
			code.Pop() // Eat opening "
			str := code.EatUntil("\"")

			tokens = append(tokens, Token{
				Type:   "string",
				Value:  str,
				Line:   line,
				Column: col,
			})

			code.Pop() // Eat closing "
		} else if strings.ContainsRune(numerics, firstRune) {
			// Numbers
			num := code.EatWhile(numerics + ".")

			tokens = append(tokens, Token{
				Type:   "number",
				Value:  num,
				Line:   code.CurrentLine,
				Column: code.CurrentColumn,
			})
		} else {
			// Rune-based infix operator
			runeFound := false
			for _, r := range runeInfixOperators {
				if firstRune == r {
					tokens = append(tokens, Token{
						Type:   "infix_operator",
						Value:  string(firstRune),
						Line:   line,
						Column: col,
					})

					code.Pop()
					runeFound = true
					break
				}
			}

			if runeFound {
				continue
			}

			// Identifier/keyword
			id := code.EatUntil(" \n#\"(){},;.=+-/*")

			if id == "true" || id == "false" {
				tokens = append(tokens, Token{
					Type:   "boolean",
					Value:  id,
					Line:   line,
					Column: col,
				})

				continue
			}

			// Check keywords
			keywordFound := false
			for _, keyword := range keywords {
				if id == keyword {
					tokens = append(tokens, Token{
						Type:   keyword,
						Line:   line,
						Column: col,
					})

					keywordFound = true
					break
				}
			}

			if keywordFound {
				continue
			}

			// Check infix operators
			infixOpFound := false
			for _, infixOp := range stringInfixOperators {
				if id == infixOp {
					tokens = append(tokens, Token{
						Type:   "infix_operator",
						Value:  id,
						Line:   line,
						Column: col,
					})

					infixOpFound = true
					break
				}
			}

			if !infixOpFound {
				tokens = append(tokens, Token{
					Type:   "identifier",
					Value:  id,
					Line:   line,
					Column: col,
				})
			}
		}
	}

	return tokens
}
