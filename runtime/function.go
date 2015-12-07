package runtime

import (
	"errors"
	"fmt"
	"strings"
)

// Alias for a function that can be used within fn.
type fnFunc func([]fnScope) (fnScope, error)

type argNames []string

func (names argNames) String() string {
	return fmt.Sprintf("(%s)", strings.Join(names, ", "))
}

// A functionScope wraps internal functions as scopes.
type functionScope struct {
	ArgumentNames argNames
	value         fnFunc
}

func (fn functionScope) Definitions() defMap {
	return nil
}

func (fn functionScope) Define(id string, value fnScope) (fnScope, error) {
	panic("Attempted definition on a function!")
}

func (fn functionScope) String() string {
	return fmt.Sprintf("%s { ... }", fn.ArgumentNames)
}

func (fn functionScope) Call(args []fnScope) (fnScope, error) {
	if len(args) != len(fn.ArgumentNames) {
		return nil, errors.New(fmt.Sprintf(
			"Argument number mismatch: got %i, need %i",
			len(args),
			len(fn.ArgumentNames),
		))
	}

	return fn.value(args)
}

// Helper for use when defining built-in functions.
func fn(args []string, value fnFunc) functionScope {
	return functionScope{
		ArgumentNames: args,
		value:         value,
	}
}
