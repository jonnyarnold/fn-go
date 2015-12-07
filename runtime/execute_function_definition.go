package runtime

import (
	"errors"
	"fmt"
	. "github.com/jonnyarnold/fn-go/parser"
)

// Converts a FunctionPrototypeExression into a runtime function.
func execFunctionPrototype(expr FunctionPrototypeExpression, scope fnScope) EvalResult {
	var argNames []string
	for _, argExpr := range expr.Arguments {
		argNames = append(argNames, argExpr.Name)
	}

	value := fn(argNames, func(argValues []fnScope) (fnScope, error) {
		if len(argValues) != len(argNames) {
			return nil, errors.New(fmt.Sprintf(
				"Argument number mismatch: got %i, need %i",
				len(argValues),
				len(argNames),
			))
		}

		// Assign args to the scope
		var err error
		for idx, name := range argNames {
			scope, err = scope.Define(name, argValues[idx])
			if err != nil {
				return nil, err
			}
		}

		// Evaluate the function!
		result := ExecuteIn(expr.Body.Body, scope)
		if result.Error != nil {
			return nil, err
		}

		return result.Value, nil
	})

	return EvalResult{Value: value, Scope: scope}
}
