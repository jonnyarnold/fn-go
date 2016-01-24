package tokeniser

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
		stripIgnored(&code)
		if code.End() {
			continue
		}

		tokens = append(tokens, nextToken(&code))
	}

	return tokens
}

// A symbol tokeniser works on the first token of the code.
// They run before identifier tokenisers,
// and require the CodeReader object for eating.
type symbolTokeniser func(*CodeReader) *Token

// An identifier tokeniser operates on a collected identifier
// from the code. This simplifies operation, as we can work
// on the identifier and not the full CodeReader.
// They run after symbolTokenisers.
type identifierTokeniser func(string) *Token

func nextToken(code *CodeReader) Token {
	line, col := code.CurrentLine, code.CurrentColumn

	// Iterate through all of our possible symbol tokenisers,
	// looking for the first that gives us a real token.
	symbolTokenisers := []symbolTokeniser{
		tryBasicTokens,
		tryString,
		tryNumber,
		trySymbolInfixOperator,
	}

	var token *Token
	for _, tryTokeniser := range symbolTokenisers {
		token = tryTokeniser(code)
		if token != nil {
			break
		}
	}

	if token == nil {

		// Iterate through all of our possible identifier tokenisers,
		// looking for the first that gives us a real token.
		identifierTokenisers := []identifierTokeniser{
			tryBoolean,
			tryKeyword,
			tryStringInfixOperator,
		}

		// Identifier/keyword
		id := code.EatUntil(" \r\n#\"(){},;.=+-/*")
		if id == "" {
			panic("Empty identifier!")
		}

		for _, tryTokeniser := range identifierTokenisers {
			token = tryTokeniser(id)
			if token != nil {
				break
			}
		}

		// If no other tokenisers worked, we have a generic identifier.
		if token == nil {
			idToken := Token{
				Type:  "identifier",
				Value: id,
			}

			token = &idToken
		}
	}

	// The Line and Column come from the first character of the token.
	token.Line = line
	token.Column = col

	return *token
}
