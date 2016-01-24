package tokeniser

var basicTokens = map[rune]string{
	'(': "bracket_open",
	')': "bracket_close",
	',': "comma",
	';': "end_statement",
	'{': "block_open",
	'}': "block_close",
}

func tryBasicTokens(code *CodeReader) *Token {
	symbolMatch := basicTokens[code.Next()]

	if symbolMatch == "" {
		return nil
	}

	token := Token{
		Type: symbolMatch,
	}

	code.Pop() // Eat token

	return &token
}
