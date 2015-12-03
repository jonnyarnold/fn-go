package runtime

import (
	"errors"
	"fmt"
	. "github.com/jonnyarnold/fn-go/parser"
)

// Executes the given expressions in the default scope.
func Execute(exprs []Expression) error {
	_, _, err := ExecuteIn(exprs, DefaultScope)
	return err
}

// Executes the expressions in the Scope.
func ExecuteIn(exprs []Expression, scope fnScope) (fnScope, fnScope, error) {
	var (
		value fnScope
		err   error
	)

	for _, expr := range exprs {
		value, scope, err = exec(expr, scope)
		if err != nil {
			return value, scope, err
		}
		fmt.Println(value)
	}

	return nil, scope, nil
}

// Executes a single expression in the Scope.
func exec(expr Expression, scope fnScope) (fnScope, fnScope, error) {
	switch expr.(type) {
	case NumberExpression:
		return execNumber(expr.(NumberExpression)), scope, nil
	case StringExpression:
		return execString(expr.(StringExpression)), scope, nil
	case BooleanExpression:
		return execBool(expr.(BooleanExpression)), scope, nil
	case FunctionCallExpression:
		return execFunctionCall(expr.(FunctionCallExpression), scope)
	}

	ignore(expr)
	return nil, scope, nil
}

func execNumber(expr NumberExpression) number {
	return Number(expr.Value)
}

func execString(expr StringExpression) fnString {
	return FnString(expr.Value)
}

func execBool(expr BooleanExpression) fnBool {
	return FnBool(expr.Value)
}

// Executes the given function in the scope.
// Returns the value of the function expression.
func execFunctionCall(expr FunctionCallExpression, scope fnScope) (fnScope, fnScope, error) {
	id, args := expr.Identifier.Name, expr.Arguments

	fnToCall := scope.Definitions()[id]
	if fnToCall == nil {
		return nil, scope, errors.New(fmt.Sprintf("%s is not defined.", id))
	}

	// TODO: Lazy evaluation?
	evalArgs := []fnScope{}
	for _, arg := range args {
		evalArg, _, err := exec(arg, scope)
		if err != nil {
			return nil, scope, err
		}

		evalArgs = append(evalArgs, evalArg)
	}

	value, err := fnToCall.Call(evalArgs)
	if err != nil {
		return value, scope, err
	}

	return value, scope, nil
}

// Ignores the current expression.
func ignore(expr Expression) {
	fmt.Printf("Ignoring %s\n", expr)
}
