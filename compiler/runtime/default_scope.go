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

func (scope defaultScope) Value() interface{} {
	return scope
}

// The top scope is the default scope used by files and REPLs.
var topScope = defaultScope{
	definitions: defMap{
		"Boolean": fn([]string{"obj"}, asBool),
		"List":    fnList{},
		"String":  fn([]string{"obj"}, callOnFirstArgument("asString")),

		"not": fn([]string{"a"}, not),
		"and": fn([]string{"a", "b"}, and),
		"or":  fn([]string{"a", "b"}, or),
		"eq":  fn([]string{"a", "b"}, eq),

		"print": fn([]string{"a"}, fnPrint),

		"+": fn([]string{"a", "b"}, callOnFirstArgument("+")),
		"-": fn([]string{"a", "b"}, callOnFirstArgument("-")),
		"*": fn([]string{"a", "b"}, callOnFirstArgument("*")),
		"/": fn([]string{"a", "b"}, callOnFirstArgument("/")),
	},
}

func DefaultScope() Scope {
	var topFnScope fnScope = topScope

	return Scope{
		parent:      &topFnScope,
		definitions: defMap{},
	}
}

func asBool(args []fnScope) (fnScope, error) {
	return FnBool(AsBool(args[0])), nil
}

func not(args []fnScope) (fnScope, error) {
	return FnBool(!AsBool(args[0])), nil
}

func callOnFirstArgument(op string) fnFunc {
	return func(args []fnScope) (fnScope, error) {
		return args[0].Definitions()[op].Call(args[1:2])
	}
}

func and(args []fnScope) (fnScope, error) {
	return FnBool(AsBool(args[0]) && AsBool(args[1])), nil
}

func or(args []fnScope) (fnScope, error) {
	return FnBool(AsBool(args[0]) || AsBool(args[1])), nil
}

func eq(args []fnScope) (fnScope, error) {
	return FnBool(args[0].Value() == args[1].Value()), nil
}

func fnPrint(args []fnScope) (fnScope, error) {
	fmt.Println(args[0].String())
	return nil, nil
}
