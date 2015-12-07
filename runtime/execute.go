package runtime

import (
	"errors"
	"fmt"
	. "github.com/jonnyarnold/fn-go/parser"
)

// The result of an evaluation
type EvalResult struct {
	Value fnScope
	Scope fnScope
	Error error
}

// Executes the given expressions in the default scope.
func Execute(exprs []Expression) error {
	result := ExecuteIn(exprs, DefaultScope)
	return result.Error
}

// Executes the expressions in the Scope.
func ExecuteIn(exprs []Expression, scope fnScope) EvalResult {
	for _, expr := range exprs {
		result := exec(expr, scope)
		if result.Error != nil {
			return result
		}

		fmt.Println(result.Value)

		scope = result.Scope
	}

	return EvalResult{Scope: scope}
}

// Executes a single expression in the Scope.
func exec(expr Expression, scope fnScope) EvalResult {
	switch expr.(type) {
	case NumberExpression:
		return EvalResult{Value: execNumber(expr.(NumberExpression)), Scope: scope}
	case StringExpression:
		return EvalResult{Value: execString(expr.(StringExpression)), Scope: scope}
	case BooleanExpression:
		return EvalResult{Value: execBool(expr.(BooleanExpression)), Scope: scope}
	case FunctionCallExpression:
		return execFunctionCall(expr.(FunctionCallExpression), scope)
	}

	ignore(expr)
	return EvalResult{Scope: scope}
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
func execFunctionCall(expr FunctionCallExpression, scope fnScope) EvalResult {
	id, args := expr.Identifier.Name, expr.Arguments

	if id == "=" {
		return execDefinition(args[0].(IdentifierExpression), args[1], scope)
	}

	fnToCall := scope.Definitions()[id]
	if fnToCall == nil {
		return EvalResult{Error: errors.New(fmt.Sprintf("%s is not defined.", id))}
	}

	// TODO: Lazy evaluation?
	evalArgs, err := execArgs(args, scope)
	if err != nil {
		return EvalResult{Error: err}
	}

	value, err := fnToCall.Call(evalArgs)
	if err != nil {
		return EvalResult{Error: err}
	}

	return EvalResult{Value: value, Scope: scope}
}

func execArgs(args []Expression, scope fnScope) ([]fnScope, error) {
	// TODO: Lazy evaluation?
	evalArgs := []fnScope{}
	for _, arg := range args {
		result := exec(arg, scope)
		if result.Error != nil {
			return nil, result.Error
		}

		evalArgs = append(evalArgs, result.Value)
	}

	return evalArgs, nil
}

func execDefinition(id IdentifierExpression, value Expression, scope fnScope) EvalResult {
	execValue := exec(value, scope)
	if execValue.Error != nil {
		return execValue
	}

	return EvalResult{
		Scope: scope.Define(id.Name, execValue.Value),
	}
}

// Ignores the current expression.
func ignore(expr Expression) {
	fmt.Printf("Ignoring %s\n", expr)
}
