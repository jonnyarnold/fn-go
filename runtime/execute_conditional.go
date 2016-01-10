package runtime

import (
	"errors"
	. "github.com/jonnyarnold/fn-go/parser"
)

// Converts a FunctionPrototypeExression into a runtime function.
func execConditional(expr ConditionalExpression, scope fnScope) EvalResult {
	for _, branch := range expr.Branches {
		result := execBranch(branch, scope)

		if result.Error != nil || result.Value != nil {
			return result
		}
	}

	return EvalResult{Error: errors.New("End of when{} reached without matching branch!")}
}

func execBranch(expr ConditionalBranchExpression, scope fnScope) EvalResult {
	conditionResult := exec(expr.Condition, scope)

	if conditionResult.Error != nil {
		return conditionResult
	}

	// If the condition is true, execute the block
	if AsBool(conditionResult.Value) {
		return ExecuteIn(expr.Body.Body, scope)
	}

	return EvalResult{Value: nil}
}
