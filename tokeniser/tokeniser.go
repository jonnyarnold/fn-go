package tokeniser

// Tokenise() converts a string input into an
// ordered array of Tokens.
func Tokenise(input string) []Token {
	var tokens = []Token{}

	// Keep looping until we have eaten the whole array.
	for input != "" {
		for _, tokenType := range TokenTypes {

			matches := tokenType.Matcher.FindStringSubmatch(input)

			if matches != nil {

				if !tokenType.Ignore {
					tokens = append(tokens, Token{
						Type:  tokenType.Tag,
						Value: matches[0],
					})
				}

				// Trim the length of the match from the front of the input.
				input = input[len(matches[0]):]

				break
			}
		}
	}

	return tokens
}
