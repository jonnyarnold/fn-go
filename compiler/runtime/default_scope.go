package runtime

import (
	"bytes"
	"errors"
	"fmt"
)

type defaultScope struct {
	definitions defMap
}

func (scope defaultScope) Definitions() defMap {
	return scope.definitions
}

func (scope defaultScope) Define(id string, value fnScope) (fnScope, error) {
	if scope.definitions[id] != nil {
		return scope, errors.New(fmt.Sprintf("%s is already defined!", id))
	}

	scope.definitions[id] = value
	return scope, nil
}

func (scope defaultScope) String() string {
	if scope.definitions["value"] != nil {
		return scope.definitions["value"].String()
	} else {
		var str bytes.Buffer
		str.WriteString("{\n")

		for id, value := range scope.definitions {
			str.WriteString("  " + id + ": ")
			str.WriteString(value.String())
			str.WriteString("\n")
		}

		str.WriteString("}")

		return str.String()
	}
}

func (scope defaultScope) Call(args []fnScope) (fnScope, error) {
	return nil, errors.New("Default scope called as a function!")
}

// The top scope is the default scope used by files and REPLs.
var topScope = defaultScope{
	definitions: defMap{
		"+":     fn([]string{"a", "b"}, add),
		"-":     fn([]string{"a", "b"}, subtract),
		"*":     fn([]string{"a", "b"}, multiply),
		"/":     fn([]string{"a", "b"}, divide),
		"and":   fn([]string{"a", "b"}, and),
		"or":    fn([]string{"a", "b"}, or),
		"not":   fn([]string{"a"}, not),
		"eq":    fn([]string{"a", "b"}, eq),
		"print": fn([]string{"a"}, fnPrint),
		"List":  fnList{},
	},
}

func DefaultScope() Scope {
	var topFnScope fnScope = topScope

	return Scope{
		parent:      &topFnScope,
		definitions: defMap{},
	}
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

func eq(args []fnScope) (fnScope, error) {
	return FnBool(args[0] == args[1]), nil
}

func fnPrint(args []fnScope) (fnScope, error) {
	fmt.Println(args[0].String())
	return nil, nil
}
