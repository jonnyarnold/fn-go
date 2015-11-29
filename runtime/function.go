package runtime

import (
	"errors"
	"fmt"
)

// Alias for a function that can be used within fn.
type fnFunc func([]fnScope) (fnScope, error)

// A functionScope wraps internal functions as scopes.
type functionScope struct {
	ArgumentNames []string
	value         fnFunc
}

func (fn functionScope) Definitions() defMap {
	return nil
}

func (fn functionScope) String() string {
	return fmt.Sprintf("(%s) { ... }", fn.ArgumentNames)
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
