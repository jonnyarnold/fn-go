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
func Execute(exprs []Expression) EvalResult {
	return ExecuteIn(exprs, DefaultScope())
}

// Executes the expressions in the Scope.
func ExecuteIn(exprs []Expression, scope fnScope) EvalResult {
	var lastResult EvalResult

	for _, expr := range exprs {
		result := exec(expr, scope)
		if result.Error != nil {
			return result
		}

		lastResult = result
	}

	return lastResult
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
	case IdentifierExpression:
		return execIdentifier(expr.(IdentifierExpression), scope)
	case FunctionPrototypeExpression:
		return execFunctionPrototype(expr.(FunctionPrototypeExpression), scope)
	case BlockExpression:
		return execBlock(expr.(BlockExpression), scope)
	case FunctionCallExpression:
		return execFunctionCall(expr.(FunctionCallExpression), scope)
	case ConditionalExpression:
		return execConditional(expr.(ConditionalExpression), scope)
	}

	ignore(expr)
	return EvalResult{Scope: scope}
}

func execIdentifier(expr IdentifierExpression, scope fnScope) EvalResult {
	value := scope.Definitions()[expr.Name]
	if value == nil {
		return EvalResult{Error: errors.New(fmt.Sprintf("%s is not defined.", expr.Name))}
	}

	return EvalResult{Value: value, Scope: scope}
}

func execBlock(expr BlockExpression, scope fnScope) EvalResult {
	newBlock := Scope{
		parent:      &scope,
		definitions: defMap{},
	}

	result := ExecuteIn(expr.Body, newBlock)
	if result.Error != nil {
		return result
	}

	return EvalResult{Value: newBlock, Scope: scope}
}

// Ignores the current expression.
func ignore(expr Expression) {
	fmt.Printf("Ignoring %s\n", expr)
}
