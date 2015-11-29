package runtime

import (
	"errors"
	"fmt"
	. "github.com/jonnyarnold/fn-go/parser"
)

// Executes the given expressions in the default scope.
func Execute(exprs []Expression) error {
	_, err := ExecuteIn(exprs, DefaultScope)
	return err
}

// Executes the expressions in the Scope.
func ExecuteIn(exprs []Expression, scope fnScope) (fnScope, error) {
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
func exec(expr Expression, scope fnScope) (fnScope, error) {
	switch expr.(type) {
	case NumberExpression:
		return execNumber(expr.(NumberExpression))
	case StringExpression:
		return execString(expr.(StringExpression))
	case BooleanExpression:
		return execBool(expr.(BooleanExpression))
	case FunctionCallExpression:
		return execFunctionCall(expr.(FunctionCallExpression), scope)
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

// Executes the given function in the scope.
// Returns the value of the function expression.
func execFunctionCall(expr FunctionCallExpression, scope fnScope) (fnScope, error) {
	id, args := expr.Identifier.Name, expr.Arguments

	fnToCall := scope.Definitions()[id]
	if fnToCall == nil {
		return nil, errors.New(fmt.Sprintf("%s is not defined.", id))
	}

	// TODO: Lazy evaluation?
	evalArgs := []fnScope{}
	for _, arg := range args {
		evalArg, err := exec(arg, scope)
		if err != nil {
			return nil, err
		}

		evalArgs = append(evalArgs, evalArg)
	}

	return fnToCall.Call(evalArgs)
}

// Ignores the current expression.
func ignore(expr Expression) {
	fmt.Printf("Ignoring %s\n", expr)
}
