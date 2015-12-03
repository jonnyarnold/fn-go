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
		"+":   fn([]string{"a", "b"}, add),
		"-":   fn([]string{"a", "b"}, subtract),
		"*":   fn([]string{"a", "b"}, multiply),
		"/":   fn([]string{"a", "b"}, divide),
		"and": fn([]string{"a", "b"}, and),
		"or":  fn([]string{"a", "b"}, or),
		"not": fn([]string{"a"}, not),
	},
}

func add(args []fnScope) (fnScope, error) {
	return NumberFromFloat(args[0].(number).AsFloat() + args[1].(number).AsFloat()), nil
}

func subtract(args []fnScope) (fnScope, error) {
	return NumberFromFloat(args[0].(number).AsFloat() - args[1].(number).AsFloat()), nil
}

func multiply(args []fnScope) (fnScope, error) {
	return NumberFromFloat(args[0].(number).AsFloat() * args[1].(number).AsFloat()), nil
}

func divide(args []fnScope) (fnScope, error) {
	return NumberFromFloat(args[0].(number).AsFloat() / args[1].(number).AsFloat()), nil
}

func and(args []fnScope) (fnScope, error) {
	return FnBool(args[0].(fnBool).value && args[1].(fnBool).value), nil
}

func or(args []fnScope) (fnScope, error) {
	return FnBool(args[0].(fnBool).value || args[1].(fnBool).value), nil
}

func not(args []fnScope) (fnScope, error) {
	return FnBool(!args[0].(fnBool).value), nil
}
