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
		"bool":  fn([]string{"a"}, asBool),
		"not":   fn([]string{"a"}, not),
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

func asBool(args []fnScope) (fnScope, error) {
	return FnBool(AsBool(args[0])), nil
}

func not(args []fnScope) (fnScope, error) {
	return FnBool(!AsBool(args[0])), nil
}

func fnPrint(args []fnScope) (fnScope, error) {
	fmt.Println(args[0].String())
	return nil, nil
}
