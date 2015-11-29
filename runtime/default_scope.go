package runtime

import (
	"errors"
)

type defaultScope struct {
	definitions defMap
}

func (scope defaultScope) Definitions() defMap {
	return scope.definitions
}

func (scope defaultScope) String() string {
	return ""
}

func (scope defaultScope) Call(args []fnScope) (fnScope, error) {
	return nil, errors.New("Default scope called as a function!")
}

var DefaultScope = defaultScope{
	definitions: defMap{
		"+": fn([]string{"a", "b"}, add),
	},
}

func add(args []fnScope) (fnScope, error) {
	return NumberFromFloat(args[0].(number).AsFloat() + args[1].(number).AsFloat()), nil
}
