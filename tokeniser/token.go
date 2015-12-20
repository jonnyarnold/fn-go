package tokeniser

import (
	"fmt"
)

// A Token is an identified part of a code listing.
type Token struct {
	Type  string
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("[%s: %s]", t.Type, t.Value)
}

var NonToken = Token{}
