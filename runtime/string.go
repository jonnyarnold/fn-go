package runtime

import (
	"errors"
)

// A fnString is a scope representing a string.
type fnString struct {
	value string
}

func (str fnString) Definitions() defMap {
	return nil
}

func (str fnString) Define(id string, value fnScope) (fnScope, error) {
	panic("Attempted definition on a number!")
}

func (str fnString) String() string {
	return str.value
}

func (str fnString) Call(args []fnScope) (fnScope, error) {
	return nil, errors.New("String called as a function!")
}

func FnString(str string) fnString {
	return fnString{value: str}
}
