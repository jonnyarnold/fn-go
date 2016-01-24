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
	return defMap{}
}

func (b fnBool) Define(id string, value fnScope) (fnScope, error) {
	return nil, errors.New("Attempted definition on a boolean!")
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

// The definition of truth.
//
// We follow Ruby's convention - only false is false.
func AsBool(value fnScope) bool {
	switch value.(type) {
	case fnBool:
		return value.(fnBool).value
	}

	return value != nil
}
