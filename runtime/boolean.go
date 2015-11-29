package runtime

import (
	"strconv"
)

// A fnString is a scope representing a string.
type fnBool struct {
	value bool
}

func (b fnBool) Definitions() map[string]*scope {
	return nil
}

func (b fnBool) String() string {
	return strconv.FormatBool(b.value)
}

func FnBool(value bool) fnBool {
	return fnBool{value: value}
}
