package runtime

import (
	"fmt"
	. "github.com/jonnyarnold/fn-go/parser"
)

// The Scope is the single object of Fn;
// it is used to represent all runtime values.
type scope interface {
	Definitions() map[string]*scope
	String() string
}

type Scope struct {
	parent      *Scope
	definitions map[string]*scope
}

func (scope Scope) Definitions() map[string]*scope {
	return scope.definitions
}

func (scope Scope) String() string {
	if scope.definitions["value"] != nil {
		return (*scope.definitions["value"]).String()
	} else {
		return "scope{}"
	}
}

// Executes the expressions in the Scope.
func ExecuteIn(exprs []Expression, scope scope) (scope, error) {
	var err error

	for _, expr := range exprs {
		scope, err = exec(expr, scope)
		if err != nil {
			return scope, err
		}
		fmt.Println(scope)
	}

	return scope, nil
}

// Executes a single expression in the Scope.
func exec(expr Expression, scope scope) (scope, error) {
	switch expr.(type) {
	case NumberExpression:
		return execNumber(expr.(NumberExpression))
	case StringExpression:
		return execString(expr.(StringExpression))
	case BooleanExpression:
		return execBool(expr.(BooleanExpression))
	}

	ignore(expr)
	return scope, nil
}

func execNumber(expr NumberExpression) (number, error) {
	return Number(expr.Value), nil
}

func execString(expr StringExpression) (fnString, error) {
	return FnString(expr.Value), nil
}

func execBool(expr BooleanExpression) (fnBool, error) {
	return FnBool(expr.Value), nil
}

// Ignores the current expression.
func ignore(expr Expression) {
	fmt.Printf("Ignoring %s\n", expr)
}
