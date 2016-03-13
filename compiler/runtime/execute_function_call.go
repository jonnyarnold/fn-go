package runtime

import (
	"errors"
	"fmt"
	. "github.com/jonnyarnold/fn-go/compiler/parser"
)

// Executes the given function in the scope.
// Returns the value of the function expression.
func execFunctionCall(expr FunctionCallExpression, scope fnScope) EvalResult {
	id, args := expr.Identifier.Name, expr.Arguments

	// Special cases
	switch id {
	case "=":
		return execDefinition(args[0].(IdentifierExpression), args[1], scope)
	case ".":
		return execDereference(args[0], args[1], scope)
	case "import!":
		return execInternalImport(args[0].(StringExpression), scope)
	case "import":
		return execVariableImport(args[0].(StringExpression), scope)
	}

	fnToCall := scope.Definitions()[id]
	if fnToCall == nil {
		return EvalResult{Error: errors.New(fmt.Sprintf("%s is not a defined function on:\n%s", id, scope.String()))}
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

// Execute a `=` function call.
func execDefinition(id IdentifierExpression, value Expression, scope fnScope) EvalResult {
	execValue := exec(value, scope)
	if execValue.Error != nil {
		return execValue
	}

	newScope, err := scope.Define(id.Name, execValue.Value)
	if err != nil {
		return EvalResult{Error: err}
	}

	return EvalResult{
		Value: newScope,
		Scope: newScope,
	}
}

// Execute a `.` function call.
func execDereference(parent Expression, child Expression, scope fnScope) EvalResult {
	execParent := exec(parent, scope)
	if execParent.Error != nil {
		return execParent
	}

	result := exec(child, execParent.Value)
	return EvalResult{
		Value: result.Value,
		Scope: scope,
		Error: result.Error,
	}
}
