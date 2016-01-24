package tokeniser

func stripIgnored(code *CodeReader) {
	// Strip comments and spaces until we hit something else.
	for stripComments(code) || stripSpaces(code) {
		continue
	}
}

// Strips spaces.
// Returns true if spaces were stripped.
func stripSpaces(code *CodeReader) bool {
	firstRune := code.Next()

	if firstRune == ' ' || firstRune == '\n' || firstRune == '\r' {
		_ = code.EatWhile(" \n\r")
		return true
	}

	return false
}

// Strips comments.
// Returns true if comments were stripped.
func stripComments(code *CodeReader) bool {
	if code.Next() == '#' {
		_ = code.EatUntil("\n\r")
		return true
	}

	return false
}
