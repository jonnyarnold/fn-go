package runtime

import (
	"errors"
	"strconv"
)

// A fnString is a scope representing a string.
type fnBool struct {
	value bool
}

func (b fnBool) Definitions() defMap {
	return nil
}

func (b fnBool) String() string {
	return strconv.FormatBool(b.value)
}

func (b fnBool) Call(args []fnScope) (fnScope, error) {
	return nil, errors.New("Bool called as a function!")
}

func FnBool(value bool) fnBool {
	return fnBool{value: value}
}
