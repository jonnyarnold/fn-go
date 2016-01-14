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

	innerScope := Scope{
		parent:      &scope,
		definitions: defMap{},
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
		for idx, name := range argNames {
			innerScope.definitions[name] = argValues[idx]
		}

		// Evaluate the function!
		result := ExecuteIn(expr.Body.Body, innerScope)
		if result.Error != nil {
			return nil, result.Error
		}

		return result.Value, nil
	})

	return EvalResult{Value: value, Scope: scope}
}
