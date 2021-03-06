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

func (fs functionScope) Definitions() defMap {
	return defMap{
		"asString": fn([]string{}, fs.asString),
	}
}

func (fn functionScope) Define(id string, value fnScope) (fnScope, error) {
	return nil, errors.New("Attempted definition on a function!")
}

func (fn functionScope) String() string {
	return fmt.Sprintf("%s { ... }", fn.ArgumentNames.String())
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

func (fn functionScope) Value() interface{} {
	return fn
}

// Helper for use when defining built-in functions.
func fn(args []string, value fnFunc) functionScope {
	return functionScope{
		ArgumentNames: args,
		value:         value,
	}
}

func (fn functionScope) asString(args []fnScope) (fnScope, error) {
	return FnString(fn.String()), nil
}
