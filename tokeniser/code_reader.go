package tokeniser

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// Used to read fn code.
//
// Allows reading once only.
// Keeps track of things like line/column numbers.
type CodeReader struct {
	code string

	CurrentLine   int
	CurrentColumn int
}

func NewCodeReader(code string) (CodeReader, error) {
	if !utf8.Valid([]byte(code)) {
		return CodeReader{}, errors.New("String given is not valid Unicode.")
	}

	return CodeReader{
		code:          code,
		CurrentLine:   1,
		CurrentColumn: 1,
	}, nil
}

// Gets the next rune in the code.
// Does not advance the code pointer.
func (reader CodeReader) Next() rune {
	r, _ := utf8.DecodeRune([]byte(reader.code))
	return r
}

// Advances to the next rune.
// Returns the rune that was initially pointed to.
func (reader *CodeReader) Pop() rune {
	eaten := reader.Next()

	if !reader.End() {
		reader.code = reader.code[1:]
		reader.updateCurrent(string(eaten))
	}

	return eaten
}

// Returns true if the code pointer has reached
// the end of the code.
func (reader CodeReader) End() bool {
	return reader.code == ""
}

// Updates the currentIndex, CurrentLine and CurrentColumn
// after eating some characters from the code.
func (reader *CodeReader) updateCurrent(eaten string) {
	for _, r := range eaten {
		if r == '\n' {
			reader.CurrentLine += 1
			reader.CurrentColumn = 1
		} else {
			reader.CurrentColumn += 1
		}
	}
}

// Advances through the code until it hits one of the stopTokens.
// Returns all code that was advanced through.
//
// The code pointer is left on the stopToken.
func (reader *CodeReader) EatUntil(stopTokens string) string {
	var eaten string

	stopIdx := strings.IndexAny(reader.code, stopTokens)

	if stopIdx == 0 {
		// stopToken is first - Eat nothing...
		return ""
	}

	if stopIdx == -1 {
		// stopToken never hit - Eat everything...
		eaten, reader.code = reader.code, ""
	} else {
		eaten, reader.code = reader.code[:stopIdx], reader.code[stopIdx:]
	}

	reader.updateCurrent(eaten)
	return eaten
}

// Advances through the code while the continueTokens are seen.
// Returns all code that was advanced through.
//
// The code pointer is left on the first non-continueToken.
func (reader *CodeReader) EatWhile(continueTokens string) string {
	var eaten string

	stopIdx := strings.IndexFunc(reader.code, func(r rune) bool {
		return !strings.ContainsRune(continueTokens, r)
	})

	if stopIdx == -1 {
		// Eat everything
		eaten, reader.code = reader.code, ""
	} else {
		eaten, reader.code = reader.code[:stopIdx], reader.code[stopIdx:]
	}

	reader.updateCurrent(eaten)
	return eaten
}
