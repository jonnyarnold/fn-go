package tokeniser

import (
	"fmt"
)

// A Token is an identified part of a code listing.
type Token struct {
	Type   string
	Value  string
	Line   int
	Column int
}

func (t Token) String() string {
	return fmt.Sprintf("[%s %s (L%dC%d)]", t.Type, t.Value, t.Line, t.Column)
}

var NonToken = Token{}
