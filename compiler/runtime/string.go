package runtime

import (
	"errors"
)

// A fnString is a scope representing a string.
type fnString struct {
	value string
}

func (str fnString) Definitions() defMap {
	return defMap{
		"eq":       fn([]string{"other"}, str.eq),
		"and":      fn([]string{"other"}, str.and),
		"or":       fn([]string{"other"}, str.or),
		"asString": fn([]string{}, str.asString),
	}
}

func (str fnString) Define(id string, value fnScope) (fnScope, error) {
	return nil, errors.New("Attempted definition on a string!")
}

func (str fnString) String() string {
	return str.value
}

func (str fnString) Call(args []fnScope) (fnScope, error) {
	return nil, errors.New("String called as a function!")
}

func (str fnString) Value() interface{} {
	return str.value
}

func FnString(str string) fnString {
	return fnString{value: str}
}

func (self fnString) and(args []fnScope) (fnScope, error) {
	return FnBool(AsBool(self) && AsBool(args[0])), nil
}

func (self fnString) or(args []fnScope) (fnScope, error) {
	return FnBool(AsBool(self) || AsBool(args[0])), nil
}

func (self fnString) eq(args []fnScope) (fnScope, error) {
	return FnBool(self.Value() == args[0].Value()), nil
}

func (self fnString) asString(args []fnScope) (fnScope, error) {
	return self, nil
}
