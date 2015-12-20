package runtime

import (
	"bytes"
	"errors"
)

type defMap map[string]fnScope

// The Scope is the single object of Fn;
// it is used to represent all runtime values.
type fnScope interface {
	// Returns the definitions accessible in this scope.
	Definitions() defMap

	// Returns the current scope with the value defined.
	Define(id string, value fnScope) (fnScope, error)

	// Returns a string representation of the scope.
	String() string

	// Evalutes the scope as a function.
	Call([]fnScope) (fnScope, error)
}

type Scope struct {
	parent      *fnScope
	definitions defMap
}

func (scope Scope) Definitions() defMap {
	allDefs := scope.definitions

	if scope.parent != nil {
		for key, value := range (*scope.parent).Definitions() {

			_, ok := allDefs[key]
			if !ok {
				allDefs[key] = value
			}
		}
	}

	return allDefs
}

func (scope Scope) Define(id string, value fnScope) (fnScope, error) {
	scope.definitions[id] = value
	return scope, nil
}

func (scope Scope) String() string {
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

func (scope Scope) Call(args []fnScope) (fnScope, error) {
	if scope.definitions["call"] != nil {
		return scope.definitions["call"].Call(args)
	}

	return nil, errors.New("Scope cannot be called")
}
